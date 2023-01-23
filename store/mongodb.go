package store

import (
	"context"

	"github.com/pallat/todoapi/todo"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStore struct {
	*mongo.Collection
}

func NewMongoDBStore(col *mongo.Collection) *MongoDBStore {
	return &MongoDBStore{Collection: col}
}

func (s *MongoDBStore) New(todo *todo.Todo) error {
	_, err := s.Collection.InsertOne(context.Background(), todo)
	return err
}

func (s *MongoDBStore) Find(todo *[]todo.Todo) error {
	_, err := s.Collection.Find(context.Background(), todo)
	return err
}

func (s *MongoDBStore) Delete(todo *todo.Todo, id int) error {
	todo.ID = uint(id)
	_, err := s.Collection.DeleteOne(context.Background(), todo)
	return err
}
