package DAL

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Kibuns/SearchService/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// global variable mongodb connection client
var client mongo.Client = NewClient()

// ----Create----
func InsertTwoot(twoot Models.Twoot) {
	twootCollection := client.Database("HashtagTwootDB").Collection("twoots")
	twoot.Created = time.Now()
	_, err := twootCollection.InsertOne(context.TODO(), twoot)
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	// return the ID of the newly inserted script
	fmt.Println("New twoot inserted for the user named: " + twoot.UserName)
}

//----Read----

func SearchTwootsByHashtag(hashtag string) ([]bson.M, error) {
    twootCollection := client.Database("HashtagTwootDB").Collection("twoots")
    
    // construct the query to search for the hashtag in the twoot.Hashtags array
    filter := bson.M{"hashtags": bson.M{"$in": []string{hashtag}}}
    
    // execute the query and retrieve the results
    cursor, err := twootCollection.Find(context.TODO(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    
    // iterate through the cursor and decode the documents into a slice of bson.M
    var results []bson.M
    for cursor.Next(context.Background()) {
        var doc bson.M
        if err := cursor.Decode(&doc); err != nil {
            return nil, err
        }
        results = append(results, doc)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    
    return results, nil
}

func ReadAllTwoots() (values []primitive.M) {
	twootCollection := client.Database("HashtagTwootDB").Collection("twoots")
	// retrieve all the documents (empty filter)
	cursor, err := twootCollection.Find(context.TODO(), bson.D{})
	// check for errors in the finding
	if err != nil {
		panic(err)
	}

	// convert the cursor result to bson
	var results []bson.M
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// display the documents retrieved
	fmt.Println("displaying all results from the search query")
	for _, result := range results {
		fmt.Println(result)
	}

	values = results
	return
}



func ReadSingleTwoot(id string) (value primitive.M) {
	twootCollection := client.Database("TwootDB").Collection("twoots")
	// convert the hexadecimal string to an ObjectID type
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	// retrieve the document with the specified _id
	var result bson.M
	err = twootCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&result)
	if err != nil {
		panic(err)
	}

	// display the retrieved document
	fmt.Println("displaying the result from the search query")
	fmt.Println(result)
	value = result

	return value
}

//----Update----

//----Delete----

// other
func NewClient() (value mongo.Client) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://ninoverhaegh:6P77TACMZwsd8pb4@twotterdb.jfx1rk2.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	value = *client

	return
}
