package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// nolint
func NewDynamoDBClient(awsCfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(awsCfg)
}

// nolint
func NewTableBasics(dynamoDbClient *dynamodb.Client, tableName string) TableBasics {
	return TableBasics{
		DynamoDbClient: dynamoDbClient,
		TableName:      tableName,
	}
}
