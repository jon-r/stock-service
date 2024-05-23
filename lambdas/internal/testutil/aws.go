package testutil

import "context"

type mockDynamoDbClient func(ctx context.Context)
