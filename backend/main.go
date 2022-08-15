package main

//port 8069
import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Student struct {
	ID     primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	Name   string             `json:"name, omitempty" bson:"firstname, omitempty"`
	Branch string             `json:"branch, omitempty" bson:"branch, omitempty"`
	UserID string             `json:"userid, omitempty" bson:"userid, omitempty"`
	RollNo string             `json:"rollno, omitempty" bson:"rollno, omitempty"`
}

//global declaration
var client *mongo.Client
var database = "pclubbackend"
var collection = "students"

//endpoint: create_record
//request type: post
func add_student(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("content-type", "application/json")

	var student Student
	json.NewDecoder(request.Body).Decode(&student)
	student.ID = primitive.NewObjectIDFromTimestamp(time.Now())

	collection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	result, _ := collection.InsertOne(ctx, student)
	json.NewEncoder(resp).Encode(result)
}

//endpoint: get_all
//request type: get
//returns a list of students.
func get_students(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("content-type", "application/json")

	var students []Student
	collection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// var params Student
	// json.NewDecoder(request.Body).Decode(&params)
	// fmt.Printf("%v\n", params)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"message" : "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var student Student
		cursor.Decode(&student)
		students = append(students, student)
	}
	if err := cursor.Err(); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"message" : "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(resp).Encode(students)

}

//endpoint: filter/{by}/{value}
//gets all records which have by: value
func filter_students_by_value(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("content-type", "application/json")

	var students []Student
	collection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	params := mux.Vars(request)
	by := params["by"]
	val := params["value"]
	// fmt.Println(by)
	// fmt.Println(val)
	//

	//usually works out of the box, have to make cases for id_ and rollno
	var cursor *mongo.Cursor
	var err error
	if by == "_id" {
		id, _ := primitive.ObjectIDFromHex(val)
		cursor, err = collection.Find(ctx, bson.M{by: id})
	} else {
		cursor, err = collection.Find(ctx, bson.M{by: val})
	}

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"message" : "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	// println("Lol")
	for cursor.Next(ctx) {
		var student Student
		cursor.Decode(&student)
		// fmt.Printf("%v\n", student)
		students = append(students, student)

	}
	if err := cursor.Err(); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"message" : "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(resp).Encode(students)

}

//endpoint: update?_id=id
//only id because nothing else is unique
//fuck there are two Rahul Jhas in Y21.
func update(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("content-type", "application/json")
	query := request.URL.Query()
	// fmt.Printf("%v", query)
	// rollno, roll_present := query["rollno"]
	id_, id_present := query["_id"]
	if !id_present || len(id_) == 0 {
		resp.Write([]byte(`{"message" : "id not present"}`))
		return
	}
	lmao, _ := primitive.ObjectIDFromHex(id_[0])

	var student Student
	collection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	collection.FindOne(ctx, bson.M{"_id": lmao}).Decode(&student)

	var marshalled interface{}
	body, _ := ioutil.ReadAll(request.Body)
	_ = json.Unmarshal([]byte(body), &marshalled)

	result, _ := collection.UpdateOne(ctx, bson.M{"_id": lmao}, bson.M{"$set": marshalled})
	json.NewEncoder(resp).Encode(result)
}

//endpoint: delete
//POST Request, deletes everything that matches a filter
func delete(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("content-type", "application/json")
	var marshalled interface{}
	body, _ := ioutil.ReadAll(request.Body)
	_ = json.Unmarshal([]byte(body), &marshalled)

	collection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	result, _ := collection.DeleteMany(ctx, marshalled)
	json.NewEncoder(resp).Encode(result)

}
func main() {
	// fmt.Println("Rahhul Best")
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
	// fmt.Println("Rahhul Best")
	router := mux.NewRouter()
	router.HandleFunc("/create_record", add_student).Methods("POST")
	router.HandleFunc("/get_all", get_students).Methods("GET")
	router.HandleFunc("/filter/{by}/{value}", filter_students_by_value).Methods("GET")
	router.HandleFunc("/update", update).Methods("POST")
	router.HandleFunc("/delete", delete).Methods("POST")
	http.ListenAndServe(":12345", router)
}
