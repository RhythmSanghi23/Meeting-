package main

import (
	"context"
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Meeting struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             `json:"title,omitempty" bson:"title,omitempty"`
	Participants  string             `json:"participants,omitempty" bson:"participants,omitempty"`
	Start string `json:"start,omitempty" bson:"start,omitempty"`
	End string `json:"end,omitempty" bson:"end,omitempty"`
	Creation_timestamp string `json:"creation_timestamp,omitempty" bson:"creation_timestamp,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	RSVP string `json:"rsvp,omitempty" bson:"rsvp,omitempty"`


}
func ScheduleMeeting(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meeting Meeting
	_ = json.NewDecoder(request.Body).Decode(&meeting)
	collection := client.Database("rhythm").Collection("meeting")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var meeting Meeting
	collection := client.Database("rhythm").Collection("meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting) }


func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) { 
response.Header().Set("content-type", "application/json")
	var people []Meeting
	collection := client.Database("rhythm").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		people = append(people, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)

	}

func GetTimeFrame(response http.ResponseWriter, request *http.Request) {
response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	start, _ := primitive.ObjectIDFromHex(params["start"])
	end, _ := primitive.ObjectIDFromHex(params["end"])
	var meeting Meeting
	collection := client.Database("rhythm").Collection("meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{Start: start, End: end}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting) }

func GetParticipants(response http.ResponseWriter, request *http.Request) {
response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	email, _ := primitive.ObjectIDFromHex(params["email"])
	var meeting Meeting
	collection := client.Database("rhythm").Collection("meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{Email:email}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting) }


func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/meetings", ScheduleMeeting).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/meetings/{ID}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/meeting/{start & end}", GetTimeFrame).Method("GET")
	router.HandleFunc("/participants/{Email}", GetParticipants).Methods("GET")

	http.ListenAndServe(":12345", router)
	fmt.Println("Ending the application...")

}


