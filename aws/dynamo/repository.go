package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// AddItemToDB uploads struct type item to dynamodb
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

// GetItemsByPartitionKey returns items from dynamodb by partition key.
// Its return type is ([]map[string]any, error).
// Usage: MapToMessage() func in chat/util.go converts []map[string]any to []Message type
func (basics TableBasics) GetItemsByPartitionKey(ctx context.Context, partitionKey string, value types.AttributeValue) ([]map[string]any, error) {
	keyCond := fmt.Sprintf("%s = :pkVal", partitionKey)

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

func (basics TableBasics) DeleteItemsByPartitionKey(ctx context.Context, partitionKeyName string, value types.AttributeValue) error {
	keyCond := fmt.Sprintf("%s = :pkVal", partitionKeyName)

	// Get all items with the partition key
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(basics.TableName),
		KeyConditionExpression: aws.String(keyCond),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkVal": value,
		},
	}

	queryOutput, err := basics.DynamoDbClient.Query(ctx, queryInput)
	if err != nil {
		return fmt.Errorf("failed to query items: %w", err)
	}

	if len(queryOutput.Items) == 0 {
		return nil // nothing to delete
	}

	// Split items into batches of 25 (BatchWrite limit)
	const batchSize = 25
	for i := 0; i < len(queryOutput.Items); i += batchSize {
		end := i + batchSize
		if end > len(queryOutput.Items) {
			end = len(queryOutput.Items)
		}
		batch := queryOutput.Items[i:end]

		writeRequests := make([]types.WriteRequest, 0, len(batch))

		for _, item := range batch {
			// Extract both PK and SK
			key := make(map[string]types.AttributeValue)
			key[partitionKeyName] = item[partitionKeyName]

			// find sort key automatically
			for k := range item {
				if k != partitionKeyName {
					key[k] = item[k]
				}
			}

			writeRequests = append(writeRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: key,
				},
			})
		}

		// Execute BatchWriteItem
		_, err := basics.DynamoDbClient.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				basics.TableName: writeRequests,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to batch delete items: %w", err)
		}
	}

	return nil
}
