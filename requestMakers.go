package gotaiga

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Evaluation Tools
const noContent204 = 204

func newRawRequest(requestType string, c *Client, responseBody interface{}, url string, payload interface{}) error {
	// New RAW request
	var request *http.Request
	if payload == nil {
		request, _ = http.NewRequest(requestType, url, nil)
	} else if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		request, _ = http.NewRequest(requestType, url, bytes.NewBuffer(body))
	} else {
		c.Logger.Panicln("Failed to build request in newRawRequest. Received payload could not be processed!")
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
	if successfulRequest(resp) {
		if resp.StatusCode == noContent204 { //  There's no body returned for 204 responses
			return nil
		}
		// Collect returned Pagination headers
		p := Pagination{}
		p.LoadFromHeaders(c, resp)
		// We expect content so convert response JSON string to struct
		json.Unmarshal([]byte(body), &responseBody) // responseBody contains a pointer to a struct
		return nil
	}

	c.Logger.Println("Failed Request!")
	return errors.New(string(body))
}

func successfulRequest(response *http.Response) bool {
	for _, okStatusCode := range [4]int{200, 201, 202, 204} {
		if response.StatusCode == okStatusCode {
			return true
		}
	}
	return false
}

func getRequest(c *Client, responseBody interface{}, url string) error {
	return newRawRequest("GET", c, responseBody, url, nil)
}

func headRequest() {}

func postRequest(c *Client, responseBody interface{}, url string, payload interface{}) error {
	// NOTE: responseBody must always be a pointer otherwise we lose the response data!
	// New POST request
	return newRawRequest("POST", c, responseBody, url, payload)
}

func putRequest(c *Client, responseBody interface{}, url string, payload interface{}) error {
	// NOTE: responseBody must always be a pointer otherwise we lose the response data!
	// New PUT request
	return newRawRequest("PUT", c, responseBody, url, payload)
}

func patchRequest(c *Client, responseBody interface{}, url string, payload interface{}) error {
	/*
		patchRequest is proto-function to keep the code DRY and write things only once
		patchRequest takes the following arguments:
		  * Client (pointer)
		  * responseBody (interface storing a pointer to a struct)
		  * url (string)
		  * payload (interface storing a pointer to a struct)
	*/
	// New PATCH request
	return newRawRequest("PATCH", c, responseBody, url, payload)
}

func deleteRequest(c *Client, url string) error {
	// New DELETE request
	return newRawRequest("DELETE", c, nil, url, nil)
}

func connectRequest() {}

func optionsRequest() {}

func traceRequest() {}

func evaluateResponseAndStatusCode() {
	// Consider moving here all parts after `defer resp.Body.Close()` until the very last return statement. Keep it DRY!
}

// // NOTE: responseBody must always be a pointer otherwise we lose the response data!
func newfileUploadRequest(c *Client, url string, attachment *Attachment) (*Attachment, error) {
	// Open file
	f, _ := os.Open(attachment.filePath)
	fileName := filepath.Base(attachment.filePath)
	defer f.Close()

	// Prepare request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("attached_file", fileName)
	io.Copy(part, f)

	// Add object_id & project to the form-data
	writer.WriteField("object_id", strconv.Itoa(attachment.ObjectID))
	writer.WriteField("project", strconv.Itoa(attachment.Project))
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
	if successfulRequest(rawResponse) {
		var responseBody Attachment
		// We expect content so convert response JSON string to struct
		json.Unmarshal([]byte(rawResponseBody), &responseBody) // responseBody contains a pointer to a struct
		return &responseBody, nil
	}

	c.Logger.Println("Failed Request!")
	return nil, errors.New(string(rawResponseBody))
}
