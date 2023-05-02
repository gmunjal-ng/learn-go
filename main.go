package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//~* Creating a todo struct

type todo struct{
	ID string `json:"id"`
	Text string `json:"text"`
	Completed bool `json:"completed"`
}

//~* Creating a slice of the todo struct

var todos = []todo {
	{ID: uuid.New().String(), Text: "Buy Milk", Completed: false},
	{ID: uuid.New().String(), Text: "Write a web server in Go", Completed: true},
}

func main(){
	router := gin.Default()

	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodosById)
	router.POST("/todos", postTodos)
	router.PUT("/todos/:id", putTodos)
	router.DELETE("/todos/:id", deleteTodos)

	router.Run("localhost:8080")
}

//~* getTodos function serialize the todos struct into JSON and add it to the response.

func getTodos(c *gin.Context){
	filteredTodos := Filter(todos, func(td todo) bool {
		return strconv.FormatBool(td.Completed) == (c.Query("completed"))
	})
	c.IndentedJSON(http.StatusOK, filteredTodos)
}

func getTodosById(c *gin.Context){
	id := c.Param("id")
	for _, td := range todos{
		if td.ID == id{
			c.IndentedJSON(http.StatusOK, td)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Todo item not found!"} )
}

//~* postTodos function de-serialize the JSON request body into the todo struct, add it to the todos slice

func postTodos(c *gin.Context){
	var newTodo todo
	newTodo.ID = uuid.New().String()
	if err:= c.BindJSON(&newTodo); err != nil{
		return
	}
	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

//~* Update toggles the completed value for the todo with specified ID

func putTodos(c *gin.Context){
	id := c.Param("id")
	for i, td := range todos{
		if td.ID == id{
			todos[i].Completed = !td.Completed
			c.IndentedJSON(http.StatusOK, todos[i])
			return
		}
	}
}

//~* Deleted the value with specified ID
func deleteTodos(c *gin.Context){
	id := c.Param("id")
	for i, td := range todos{
		if td.ID == id{
			todos = remove(todos, i)
			c.IndentedJSON(http.StatusOK, td)
			return
		}
	}
}

//~* Helper Functions

// Filter returns a new slice holding only the elements of s that satisfy fn()
func Filter(s []todo, fn func(todo) bool) []todo {
    var p []todo // == nil
    for _, v := range s {
        if fn(v) {
            p = append(p, v)
        }
    }
	if p == nil {return todos} else {return p}
}

func remove(slice []todo, i int) []todo {
    copy(slice[i:], slice[i+1:])
    return slice[:len(slice)-1]
}
