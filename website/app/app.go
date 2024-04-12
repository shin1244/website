package app

import (
	"encoding/json"
	"fmt"
	"mongo/website/mongo"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/*
기본 사이트 접근시
*/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

/*
/users에 접근시 userMap이 비어있으면 No Users를 반환하고
비어있지 않으면 유저들의 정보를 반환함
*/
func userHandler(w http.ResponseWriter, r *http.Request) {
	MongoClient := mongo.MongoOpen()
	defer mongo.MongoClose(MongoClient)
	users := mongo.FindAllUserMongo(MongoClient, "test")
	if len(users) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	data, _ := json.Marshal(users)
	fmt.Fprint(w, string(data))

	w.WriteHeader(http.StatusCreated)
}

/*
 */
func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := new(mongo.User)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return
	}

	MongoClient := mongo.MongoOpen()
	user = mongo.FindOneUserMongo(MongoClient, "test", id)
	defer mongo.MongoClose(MongoClient)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(mongo.User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user.CreatedAt = time.Now()
	MongoClient := mongo.MongoOpen()
	defer mongo.MongoClose(MongoClient)

	mongo.InsertOneUserMongo(MongoClient, "test", user)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(mongo.User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	MongoClient := mongo.MongoOpen()
	defer mongo.MongoClose(MongoClient)
	delBool := mongo.DeleteOneUserMongo(MongoClient, "test", user.Id)

	if delBool {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Deleted Id:", user.Id)
	} else {
		fmt.Fprintf(w, "Not Found User %v\n", user.Id)
	}

}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	updateUser := new(mongo.User)
	err := json.NewDecoder(r.Body).Decode(updateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	MongoClient := mongo.MongoOpen()
	defer mongo.MongoClose(MongoClient)
	mongo.UpdateOneUserMongo(MongoClient, "test", updateUser.Id, updateUser.UserName)

	fmt.Fprintf(w, "User %v, UserName Change: %v\n", updateUser.Id, updateUser.UserName)
}

func NewHandler() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", userHandler).Methods("GET")
	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:.+}", getUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users", deleteUserHandler).Methods("DELETE")
	return mux
}
