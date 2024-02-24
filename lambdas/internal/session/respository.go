package awsSession

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

//type SessionService interface {
//	NewAwsSession() session.Session
//}

func NewAwsSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
