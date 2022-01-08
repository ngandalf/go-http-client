package gohttpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/google/go-querystring/query"
)

var timeout = time.Duration(2 * time.Second)

// New func returns a Client interface
func New(baseUrl string, token string) Client {
	return &client{
		BaseUrl: baseUrl,
		Token:   token}
}

// Get func returns a request
func (h client) Get(endpoint string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", h.BaseUrl, endpoint), bytes.NewBuffer([]byte{}))
}

// GetWith func returns a request
func (h client) GetWith(endpoint string, params interface{}) (*http.Request, error) {
	queryString, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s?%s", h.BaseUrl, endpoint, queryString.Encode()), bytes.NewBuffer([]byte{}))
}

// Post func returns a request
func (h client) Post(endpoint string) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, h.BaseUrl+endpoint, bytes.NewBuffer([]byte{}))
}

// PostWith func returns a request
func (h client) PostWith(endpoint string, params interface{}) (*http.Request, error) {
	json, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPost, h.BaseUrl+endpoint, bytes.NewBuffer(json))
}

// Put func returns a request
func (h client) Put(endpoint string) (*http.Request, error) {
	return http.NewRequest(http.MethodPut, h.BaseUrl+endpoint, bytes.NewBuffer([]byte{}))
}

// PutWith func returns a request
func (h client) PutWith(endpoint string, params interface{}) (*http.Request, error) {
	json, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPut, h.BaseUrl+endpoint, bytes.NewBuffer(json))
}

// Patch func returns a request
func (h client) Patch(endpoint string) (*http.Request, error) {
	return http.NewRequest(http.MethodPatch, h.BaseUrl+endpoint, bytes.NewBuffer([]byte{}))
}

// PatchWith func returns a request
func (h client) PatchWith(endpoint string, params interface{}) (*http.Request, error) {
	json, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPatch, h.BaseUrl+endpoint, bytes.NewBuffer(json))
}

// Delete func returns a request
func (h client) Delete(endpoint string) (*http.Request, error) {
	return http.NewRequest(http.MethodDelete, h.BaseUrl+endpoint, bytes.NewBuffer([]byte{}))
}

// DeleteWith func returns a request
func (h client) DeleteWith(endpoint string, params interface{}) (*http.Request, error) {
	json, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodDelete, h.BaseUrl+endpoint, bytes.NewBuffer(json))
}

//func dialTimeout(network, addr string) (net.Conn, error) {
//	return net.DialTimeout(network, addr, timeout)
//}

// Do func returns a response with your data
func (h client) Do(request *http.Request) (Response, error) {

	fmt.Println("define transport")
	transport := http.Transport{
		//Dial: dialTimeout,
		Dial: (&net.Dialer{
			// Modify the time to wait for a connection to establish
			Timeout:   1 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: &transport,
		Timeout:   4 * time.Second,
	}

	// Create a Bearer string by appending string access token
	var auth = "Token " + h.Token

	// add authorization header to the req
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", auth)

	dump, err := httputil.DumpRequestOut(request, true)
	if err == nil {
		fmt.Printf("%s\n", dump)
	}

	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	fmt.Printf("request: %v %v, Time taken: %v\n", request.Method, request.RequestURI, elapsed)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &ResponseStruct{
		Status:        response.Status,
		StatusCode:    response.StatusCode,
		Header:        response.Header,
		ContentLength: response.ContentLength,
		Body:          body,
	}, nil
}

// Get func returns ResponseStruct struct of request
func (r ResponseStruct) Get() ResponseStruct {
	return r
}

// To func returns converts string to struct
func (r ResponseStruct) To(value interface{}) {
	err := json.Unmarshal(r.Body, &value)
	if err != nil {
		value = nil
	}
}
