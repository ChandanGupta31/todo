package usecases

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/models"
	"todo/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoService struct {
	MongoCollection *mongo.Collection
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (svc *TodoService) CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	var todo *models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Print("decoding error", err)
		w.WriteHeader(http.StatusBadRequest)
		res.Error = err.Error()
		return
	}
	todo.ID = uuid.NewString()

	repo := repository.TodoRepo{MongoCollection: svc.MongoCollection}

	result, err := repo.CreateTodo(todo)
	if err != nil {
		log.Print("insertion error", err)
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	w.WriteHeader(http.StatusOK)
	res.Data = result
}

func (svc *TodoService) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	todoID := mux.Vars(r)["id"]
	if todoID == "" {
		log.Print("invalid id")
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "invalid id"
		return
	}

	var todo *models.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = todoID

	repo := repository.TodoRepo{MongoCollection: svc.MongoCollection}
	result, err := repo.UpdateTodo(todoID, todo)
	if err != nil {
		log.Print("internal server error", err)
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	w.WriteHeader(http.StatusOK)
	res.Data = result
}

func (svc *TodoService) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	todoID := mux.Vars(r)["id"]
	if todoID == "" {
		log.Print("invalid id")
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "invalid id"
		return
	}

	repo := repository.TodoRepo{MongoCollection: svc.MongoCollection}
	result, err := repo.DeleteTodo(todoID)
	if err != nil {
		log.Print("internal sever error", err)
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	w.WriteHeader(http.StatusOK)
	res.Data = result
}

func (svc *TodoService) GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.TodoRepo{MongoCollection: svc.MongoCollection}
	todos, err := repo.GetTODOs()
	if err != nil {
		log.Print("internal server error", err)
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	w.WriteHeader(http.StatusOK)
	res.Data = todos
}
