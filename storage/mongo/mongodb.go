package mongo

import (
	"log"
	"net/url"
	"promhsd/db"
	"strings"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	StorageID      = "mongodb"
	collectionName = "targets"
)

type MongoDB struct {
	dbName string
	client *mongo.Client
}

func (c *MongoDB) IsHealthy() bool {
	err := c.client.Ping(context.TODO(), nil)
	if err != nil {
		return false
	}
	return true

}

func (c *MongoDB) Create(target *db.Target) error {
	coll := c.client.Database(c.dbName).Collection(collectionName)
	_, err := coll.InsertOne(context.TODO(), target)
	if err != nil {
		log.Println("Failed to Insert the document:", err)
		return err
	}
	return nil
}

func (c *MongoDB) Delete(target *db.Target) error {
	id, _ := primitive.ObjectIDFromHex(target.ID.String())
	filter := bson.D{{"_id", id}}
	coll := c.client.Database(c.dbName).Collection(collectionName)
	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("Failed to Delete the document:", err)
		return err
	}
	return nil
}

func (c *MongoDB) Get(target *db.Target) error {
	id, _ := primitive.ObjectIDFromHex(target.ID.String())
	filter := bson.D{{"_id", id}}
	coll := c.client.Database(c.dbName).Collection(collectionName)
	err := coll.FindOne(context.TODO(), filter).Decode(target)
	if err != nil {
		log.Println("Failed to Find the document:", err)
		return db.ErrNotFound
	}
	return nil
}

func (c *MongoDB) Update(target *db.Target) error {
	id, _ := primitive.ObjectIDFromHex(target.ID.String())
	filter := bson.D{{"_id", id}}
	coll := c.client.Database(c.dbName).Collection(collectionName)
	target.ID = db.ID("")
	_, err := coll.ReplaceOne(context.TODO(), filter, target)
	if err != nil {
		log.Println("Failed to Replace the document:", err)
		return err
	}
	return nil
}

func (c *MongoDB) GetAll(list *[]db.Target) error {
	coll := c.client.Database(c.dbName).Collection(collectionName)
	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Failed to Find documents", err)
		return err
	}
	if err = cur.All(context.TODO(), list); err != nil {
		log.Println("Failed to GetAll documents", err)
		return err
	}
	return nil
}

type StorageService struct{}

func (s *StorageService) ServiceID() string {
	return StorageID
}

func (s *StorageService) New(connUri string) (db.Storage, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connUri))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	u, _ := url.Parse(connUri)
	db := new(MongoDB)
	db.dbName = strings.Replace(u.Path, "/", "", 1)
	db.client = client

	return db, nil
}

func init() {
	db.RegisterStorage(&StorageService{})
}
