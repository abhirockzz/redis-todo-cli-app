package db

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

const redisHostEnvVar = "REDIS_HOST"
const redisPasswordEnvVar = "REDIS_PASSWORD"
const redisPasswordRequiredEnvVar = "REDIS_PASSWORD_REQUIRED"
const sslRequiredEnvVar = "REDIS_SSL_REQUIRED"
const defaultRedisHost = "localhost:6379"

var redisHost string
var redisPassword string

// true or false
var redisPasswordRequired = "true"

// true or false
var sslRequired = "true"

func init() {
	redisHost = os.Getenv(redisHostEnvVar)
	if redisHost == "" {
		redisHost = defaultRedisHost
	}

	redisPasswordRequired = os.Getenv(redisPasswordRequiredEnvVar)
	if redisPasswordRequired == "" {
		redisPasswordRequired = "false"
	}

	if redisPasswordRequired == "true" {
		redisPassword = os.Getenv(redisPasswordEnvVar)
		if redisPassword == "" {
			log.Fatal("your redis instance requires a password. please provide set environment variable ", redisPasswordEnvVar)
		}
	}

	sslRequired = os.Getenv(sslRequiredEnvVar)
	if sslRequired == "" {
		sslRequired = "false"
	}
}

const todoIDCounter = "todoid"
const todoIDsSet = "todos-id-set"
const statusPending = "pending"

func getClient() *redis.Client {
	opts := &redis.Options{Addr: redisHost}

	if redisPasswordRequired == "true" {
		opts.Password = redisPassword
	}
	if sslRequired == "true" {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	c := redis.NewClient(opts)
	err := c.Ping().Err()
	if err != nil {
		log.Fatal("redis connect failed", err)
	}
	return c
}

// CreateTodo creates a todo in Redis
func CreateTodo(desc string) {
	c := getClient()
	defer c.Close()

	//increment todo id counter
	id, err := c.Incr(todoIDCounter).Result()
	if err != nil {
		log.Fatal("failed to increment todo id counter", err)
	}
	todoid := "todo:" + strconv.Itoa(int(id))

	//store ID in a SET for other operations
	err = c.SAdd(todoIDsSet, todoid).Err()
	if err != nil {
		log.Fatal("failed to add todo id to SET", err)
	}

	//save todo in a HASH
	todo := map[string]interface{}{"desc": desc, "status": statusPending}
	err = c.HMSet(todoid, todo).Err()
	if err != nil {
		log.Fatal("failed to save todo")
	}
	fmt.Println("created todo! use 'todo list' to get your todos")
}

// ListTodos gets todos from Redis based on status (if not, gets all)
func ListTodos(status string) []Todo {
	c := getClient()
	defer c.Close()

	todoHashNames, err := c.SMembers(todoIDsSet).Result()
	if err != nil {
		log.Fatal("failed to get todo IDs", err)
	}

	todos := []Todo{}
	for _, todoHashName := range todoHashNames {
		id := strings.Split(todoHashName, ":")[1]

		todoMap, err := c.HGetAll(todoHashName).Result()
		if err != nil {
			log.Fatalf("failed to get todo from %s - %v\n", todoHashName, err)
		}

		var todo Todo
		if status == "" {
			todo = Todo{id, todoMap["desc"], todoMap["status"]}
			todos = append(todos, todo)
		} else {
			if status == todoMap["status"] {
				todo = Todo{id, todoMap["desc"], todoMap["status"]}
				todos = append(todos, todo)
			}
		}
	}
	if len(todos) == 0 {
		fmt.Println("no todos founds")
		return nil
	}
	return todos
}

// DeleteTodo deletes a todo from redis
func DeleteTodo(id string) {
	c := getClient()
	defer c.Close()

	//delete HASH
	n, err := c.Del("todo:" + id).Result()
	if err != nil {
		log.Fatalf("failed to delete todo with id %s - %v\n", id, err)
	}

	//if ID was valid and HASH got deleted
	if n > 0 {
		//delete from SET
		err = c.SRem(todoIDsSet, "todo:"+id).Err()
		if err != nil {
			log.Fatalf("failed to delete todo from SET %s - %v\n", id, err)
		}
		fmt.Println("deleted todo. use 'todo list' to fetch all todos")
	} else {
		fmt.Printf("todo with id %s not found\n", id)
	}
}

// UpdateTodo updates status, description or both
func UpdateTodo(id, desc, status string) {
	c := getClient()
	defer c.Close()

	//confirm whether todo exists
	exists, err := c.SIsMember(todoIDsSet, "todo:"+id).Result()
	if err != nil {
		log.Fatalf("cannot confirm whether todo %s exists %v", id, err)
	}

	if !exists {
		log.Fatalf("todo with id %s does not exist\n", id)
	}
	updatedTodo := map[string]interface{}{}
	if status != "" {
		updatedTodo["status"] = status
	}

	if desc != "" {
		updatedTodo["desc"] = desc
	}
	err = c.HMSet("todo:"+id, updatedTodo).Err()
	if err != nil {
		log.Fatal("failed to update todo id", id)
	}
	fmt.Printf("updated todo id %s. use 'todo list' to fetch all todos\n", id)
}

// Todo holds todo information
type Todo struct {
	ID     string
	Desc   string
	Status string
}
