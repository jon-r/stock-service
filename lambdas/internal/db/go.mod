module jon-richards.com/stock-app/db

go 1.22.0

replace jon-richards.com/stock-app/remote => ../remote

require (
	github.com/aws/aws-sdk-go v1.50.25
	github.com/google/uuid v1.6.0
	jon-richards.com/stock-app/remote v0.0.0-00010101000000-000000000000
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect
