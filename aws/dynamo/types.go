package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PartiQLRunner struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}
