package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id        string    `json:"id"`
	Password  int       `json:"password"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

/*
mongoDB 클라이언트 연결
*/
func MongoOpen() mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// MongoDB 연결
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return *client
}

/*
mongoDB 클라이언트 종료
*/
func MongoClose(client mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

/*
mongoDB 저장소에 데이터를 하나 저장한다.
1. client인자는 연결된 mongoDB 클라이언트를 받는다.
2. cn인자는 데이터를 저장할 저장소를 받는다.
3. user인자는 저장소에 넣을 데이터를 받는다.
*/
func InsertOneUserMongo(client mongo.Client, cn string, user *User) {
	collection := client.Database("monGo").Collection(cn)
	createUser, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Create User: %s\n", createUser)
}

/*
 */
func FindOneUserMongo(client mongo.Client, cn string, id string) *User {
	filter := bson.D{{"id", id}}
	collection := client.Database("monGo").Collection(cn)
	var findUser *User
	err := collection.FindOne(context.TODO(), filter).Decode(findUser)
	if err != nil {
		log.Fatal(err)
	}
	return findUser
}

func UpdateOneUserMongo(client mongo.Client, cn string, id string, newName string) {
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"username", newName}}}}
	collection := client.Database("monGo").Collection(cn)
	updateUser, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if updateUser.MatchedCount == 0 {
		fmt.Printf("Not Found User %v\n", id)
		return
	}
	fmt.Printf("User %s, UserName Change: %s\n", id, newName)
}

func DeleteOneUserMongo(client mongo.Client, cn string, id string) bool {
	filter := bson.D{{"id", id}}
	collection := client.Database("monGo").Collection(cn)
	deleteUser, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	if deleteUser.DeletedCount == 0 {
		fmt.Printf("Not Found User %v\n", id)
		return false
	}
	fmt.Printf("Delete User %s\n", id)
	return true
}
func FindAllUserMongo(client mongo.Client, cn string) []*User {
	var users []*User
	filter := bson.D{{}}
	collection := client.Database("monGo").Collection(cn)
	findUser, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for findUser.Next(context.TODO()) {
		var user User
		if err := findUser.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, &user)
	}
	return users
}
