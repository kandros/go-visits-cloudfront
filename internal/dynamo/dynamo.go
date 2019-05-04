package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

const tableName = "go_visits"

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var db = dynamodb.New(sess)

func Store(in interface{}) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	output, err := db.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		return nil, err
	}

	return output, nil
}
