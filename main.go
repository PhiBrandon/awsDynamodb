package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var homeTemplate *template.Template
var data Data

type Data struct {
	People []Person
}
type Person struct {
	Id        int
	Age       int
	Email     string
	FirstName string
	LastName  string
}

func createPersonItem(p Person) *dynamodb.PutItemInput {
	av, err := dynamodbattribute.MarshalMap(p)
	checkError(err)
	//fmt.Println(av.String())
	item := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("People"),
	}
	return item
}

func putPersons(p []Person, dynamoSvc *dynamodb.DynamoDB) {

	for _, person := range p {
		putItemInput := createPersonItem(person)

		fmt.Printf("Adding %v to Dyanmo...\n", person.FirstName)
		_, err := dynamoSvc.PutItem(putItemInput)
		checkError(err)
		fmt.Println(person.FirstName + " has been added to dyanmo.")
	}
}

//Web page to use data
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, data); err != nil {
		panic(err)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var people []Person
	var err error
	homeTemplate, err = template.ParseFiles("views/index.gohtml")
	checkError(err)
	// Create AWS DynamoDB service resources
	persons := []Person{
		Person{
			Id:        0,
			FirstName: "Brandon",
			LastName:  "Phillips",
			Age:       30,
			Email:     "thedefinedone@gmail.com",
		},
		Person{
			Id:        1,
			FirstName: "Taylor",
			LastName:  "Scahefer",
			Age:       22,
			Email:     "tshaef@gmail.com",
		},
		Person{
			Id:        2,
			FirstName: "Joe",
			LastName:  "Ruck",
			Age:       40,
			Email:     "jru@gmail.com",
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
					AttributeName: aws.String("Id"),
					AttributeType: aws.String("N"),
				},
			},
			TableName: aws.String("People"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Id"),
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
		putPersons(persons, dynamoSvc)
	} else {
		fmt.Println("Table name: " + *checktableOutput.Table.TableName)

		putPersons(persons, dynamoSvc)
	}

	queryOutput, err := dynamoSvc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("People"),
	})
	checkError(err)

	record := []Person{}
	err = dynamodbattribute.UnmarshalListOfMaps(queryOutput.Items, &record)
	checkError(err)
	people = append(people, record...)
	for _, person := range people {
		fmt.Println(person.FirstName)
	}
	data = Data{
		People: people,
	}
	http.HandleFunc("/", home)
	http.ListenAndServe(":8080", nil)

}
