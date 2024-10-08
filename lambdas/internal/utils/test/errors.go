package test

import (
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubbedError(err error) *testtools.StubError {
	if err != nil {
		return &testtools.StubError{Err: err, ContinueAfter: true}
	}

	return nil
}
