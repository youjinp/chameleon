package copy

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/youjinp/chameleon/pkg/utils"
)

type Copy struct {
	DynamodbClient dynamodbiface.ClientAPI
	Opts           Options
}

type Options struct {
	TableName string
	Limit     int64
	Output    string
}

func (c *Copy) Start() {
	// create writer
	f, err := os.Create(c.Opts.Output)
	utils.CheckError("Failed to create file", err)

	bufferedWriter := bufio.NewWriter(f)
	encoder := json.NewEncoder(bufferedWriter)

	// dynamodb scan & write
	c.scan(func(items []map[string]dynamodb.AttributeValue) {
		var m []map[string]interface{}
		err := dynamodbattribute.UnmarshalListOfMaps(items, &m)
		utils.CheckError("unable to parse dynamodb output", err)
		for _, item := range m {
			err = encoder.Encode(item)
			utils.CheckError("unable to write to backup file", err)
		}
	})

	// close
	err = bufferedWriter.Flush()
	utils.CheckError("Failed to flush buffered writer", err)

	err = f.Close()
	utils.CheckError("Failed to close file", err)
}

func (c *Copy) scan(processItems func([]map[string]dynamodb.AttributeValue)) {
	req := c.DynamodbClient.ScanRequest(&dynamodb.ScanInput{
		TableName: &c.Opts.TableName,
		Limit:     &c.Opts.Limit,
	})

	resp, err := req.Send(context.TODO())
	utils.CheckError("Failed to perform scan", err)
	processItems(resp.Items)

	lastKey := resp.LastEvaluatedKey
	for lastKey != nil {
		req := c.DynamodbClient.ScanRequest(&dynamodb.ScanInput{
			TableName:         &c.Opts.TableName,
			ExclusiveStartKey: lastKey,
			Limit:             &c.Opts.Limit,
		})

		resp, err := req.Send(context.TODO())
		utils.CheckError("Failed to perform scan", err)
		processItems(resp.Items)
		lastKey = resp.LastEvaluatedKey
	}
}
