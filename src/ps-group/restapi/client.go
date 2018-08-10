package restapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Client - existing JSON REST API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient - creates new JSON REST API client
func NewClient(baseURL string) *Client {
	client := new(Client)
	client.httpClient = new(http.Client)
	client.baseURL = baseURL
	return client
}

// Post - calls POST method of the API, writes results into `result` parameter
func (c *Client) Post(method string, params interface{}, result interface{}) error {
	url := c.baseURL + method
	requestBytes, err := json.Marshal(params)
	if err != nil {
		return errors.Wrap(err, "cannot encode request for POST method "+method)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return errors.Wrap(err, "cannot create request for POST method "+method+", url="+url)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "cannot send request for POST method "+method)
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response for POST method "+method)
	}

	err = json.Unmarshal(responseBytes, result)
	if err != nil {
		return errors.Wrap(err, "cannot parse JSON response for POST method "+method)
	}

	return nil
}

// Get - calls GET method of the API, writes results into `result` parameter
func (c *Client) Get(method string, result interface{}) error {
	url := c.baseURL + method
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create request for GET method "+method)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "cannot send request for GET method "+method)
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response for GET method "+method)
	}

	err = json.Unmarshal(responseBytes, result)
	if err != nil {
		return errors.Wrap(err, "cannot parse JSON response for GET method "+method)
	}

	return nil
}
