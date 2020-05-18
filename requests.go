package taigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Evaluation Tools
const noContent204 = 204

var httpSuccessCodes = [...]int{200, 201, 202, 204}

// RequestService is a handle to HTTP request operations
type RequestService struct {
	client *Client
}

// MakeURL accepts an Endpoint URL and returns a compiled absolute URL
//
// For example:
//	* If the given endpoint URL is /epic-custom-attributes
//	* If the BaseURL is https://api.taiga.io
//	* It returns https://api.taiga.io/api/v1/epic-custom-attributes
func (s *RequestService) MakeURL(Endpoint string) string {
	if strings.HasPrefix(Endpoint, "/") {
		return s.client.APIURL + Endpoint
	}
	return s.client.APIURL + "/" + Endpoint
}

// SuccessfulHTTPRequest returns true if the given Response's StatusCode
// is one of `[...]int{200, 201, 202, 204}`; otherwise returns false
// Taiga does not return status codes other than above stated
func SuccessfulHTTPRequest(Response *http.Response) bool {
	for _, code := range httpSuccessCodes {
		if Response.StatusCode == code {
			return true
		}
	}
	return false
}

// GetRequest a handler for composing a new HTTP GET request
//
// URL must be an absolute (full) URL to the desired endpoint
// ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) GetRequest(URL string, ResponseBody interface{}) error {
	return newRawRequest("GET", s.client, ResponseBody, URL, nil)
}

// HeadRequest a handler for composing a new HTTP HEAD request
func (s *RequestService) HeadRequest() {
	panic("HEAD requests are not implemented")
}

// PostRequest a handler for composing a new HTTP POST request
//
// URL must be an absolute (full) URL to the desired endpoint
// Payload must be a pointer to a complete struct which will be sent to Taiga
// ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) PostRequest(URL string, Payload interface{}, ResponseBody interface{}) error {
	// NOTE: responseBody must always be a pointer otherwise we lose the response data!
	// New POST request
	return newRawRequest("POST", s.client, ResponseBody, URL, Payload)
}

// PutRequest a handler for composing a new HTTP PUT request
//
// URL must be an absolute (full) URL to the desired endpoint
// Payload must be a pointer to a complete struct which will be sent to Taiga
// ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) PutRequest(URL string, Payload interface{}, ResponseBody interface{}) error {
	// NOTE: responseBody must always be a pointer otherwise we lose the response data!
	// New PUT request
	return newRawRequest("PUT", s.client, ResponseBody, URL, Payload)
}

// PatchRequest a handler for composing a new HTTP PATCH request
//
// URL must be an absolute (full) URL to the desired endpoint
// Payload must be a pointer to a complete struct which will be sent to Taiga
// ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) PatchRequest(URL string, Payload interface{}, ResponseBody interface{}) error {
	/*
		patchRequest is proto-function to keep the code DRY and write things only once
		patchRequest takes the following arguments:
		  * Client (pointer)
		  * responseBody (interface storing a pointer to a struct)
		  * url (string)
		  * payload (interface storing a pointer to a struct)
	*/
	// New PATCH request
	return newRawRequest("PATCH", s.client, ResponseBody, URL, Payload)
}

// DeleteRequest a handler for composing a new HTTP DELETE request
//
// URL must be an absolute (full) URL to the desired endpoint
func (s *RequestService) DeleteRequest(URL string) error {
	// New DELETE request
	return newRawRequest("DELETE", s.client, nil, URL, nil)
}

// ConnectRequest a handler for composing a new HTTP CONNECT request
func (s *RequestService) ConnectRequest() {
	panic("CONNECT requests are not implemented")
}

// OptionsRequest a handler for composing a new HTTP OPTIONS request
func (s *RequestService) OptionsRequest() {
	panic("OPTIONS requests are not implemented")
}

// TraceRequest a handler for composing a new HTTP TRACE request
func (s *RequestService) TraceRequest() {
	panic("TRACE requests are not implemented")
}

/*
func evaluateResponseAndStatusCode() {
	// Consider moving here all parts after `defer resp.Body.Close()` until the very last return statement. Keep it DRY!
}
*/

// NOTE: responseBody must always be a pointer otherwise we lose the response data!
func newfileUploadRequest(c *Client, url string, attachment *Attachment) (*Attachment, error) {
	// Open file
	f, err := os.Open(attachment.filePath)
	if err != nil {
		return nil, fmt.Errorf("Could not open file at specified location: " + attachment.filePath)
	}
	fileName := filepath.Base(attachment.filePath)
	defer f.Close()

	// Prepare request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("attached_file", fileName)
	if err != nil {
		return nil, fmt.Errorf("Could not write file to buffer")
	}
	io.Copy(part, f)

	// Add object_id & project to the form-data
	writer.WriteField("object_id", strconv.Itoa(attachment.ObjectID))
	writer.WriteField("project", strconv.Itoa(attachment.Project))
	writer.WriteField("from_comment", "False")
	writer.Close()

	// Add headers && Execute Request
	req, _ := http.NewRequest("POST", url, body)
	c.setContentType(writer.FormDataContentType())
	c.loadHeaders(req)
	rawResponse, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rawResponse.Body.Close()

	// Evaluate response status code && return reponse
	rawResponseBody, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}
	if SuccessfulHTTPRequest(rawResponse) {
		var responseBody Attachment
		// We expect content so convert response JSON string to struct
		json.Unmarshal([]byte(rawResponseBody), &responseBody) // responseBody contains a pointer to a struct
		return &responseBody, nil
	}

	return nil, fmt.Errorf("Request Failed. Returned body was:\n %s", rawResponseBody)
}

func newRawRequest(RequestType string, c *Client, ResponseBody interface{}, URL string, Payload interface{}) error {
	// New RAW request
	var request *http.Request
	var err error

	switch {
	case Payload == nil:
		request, err = http.NewRequest(RequestType, URL, nil)
		if err != nil {
			return err
		}
		break

	case Payload != nil:
		body, err := json.Marshal(Payload)
		if err != nil {
			return err
		}
		request, err = http.NewRequest(RequestType, URL, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		break

	default:
		return fmt.Errorf("Failed to build request because the received payload could not be processed")
	}

	// Load Headers
	c.loadHeaders(request)

	// Execute request
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Evaluate response status code
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if SuccessfulHTTPRequest(resp) {
		if resp.StatusCode == noContent204 { //  There's no body returned for 204 responses
			return nil
		}
		// Collect returned Pagination headers
		p := Pagination{}
		p.LoadFromHeaders(c, resp)
		// We expect content so convert response JSON string to struct
		json.Unmarshal([]byte(body), &ResponseBody) // responseBody contains a pointer to a struct
		return nil
	}

	return fmt.Errorf("Request Failed. Returned body was:\n %s", body)
}
