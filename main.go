package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Create AWS DynamoDB service resources
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}
	sess := session.New(config)
	dynamoSvc := dynamodb.New(sess)

	// Create Dynamo DB table
	/* createTableOutput, err := dynamoSvc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
		},
		TableName: aws.String("People"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	checkError(err)

	// Wait until the table is created before proceeding...
	fmt.Println("Waiting for table to be active")
	dynamoTableName := *createTableOutput.TableDescription.TableName
	err = dynamoSvc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(dynamoTableName),
	})
	fmt.Println("Table has been successfully created!")
	fmt.Println(createTableOutput) */

	addItemOutput, err := dynamoSvc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":   {N: aws.String("0")},
			"Name": {S: aws.String("Brandon")},
		},
		TableName: aws.String("People"),
	})
	checkError(err)
	fmt.Println(addItemOutput.GoString())
}
