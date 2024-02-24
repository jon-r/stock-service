package awsSession

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewAwsSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
