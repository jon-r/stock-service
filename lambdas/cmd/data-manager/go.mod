module github.com/jon-r/stock-service/data-manager

go 1.22.0

replace (
	jon-richards.com/stock-app/db => ../../internal/db
	jon-richards.com/stock-app/queue => ../../internal/queue
	jon-richards.com/stock-app/remote => ../../internal/remote
)

require (
	github.com/aws/aws-lambda-go v1.46.0
	jon-richards.com/stock-app/db v0.0.0-00010101000000-000000000000
	jon-richards.com/stock-app/queue v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	jon-richards.com/stock-app/remote v0.0.0-00010101000000-000000000000 // indirect
)
