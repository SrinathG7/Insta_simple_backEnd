package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Password     string             `json:"password,omitempty" bson:"-"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	PasswordHash string             `json:"passwordhash,omitempty" bson:"passwordhash,omitempty"`
}

type Post struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName        string             `json:"username,omitempty" bson:"username,omitempty"`
	Password        string             `json:"password,omitempty" bson:"-"`
	Caption         string             `json:"caption" bson:"caption"`
	ImageURL        string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	PostedTimeStamp string             `json:"postedtimestamp,omitempty" bson:"postedtimestamp,omitempty"`
}

func Post_Details(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	var DetailsOfThePost []Post
	userInput := client.Database("Insta_data").Collection("UserDetails")
	UserPointer, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := userInput.FindOne(UserPointer, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if id == user.ID {
		collection_post := client.Database("Insta_data").Collection("DetailsOfThePost")
		post_point, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection_post.Find(post_point, Post{UserName: user.Name})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(post_point)
		for cursor.Next(post_point) {
			
			var post Post
			cursor.Decode(&post)
			DetailsOfThePost = append(DetailsOfThePost, post)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(DetailsOfThePost)
	} else {
		fmt.Println("User Not Found")
		return
	}
}

func main() {
	fmt.Println("Started")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/posts", CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPost).Methods("GET")
	router.HandleFunc("/posts/users/{id}", Post_Details).Methods("GET")
	router.HandleFunc("/UserDetails", GetUserDetails).Methods("GET")
	router.HandleFunc("/DetailsOfThePost", GetDetailsOfThePost).Methods("GET")
	http.ListenAndServe(":5000", router)
}
