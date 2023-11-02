package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/go-plugin"
	"github.com/philipsahli/cacher/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type KV struct {
	Logger hclog.Logger
}

type Item struct {
	Key   string
	Value string
}

func (kv KV) Get(key string) ([]byte, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("eu-central-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		return []byte{}, err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	item := Item{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("MyTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return []byte{}, err
	}
	if result.Item == nil {
		msg := "Could not find '" + key
		return []byte{}, errors.New(msg)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	kv.Logger.Info(fmt.Sprintf("Key %s read from backend (dynamodb)", key))

	return []byte(item.Value), nil
}

func (kv KV) Put(key string, value []byte) error {

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("eu-central-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		return err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Prepare the item input
	item := map[string]*dynamodb.AttributeValue{
		"key": {
			S: aws.String(key),
		},
		"value": {
			S: aws.String(string(value)),
		},
	}

	// Put the item
	input := &dynamodb.PutItemInput{
		TableName: aws.String("MyTable"),
		Item:      item,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:", err)
		return err
	}

	kv.Logger.Info(fmt.Sprintf("Key %s written to backend (dynamodb)", key))

	return nil

}

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	logger.Debug("message from plugin", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.KVGRPCPlugin{Impl: &KV{
				Logger: logger,
			}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})

}
