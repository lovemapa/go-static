package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	Name       string `json:name`
	City       string `json:city`
	ProfilePic string `json:profilePic`
}

var userCollection = db().Database("goTest").Collection("profile") // get collection "users" from db() which returns *mongo.Client

func createProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("static", "uploadFile-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	person := user{

		Name:       r.FormValue("name"),
		City:       r.FormValue("city"),
		ProfilePic: "/" + path.Base(tempFile.Name()),
	}

	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document

}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                   //slice for multiple documents
	cur, err := userCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)
}
