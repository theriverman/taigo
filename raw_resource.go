package taigo

import (
	"fmt"
	"net/http"
)

// RawResource is a generic JSON object container for endpoints without dedicated DTOs yet.
type RawResource map[string]any

func listRawResources(c *Client, endpoint string, defaultProjectID int, queryParams any) ([]RawResource, error) {
	url := c.MakeURL(endpoint)
	var err error
	url, err = urlWithQueryOrDefaultProject(url, queryParams, defaultProjectID)
	if err != nil {
		return nil, err
	}
	var resources []RawResource
	_, err = c.Request.Get(url, &resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func getRawResource(c *Client, endpoint string, resourceID any) (*RawResource, error) {
	url := c.MakeURL(endpoint, fmt.Sprint(resourceID))
	var resource RawResource
	_, err := c.Request.Get(url, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func createRawResource(c *Client, endpoint string, payload any) (*RawResource, error) {
	url := c.MakeURL(endpoint)
	var resource RawResource
	_, err := c.Request.Post(url, payload, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func putRawResource(c *Client, endpoint string, resourceID any, payload any) (*RawResource, error) {
	url := c.MakeURL(endpoint, fmt.Sprint(resourceID))
	var resource RawResource
	_, err := c.Request.Put(url, payload, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func patchRawResource(c *Client, endpoint string, resourceID any, payload any) (*RawResource, error) {
	url := c.MakeURL(endpoint, fmt.Sprint(resourceID))
	var resource RawResource
	_, err := c.Request.Patch(url, payload, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func deleteRawResource(c *Client, endpoint string, resourceID any) (*http.Response, error) {
	url := c.MakeURL(endpoint, fmt.Sprint(resourceID))
	return c.Request.Delete(url)
}

func getRawResourceAtPath(c *Client, endpointParts ...string) (*RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var resource RawResource
	_, err := c.Request.Get(url, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func getRawResourceAtPathWithQuery(c *Client, queryParams any, endpointParts ...string) (*RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var err error
	url, err = appendQueryParams(url, queryParams)
	if err != nil {
		return nil, err
	}
	var resource RawResource
	_, err = c.Request.Get(url, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func getRawResourceListAtPath(c *Client, endpointParts ...string) ([]RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var resources []RawResource
	_, err := c.Request.Get(url, &resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func getRawResourceListAtPathWithQuery(c *Client, queryParams any, endpointParts ...string) ([]RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var err error
	url, err = appendQueryParams(url, queryParams)
	if err != nil {
		return nil, err
	}
	var resources []RawResource
	_, err = c.Request.Get(url, &resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func postRawResourceAtPath(c *Client, payload any, endpointParts ...string) (*RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var resource RawResource
	_, err := c.Request.Post(url, payload, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func postRawResourceAtPathWithQuery(c *Client, payload any, queryParams any, endpointParts ...string) (*RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var err error
	url, err = appendQueryParams(url, queryParams)
	if err != nil {
		return nil, err
	}
	var resource RawResource
	_, err = c.Request.Post(url, payload, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func postRawResourceListAtPath(c *Client, payload any, endpointParts ...string) ([]RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var resources []RawResource
	_, err := c.Request.Post(url, payload, &resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func putRawResourceAtPath(c *Client, payload any, endpointParts ...string) (*RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var resource RawResource
	_, err := c.Request.Put(url, payload, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func patchRawResourceAtPath(c *Client, payload any, endpointParts ...string) (*RawResource, error) {
	url := c.MakeURL(endpointParts...)
	var resource RawResource
	_, err := c.Request.Patch(url, payload, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func deleteRawResourceAtPath(c *Client, endpointParts ...string) (*http.Response, error) {
	url := c.MakeURL(endpointParts...)
	return c.Request.Delete(url)
}
