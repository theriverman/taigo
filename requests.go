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
)

// Evaluation Tools
var httpSuccessCodes = [...]int{
	http.StatusOK,
	http.StatusCreated,
	http.StatusAccepted,
	http.StatusNoContent,
}

// RequestService is a handle to HTTP request operations
type RequestService struct {
	client *Client
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

// Get a handler for composing a new HTTP GET request
//
//  * URL must be an absolute (full) URL to the desired endpoint
//  * ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Get(URL string, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequest("GET", s.client, ResponseBody, URL, nil)
}

// Head a handler for composing a new HTTP HEAD request
func (s *RequestService) Head() {
	panic("HEAD requests are not implemented")
}

// Post a handler for composing a new HTTP POST request
//
//  * URL must be an absolute (full) URL to the desired endpoint
//  * Payload must be a pointer to a complete struct which will be sent to Taiga
//  * ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Post(URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequest("POST", s.client, ResponseBody, URL, Payload)
}

// Put a handler for composing a new HTTP PUT request
//
//  * URL must be an absolute (full) URL to the desired endpoint
//  * Payload must be a pointer to a complete struct which will be sent to Taiga
//  * ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Put(URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequest("PUT", s.client, ResponseBody, URL, Payload)
}

// Patch a handler for composing a new HTTP PATCH request
//
//  * URL must be an absolute (full) URL to the desired endpoint
//  * Payload must be a pointer to a complete struct which will be sent to Taiga
//  * ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Patch(URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequest("PATCH", s.client, ResponseBody, URL, Payload)
}

// Delete a handler for composing a new HTTP DELETE request
//
//  * URL must be an absolute (full) URL to the desired endpoint
func (s *RequestService) Delete(URL string) (*http.Response, error) {
	return newRawRequest("DELETE", s.client, nil, URL, nil)
}

// Connect a handler for composing a new HTTP CONNECT request
func (s *RequestService) Connect() {
	panic("CONNECT requests are not implemented")
}

// Options a handler for composing a new HTTP OPTIONS request
func (s *RequestService) Options() {
	panic("OPTIONS requests are not implemented")
}

// Trace a handler for composing a new HTTP TRACE request
func (s *RequestService) Trace() {
	panic("TRACE requests are not implemented")
}

/*
func evaluateResponseAndStatusCode() {
	// Consider moving here all parts after `defer resp.Body.Close()` until the very last return statement. Keep it DRY!
}
*/

// NOTE: responseBody must always be a pointer otherwise we lose the response data!
func newfileUploadRequest(c *Client, url string, attachment *Attachment, tgObject TaigaBaseObject) (*Attachment, error) {
	// Map Object details into *Attachment
	attachment.ObjectID = tgObject.GetID()
	attachment.Project = tgObject.GetProject()

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

	// Create POST Request
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// Add headers (manually, not calling c.loadHeaders())
	request.Header.Set("Authorization", c.GetAuthorizationHeader())  // Load token
	request.Header.Set("Content-Type", writer.FormDataContentType()) // Set Content-Type to multipart/form-data

	// Execute Request
	rawResponse, err := c.HTTPClient.Do(request)
	// c.setContentTypeToJSON()  // Reset (just in case..)
	if err != nil {
		return nil, err
	}
	defer rawResponse.Body.Close()

	// Evaluate response status code && return response
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

	return nil, fmt.Errorf("Request Failed. Returned body was:\n %s", string(rawResponseBody))
}

func newRawRequest(RequestType string, c *Client, ResponseBody interface{}, URL string, Payload interface{}) (*http.Response, error) {
	// New RAW request
	var request *http.Request
	var err error

	switch {
	case Payload == nil:
		request, err = http.NewRequest(RequestType, URL, nil)
		if err != nil {
			return nil, err
		}
		break

	case Payload != nil:
		body, err := json.Marshal(Payload)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequest(RequestType, URL, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		break

	default:
		return nil, fmt.Errorf("Failed to build request because the received payload could not be processed")
	}

	// Load Headers
	c.loadHeaders(request)

	// Execute request
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Evaluate response status code
	if SuccessfulHTTPRequest(resp) {
		if resp.StatusCode == http.StatusNoContent { //  There's no body returned for 204 responses
			return resp, nil
		}
		// We expect content so convert response JSON string to struct
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&ResponseBody)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	rawResponseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("Request Failed. Returned body was:\n %s", string(rawResponseBody))
}
