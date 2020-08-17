package paste

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

type Paste struct {
	DynamodbClient dynamodbiface.ClientAPI
	Opts           Options
}

type Options struct {
	TableName string
	File      string
}

func (p *Paste) Start() {
	f, err := os.Open(p.Opts.File)
	utils.CheckError("Failed to open file", err)

	decoder := json.NewDecoder(bufio.NewReader(f))

	getItems(decoder, func(item map[string]interface{}) {
		p.writeItem(item)
	})

	// close
	err = f.Close()
	utils.CheckError("Failed to close file", err)

	// wp := workerpool.New(4)

}

func (p *Paste) writeItem(item map[string]interface{}) {
	v, err := dynamodbattribute.MarshalMap(item)
	utils.CheckError("Failed to marshal item", err)

	req := p.DynamodbClient.PutItemRequest(&dynamodb.PutItemInput{
		TableName: &p.Opts.TableName,
		Item:      v,
	})

	_, err = req.Send(context.Background())
	utils.CheckError("Failed to perform put request", err)
}

func getItems(decoder *json.Decoder, processItem func(map[string]interface{})) {
	for decoder.More() {
		var m map[string]interface{}
		err := decoder.Decode(&m)
		utils.CheckError("Error while decoding backup file", err)
		processItem(m)
	}
}
