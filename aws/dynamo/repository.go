package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Iemt : struct type for DynamoDB item
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

func (basics TableBasics) GetItemsByPartitionKey(ctx context.Context, partitionKeyName string, value types.AttributeValue) ([]map[string]any, error) {
	keyCond := fmt.Sprintf("%s = :pkVal", partitionKeyName)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(basics.TableName),
		KeyConditionExpression: aws.String(keyCond),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkVal": value,
		},
	}

	var items []map[string]any
	paginator := dynamodb.NewQueryPaginator(basics.DynamoDbClient, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("query error: %w", err)
		}

		var pageItems []map[string]any
		err = attributevalue.UnmarshalListOfMaps(page.Items, &pageItems)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal page items: %w", err)
		}

		items = append(items, pageItems...)
	}

	return items, nil
}
