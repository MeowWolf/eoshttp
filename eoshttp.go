package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/MeowWolf/eoslog"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

// HTTPClient ...
type HTTPClient struct {
	Bearer string
	Host   string
}

// Get is a standard http GET call
func (c *HTTPClient) Get(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.Host+path, nil)
	if err != nil {
		log.Error.Printf("Problem getting new http GET request from http package")
		return nil, err
	}
	req.Header.Add("Authorization", c.Bearer)

	res, err := sendRequest(client, req)
	defer res.Body.Close()
	if err != nil {
		log.Error.Printf("Problem sending http POST request")
		return nil, err
	}

	return readResponse(res)
}

// Post is a standard http POST call
func (c *HTTPClient) Post(path string, data interface{}) ([]byte, error) {
	requestJSON, err := marshallJSON(data)
	if err != nil {
		log.Error.Printf("Problem marshalling json for POST request")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.Host+path, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Error.Println("Problem getting new http POST request from http package")
		return nil, err
	}

	req.Header.Add("Authorization", c.Bearer)

	res, err := sendRequest(client, req)
	defer res.Body.Close()
	if err != nil {
		log.Error.Printf("Problem sending http POST request")
		return nil, err
	}

	return readResponse(res)
}

// Put is a standard http PUT call
func (c *HTTPClient) Put(path string, data interface{}) ([]byte, error) {
	requestJSON, err := marshallJSON(data)
	if err != nil {
		log.Error.Printf("Problem marshalling json for PUT request")
		return nil, err
	}

	requestBuffer := bytes.NewBuffer(requestJSON)
	req, err := http.NewRequest(http.MethodPut, c.Host+path, requestBuffer)
	if err != nil {
		log.Error.Printf("Problem getting new http PUT request from http package")
		return nil, err
	}

	req.Header.Add("Authorization", c.Bearer)
	res, err := sendRequest(client, req)
	defer res.Body.Close()
	if err != nil {
		log.Error.Printf("Problem sending http PUT request")
		return nil, err
	}

	return readResponse(res)
}

// Delete is a standard http DELETE call
func (c *HTTPClient) Delete(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, c.Host+path, nil)
	if err != nil {
		log.Error.Printf("Problem getting new http DELETE request from http package")
		return nil, err
	}
	req.Header.Add("Authorization", c.Bearer)

	res, err := sendRequest(client, req)
	defer res.Body.Close()
	if err != nil {
		log.Error.Printf("Problem sending http DELETE request")
		return nil, err
	}

	return readResponse(res)
}

// Is404Error checks if an error is 404
func Is404Error(err error) bool {
	return fmt.Sprint(http.StatusNotFound) == fmt.Sprint(err)
}

func marshallJSON(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error.Printf("Problem marshalling json for http request")
		return nil, err
	}

	return jsonData, err
}

func sendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := client.Do(req)
	if err != nil {
		log.Error.Printf("Problem sending http request")
		return nil, err
	}

	return res, err
}

func readResponse(res *http.Response) ([]byte, error) {
	statusCode := res.StatusCode
	if statusCode > 299 {
		status := res.Status
		var response map[string]interface{}
		json.NewDecoder(res.Body).Decode(&response)
		log.Error.Printf("%s: %s", status, response["message"])
		return nil, errors.New(fmt.Sprint(statusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error.Printf("Problem reading http response: %s", err)
	}

	return body, err
}
