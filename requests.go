package taigo

import (
	"bytes"
	"context"
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

// APIError represents a non-2xx response from Taiga.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("taiga API error (status=%d): %s", e.StatusCode, e.Body)
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
//   - URL must be an absolute (full) URL to the desired endpoint
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Get(URL string, ResponseBody interface{}) (*http.Response, error) {
	return s.GetCtx(context.Background(), URL, ResponseBody)
}

// GetCtx composes a new HTTP GET request with context.
func (s *RequestService) GetCtx(ctx context.Context, URL string, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "GET", s.client, ResponseBody, URL, nil)
}

// Head a handler for composing a new HTTP HEAD request
func (s *RequestService) Head(URL string, ResponseBody interface{}) (*http.Response, error) {
	return s.HeadCtx(context.Background(), URL, ResponseBody)
}

// HeadCtx composes a new HTTP HEAD request with context.
func (s *RequestService) HeadCtx(ctx context.Context, URL string, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "HEAD", s.client, ResponseBody, URL, nil)
}

// Post a handler for composing a new HTTP POST request
//
//   - URL must be an absolute (full) URL to the desired endpoint
//   - Payload must be a pointer to a complete struct which will be sent to Taiga
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Post(URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return s.PostCtx(context.Background(), URL, Payload, ResponseBody)
}

// PostCtx composes a new HTTP POST request with context.
func (s *RequestService) PostCtx(ctx context.Context, URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "POST", s.client, ResponseBody, URL, Payload)
}

// Put a handler for composing a new HTTP PUT request
//
//   - URL must be an absolute (full) URL to the desired endpoint
//   - Payload must be a pointer to a complete struct which will be sent to Taiga
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Put(URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return s.PutCtx(context.Background(), URL, Payload, ResponseBody)
}

// PutCtx composes a new HTTP PUT request with context.
func (s *RequestService) PutCtx(ctx context.Context, URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "PUT", s.client, ResponseBody, URL, Payload)
}

// Patch a handler for composing a new HTTP PATCH request
//
//   - URL must be an absolute (full) URL to the desired endpoint
//   - Payload must be a pointer to a complete struct which will be sent to Taiga
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Patch(URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return s.PatchCtx(context.Background(), URL, Payload, ResponseBody)
}

// PatchCtx composes a new HTTP PATCH request with context.
func (s *RequestService) PatchCtx(ctx context.Context, URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "PATCH", s.client, ResponseBody, URL, Payload)
}

// Delete a handler for composing a new HTTP DELETE request
//
//   - URL must be an absolute (full) URL to the desired endpoint
func (s *RequestService) Delete(URL string) (*http.Response, error) {
	return s.DeleteCtx(context.Background(), URL)
}

// DeleteCtx composes a new HTTP DELETE request with context.
func (s *RequestService) DeleteCtx(ctx context.Context, URL string) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "DELETE", s.client, nil, URL, nil)
}

// Connect a handler for composing a new HTTP CONNECT request
func (s *RequestService) Connect(URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return s.ConnectCtx(context.Background(), URL, Payload, ResponseBody)
}

// ConnectCtx composes a new HTTP CONNECT request with context.
func (s *RequestService) ConnectCtx(ctx context.Context, URL string, Payload interface{}, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "CONNECT", s.client, ResponseBody, URL, Payload)
}

// Options a handler for composing a new HTTP OPTIONS request
func (s *RequestService) Options(URL string, ResponseBody interface{}) (*http.Response, error) {
	return s.OptionsCtx(context.Background(), URL, ResponseBody)
}

// OptionsCtx composes a new HTTP OPTIONS request with context.
func (s *RequestService) OptionsCtx(ctx context.Context, URL string, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "OPTIONS", s.client, ResponseBody, URL, nil)
}

// Trace a handler for composing a new HTTP TRACE request
func (s *RequestService) Trace(URL string, ResponseBody interface{}) (*http.Response, error) {
	return s.TraceCtx(context.Background(), URL, ResponseBody)
}

// TraceCtx composes a new HTTP TRACE request with context.
func (s *RequestService) TraceCtx(ctx context.Context, URL string, ResponseBody interface{}) (*http.Response, error) {
	return newRawRequestWithContext(ctx, "TRACE", s.client, ResponseBody, URL, nil)
}

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
		return nil, fmt.Errorf("could not write file to buffer")
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("could not copy file data to request body: %w", err)
	}

	// Add object_id & project to the form-data
	if err := writer.WriteField("object_id", strconv.Itoa(attachment.ObjectID)); err != nil {
		return nil, fmt.Errorf("could not set object_id field: %w", err)
	}
	if err := writer.WriteField("project", strconv.Itoa(attachment.Project)); err != nil {
		return nil, fmt.Errorf("could not set project field: %w", err)
	}
	if err := writer.WriteField("from_comment", "False"); err != nil {
		return nil, fmt.Errorf("could not set from_comment field: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("could not finalize multipart body: %w", err)
	}

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

	return nil, &APIError{
		StatusCode: rawResponse.StatusCode,
		Body:       string(rawResponseBody),
	}
}

func newRawRequest(RequestType string, c *Client, ResponseBody interface{}, URL string, Payload interface{}) (*http.Response, error) {
	return newRawRequestWithContext(context.Background(), RequestType, c, ResponseBody, URL, Payload)
}

func newRawRequestWithContext(ctx context.Context, RequestType string, c *Client, ResponseBody interface{}, URL string, Payload interface{}) (*http.Response, error) {
	// New RAW request
	var request *http.Request
	var err error

	switch Payload {
	case nil:
		request, err = http.NewRequestWithContext(ctx, RequestType, URL, nil)
		if err != nil {
			return nil, err
		}
	default:
		body, err := json.Marshal(Payload)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequestWithContext(ctx, RequestType, URL, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
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
		if ResponseBody != nil {
			// We expect content so convert response JSON string to struct
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(ResponseBody)
			if err != nil {
				return nil, err
			}
		}
		return resp, nil
	}

	rawResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	return resp, &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(rawResponseBody),
	}
}
