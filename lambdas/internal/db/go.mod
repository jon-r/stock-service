module jon-richards.com/stock-app/db

go 1.22.0

replace jon-richards.com/stock-app/remote => ../remote

require (
	github.com/aws/aws-sdk-go-v2 v1.25.2
	github.com/aws/aws-sdk-go-v2/config v1.27.4
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.13.6
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.7.6
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.30.1
	github.com/google/uuid v1.6.0
	jon-richards.com/stock-app/remote v0.0.0-00010101000000-000000000000
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.1 // indirect
	github.com/aws/smithy-go v1.20.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
