package taigo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

/*
	// Attachments Manager
	since go does not have generics, a common attachments manager had to be created
	Taiga objects (epics, tasks, issues, etc...) can wrap around this method to simplify otherwise redundant requests
*/

// listAttachmentsForEndpoint is a common method to get attachments for an endpoint (userstories, tasks, etc...)
func listAttachmentsForEndpoint(c *Client, queryParams *attachmentsQueryParams) ([]Attachment, error) {
	if queryParams == nil {
		return nil, fmt.Errorf("queryParams must not be nil")
	}
	url := c.MakeURL(queryParams.endpointURI, "attachments")
	paramValues, err := query.Values(queryParams)
	if err != nil {
		return nil, fmt.Errorf("encode attachment query params: %w", err)
	}
	if encoded := paramValues.Encode(); encoded != "" {
		url += "?" + encoded
	}
	var attachments []Attachment
	_, err = c.Request.Get(url, &attachments)
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// getAttachmentForEndpoint is a common method to get a specific attachment for an endpoint (epic, issue, etc...)
func getAttachmentForEndpoint(c *Client, attachmentID int, endpointURI string) (*Attachment, error) {
	url := c.MakeURL(endpointURI, "attachments", strconv.Itoa(attachmentID))
	var a Attachment
	_, err := c.Request.Get(url, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// convertStructViaJSON takes a model struct and converts it to another struct
//
// Since Type Conversion (https://golang.org/ref/spec#Conversions) is limited to identical types in go,
// JSON is used as an intermediate language to achieve this functionality
//
// NOTE: Both `sourcePtr` and `targetPtr` MUST BE POINTERS!
func convertStructViaJSON(sourcePtr any, targetPtr any) error {
	payloadInJSON, err := json.Marshal(sourcePtr)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(payloadInJSON), targetPtr)
	if err != nil {
		return err
	}
	return nil
}

// isEmpty is a generic-ish function to check if a struct's field is empty/default
// it is convenient when making sure the bare minimum values are set when creating an object
func isEmpty(structField any) bool {
	if structField == nil {
		return true
	}
	v := reflect.ValueOf(structField)
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice, reflect.Func:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

// projectIDQueryParam returns gives project ID formatted as a generic QueryParam
func projectIDQueryParam(ProjectID int) string {
	return "?project=" + strconv.Itoa(ProjectID)
}

// BoolPtr returns a pointer to the given bool.
func BoolPtr(v bool) *bool {
	return &v
}

// requireNonNil validates that a required input is non-nil.
func requireNonNil(name string, value any) error {
	if value == nil {
		return fmt.Errorf("%s must not be nil", name)
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice, reflect.Func:
		if v.IsNil() {
			return fmt.Errorf("%s must not be nil", name)
		}
	}
	return nil
}

// applyDefaultProjectToQuery fills the `project` query parameter when omitted.
func applyDefaultProjectToQuery(queryParams any, defaultProjectID int) {
	if queryParams == nil || defaultProjectID == 0 {
		return
	}
	v := reflect.ValueOf(queryParams)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	s := v.Elem()
	if s.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		structField := s.Type().Field(i)
		tag := structField.Tag.Get("url")
		key := strings.Split(tag, ",")[0]
		if key != "project" || !field.CanSet() {
			continue
		}
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() == 0 {
				field.SetInt(int64(defaultProjectID))
			}
		case reflect.Ptr:
			if field.IsNil() && field.Type().Elem().Kind() == reflect.Int {
				projectID := defaultProjectID
				field.Set(reflect.ValueOf(&projectID))
			}
		}
		return
	}
}

func cloneQueryParams(queryParams any) any {
	if queryParams == nil {
		return nil
	}
	v := reflect.ValueOf(queryParams)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return queryParams
	}
	clone := reflect.New(v.Elem().Type())
	clone.Elem().Set(v.Elem())
	return clone.Interface()
}

// appendQueryParams appends encoded query parameters to baseURL.
func appendQueryParams(baseURL string, queryParams any) (string, error) {
	if queryParams == nil {
		return baseURL, nil
	}
	paramValues, err := query.Values(queryParams)
	if err != nil {
		return "", fmt.Errorf("encode query params: %w", err)
	}
	if encoded := paramValues.Encode(); encoded != "" {
		return baseURL + "?" + encoded, nil
	}
	return baseURL, nil
}

// urlWithQueryOrDefaultProject applies query filters and falls back to default project.
func urlWithQueryOrDefaultProject(baseURL string, queryParams any, defaultProjectID int) (string, error) {
	if queryParams != nil {
		encodedParams := cloneQueryParams(queryParams)
		applyDefaultProjectToQuery(encodedParams, defaultProjectID)
		return appendQueryParams(baseURL, encodedParams)
	}
	if defaultProjectID != 0 {
		return baseURL + projectIDQueryParam(defaultProjectID), nil
	}
	return baseURL, nil
}

// tagsToNames extracts tag names from Taiga's [][]string tag format.
// Each tag entry is expected to have at least one element (the tag name).
func tagsToNames(tags Tags) []string {
	if len(tags) == 0 {
		return nil
	}
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag) == 0 {
			continue
		}
		names = append(names, tag[0])
	}
	return names
}

// namesToTags converts plain tag names to Taiga's [][]string tag format.
func namesToTags(names ...string) Tags {
	if len(names) == 0 {
		return nil
	}
	tags := make(Tags, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		tags = append(tags, []string{name})
	}
	if len(tags) == 0 {
		return nil
	}
	return tags
}
