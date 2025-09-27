package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (basics TableBasics) AddItemToDB(ctx context.Context, item any) error {
	it, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	_, err = basics.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: it,
	})
	if err != nil {
		return fmt.Errorf("couldn't marshal item: %w", err)
	}
	return err
}
