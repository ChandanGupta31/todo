package repository

import (
	"context"
	"log"
	"todo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepo struct {
	MongoCollection *mongo.Collection
}

func (r *TodoRepo) CreateTodo(todo *models.Todo) (interface{}, error) {
	_, err := r.MongoCollection.InsertOne(context.Background(), todo)
	if err != nil {
		log.Print("insertion error", err)
		return nil, err
	}
	return todo.ID, nil
}

func (r *TodoRepo) UpdateTodo(id string, todo *models.Todo) (int64, error) {
	result, err := r.MongoCollection.UpdateOne(context.Background(), bson.D{{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: todo}})
	if err != nil {
		log.Print("updation error", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (r *TodoRepo) DeleteTodo(id string) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		log.Print("deletion error", err)
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *TodoRepo) GetTODOs() ([]models.Todo, error) {
	result, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Print("Not getting todos", err)
		return nil, err
	}

	var todos []models.Todo
	err = result.All(context.Background(), &todos)
	if err != nil {
		log.Print("conversion fail", err)
		return nil, err
	}

	return todos, nil
}
