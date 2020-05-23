package taigo

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"
)

/*
	// Attachments Manager
	since go does not have generics, a common attachments manager had to be created
	Taiga objects (epics, tasks, issues, etc...) can wrap around this method to simplify otherwise redundant requests
*/

// listAttachmentsForEndpoint is a common method to get attachments for an endpoint (userstories, tasks, etc...)
func listAttachmentsForEndpoint(c *Client, queryParams *attachmentsQueryParams) ([]Attachment, error) {
	paramValues, _ := query.Values(queryParams)
	url := c.MakeURL(queryParams.endpointURI, "attachments?", paramValues.Encode())
	var attachments []Attachment
	_, err := c.Request.Get(url, &attachments)
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// getAttachmentForEndpoint is a common method to get a specific attachment for an endpoint (epic, issue, etc...)
func getAttachmentForEndpoint(c *Client, attachmentID int, endpointURI string) (*Attachment, error) {
	url := c.MakeURL(endpointURI, fmt.Sprintf("attachments/%d", attachmentID))
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
// JSON is used as an intermediate language to achive this functionality
//
// NOTE: Both `sourcePtr` and `targetPtr` MUST BE POINTERS!
func convertStructViaJSON(sourcePtr interface{}, targetPtr interface{}) error {
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
func isEmpty(structField interface{}) bool {
	if structField == nil {
		return true
	} else if structField == "" {
		return true
	} else if structField == false {
		return true
	}
	return false
}
