package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Person struct {
	id        int
	age       int
	email     string
	firstName string
	lastName  string
}

func createPersonItem(p *Person) *dynamodb.PutItemInput {
	item := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":        {N: aws.String(strconv.Itoa(p.id))},
			"firstName": {S: aws.String(p.firstName)},
			"lastName":  {S: aws.String(p.lastName)},
			"age":       {N: aws.String(strconv.Itoa(p.age))},
			"email":     {S: aws.String(p.email)},
		},
		TableName: aws.String("People"),
	}
	return item
}

func main() {
	// Create AWS DynamoDB service resources
	persons := []*Person{
		&Person{
			id:        0,
			firstName: "Brandon",
			lastName:  "Phillips",
			age:       30,
			email:     "thedefinedone@gmail.com",
		},
		&Person{
			id:        1,
			firstName: "Taylor",
			lastName:  "Scahefer",
			age:       22,
			email:     "tshaef@gmail.com",
		},
		&Person{
			id:        2,
			firstName: "Joe",
			lastName:  "Ruck",
			age:       40,
			email:     "jru@gmail.com",
		},
	}
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}
	sess := session.New(config)
	dynamoSvc := dynamodb.New(sess)

	// Check if people table exists...
	checktableOutput, err := dynamoSvc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String("People"),
	})
	if err != nil {
		fmt.Println("Table has not been created... Creating Table...")
		createTableOutput, err := dynamoSvc.CreateTable(&dynamodb.CreateTableInput{
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
		dynamoTableName := *createTableOutput.TableDescription.TableName
		fmt.Println("Waiting for table to be active")
		err = dynamoSvc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
			TableName: aws.String(dynamoTableName),
		})
		fmt.Println("Table has been successfully created!")
		fmt.Println(createTableOutput)
	}
	fmt.Println("Table name: " + *checktableOutput.Table.TableName)

	for _, person := range persons {
		putItemInput := createPersonItem(person)

		fmt.Printf("Adding %v to Dyanmo...\n", person.firstName)
		_, err := dynamoSvc.PutItem(putItemInput)
		checkError(err)
		fmt.Println(person.firstName + " has been added to dyanmo.")
	}
}
