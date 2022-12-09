package dynamo

import (
	"errors"
	"log"
	"promhsd/db"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	StorageID = "dynamodb"
)

type DynamoDB struct {
	tableName string
	svc       *dynamodb.DynamoDB
}

func (d *DynamoDB) Create(target *db.Target) error {
	target.ID = db.ID(target.Name)
	av, err := dynamodbattribute.MarshalMap(target)
	if err != nil {
		return err
	}

	_, err = d.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      av,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *DynamoDB) Delete(target *db.Target) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(target.ID.String()),
			},
		},
		TableName: aws.String(d.tableName),
	}

	_, err := d.svc.DeleteItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				return db.ErrNotFound
			}
		}
		return err
	}
	return nil
}

func (d *DynamoDB) Get(target *db.Target) error {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(target.ID.String()),
			},
		},
		TableName: aws.String(d.tableName),
	}

	result, err := d.svc.GetItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				return db.ErrNotFound
			}
		}
		return err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, target)
	if err != nil {
		return err
	}
	return nil
}

func (d *DynamoDB) Update(target *db.Target) error {
	av, err := dynamodbattribute.MarshalMap(target)
	if err != nil {
		return err
	}

	_, err = d.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      av,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				return db.ErrNotFound
			}
		}
		return err
	}
	return nil
}

func (d *DynamoDB) GetAll(list *[]db.Target) error {
	input := &dynamodb.ScanInput{
		TableName: aws.String(d.tableName),
	}
	result, err := d.svc.Scan(input)
	if err != nil {
		return err
	}
	targets := make([]db.Target, *result.Count, *result.Count)
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &targets)
	if err != nil {
		return err
	}
	*list = targets
	return nil
}

func (d *DynamoDB) createTable() error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
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
		TableName: aws.String(d.tableName),
	}
	_, err := d.svc.CreateTable(input)
	if err != nil {
		resourceInUseException := &dynamodb.ResourceInUseException{}
		if errors.As(err, &resourceInUseException) {
			return nil
		}
		log.Println(err.Error())
		return err
	}
	return nil
}

type StorageService struct{}

func (s *StorageService) ServiceID() string {
	return StorageID
}

func (s *StorageService) New(tableName string) (db.Storage, error) {
	db := new(DynamoDB)
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	db.svc = dynamodb.New(sess)
	db.tableName = tableName
	err = db.createTable()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func init() {
	db.RegisterStorage(&StorageService{})
}
