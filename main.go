package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"todo/usecases"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func init() {
	// loading env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error in loading env file", err)
	}
	log.Print("env file loaded")

	// connecting to database
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("error in connecting mongodb", err)
	}
	log.Print("mongodb connected successfully")

	// ping
	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("error pinging mongodb", err)
	}
	log.Print("ping successful")
}

func main() {
	// closing connection gracefully
	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			log.Print("error in disconnecting mongodb")
		}
	}()

	// accessing database
	coll := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	// setting services
	todoServices := usecases.TodoService{MongoCollection: coll}

	// initializing routers
	r := mux.NewRouter()

	// handling routes
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/todo", todoServices.CreateTodo).Methods(http.MethodPost)
	r.HandleFunc("/todo/{id}", todoServices.UpdateTodo).Methods(http.MethodPut)
	r.HandleFunc("/todo/{id}", todoServices.DeleteTodo).Methods(http.MethodDelete)
	r.HandleFunc("/todo", todoServices.GetTodos).Methods(http.MethodGet)

	// starting server
	log.Print("starting server at 3331")
	err := http.ListenAndServe(":3331", r)
	if err != nil {
		log.Fatal("error in starting server", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("services are working ..."))
}
