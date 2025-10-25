package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// nolint
type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}
