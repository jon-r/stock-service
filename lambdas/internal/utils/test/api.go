package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type ReqStub struct {
	URL    string
	Method string
	Input  string
	Output string
}

type ApiRequestStubber struct {
	requests []ReqStub
	reqIndex int
}

func (stubber *ApiRequestStubber) AddRequest(stub ReqStub) {
	stubber.requests = append(stubber.requests, stub)
}

func (stubber *ApiRequestStubber) Next() *ReqStub {
	if stubber.reqIndex < len(stubber.requests) {
		stub := stubber.requests[stubber.reqIndex]
		stubber.reqIndex += 1
		return &stub
	} else {
		return nil
	}
}

func (stubber *ApiRequestStubber) NewTestClient() *http.Client {
	handler := func(req *http.Request) *http.Response {
		var err error
		stub := stubber.Next()

		if stub == nil {
			fmt.Printf("no stub found for request URL %s\n", req.URL.Path)
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}
		}

		err = stubber.compare(*stub, *req)

		if err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: io.NopCloser(bytes.NewBufferString(stub.Output)),
		}
	}

	return &http.Client{
		Transport: RoundTripFunc(handler),
	}
}

func (stubber *ApiRequestStubber) compare(expected ReqStub, actual http.Request) error {
	var err error

	// todo use assert test stuff?
	if expected.Method != actual.Method {
		err = fmt.Errorf("expected method %s, got %s", expected.Method, actual.Method)
	}
	if expected.URL != actual.URL.String() {
		err = fmt.Errorf("expected URL %s, got %s", expected.URL, actual.URL.String())
	}

	// todo handle body? (not needed as api requests are all get atm

	return err
}

func NewApiStubber() *ApiRequestStubber {
	return &ApiRequestStubber{}
}

func (stubber *ApiRequestStubber) VerifyAllStubsCalled() error {
	var err error
	next := stubber.Next()
	if next != nil {
		err = fmt.Errorf("remaining stub %v was never called", next.URL)
	}
	return err
}
