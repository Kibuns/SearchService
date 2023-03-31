package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kibuns/SearchService/DAL"
	"github.com/Kibuns/SearchService/Models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	handleRequests()
}

//controllers

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAll")
	json.NewEncoder(w).Encode(getAllTwoots())
}

func returnTwoot(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	json.NewEncoder(w).Encode(DAL.ReadSingleTwoot(idParam))
}

func storeTwoot(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("Storing Twoot")
	// parse the request body into a Twoot struct
	var twoot Models.Twoot
	err := json.NewDecoder(body).Decode(&twoot)
	fmt.Println(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// insert the twoot into the database
	DAL.InsertTwoot(twoot, w)

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.Use(CORS)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAll)
	myRouter.HandleFunc("/get/{id}", returnTwoot)
	myRouter.HandleFunc("/create", storeTwoot)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}



//service functions

func getAllTwoots() (values []primitive.M) {
	return DAL.ReadAllTwoots()
}

// other
// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		//return
	})

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
