package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type people struct {
	Id    int
	First string
	Last  string
}

var id int
var firstname string
var lastname string

func StartUpDatabase() *sql.DB {
	database, _ := sql.Open("sqlite3", "./testsql.db")
	return database
}

func CreateDb() {
	database := StartUpDatabase()
	defer database.Close()
	QueryCreate := "CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT,lastname TEXT );"
	statement, _ := database.Prepare(QueryCreate)
	statement.Exec()
}
func InsertDb(firstname string, lastname string) {
	database := StartUpDatabase()
	defer database.Close()
	QueryInsert := "INSERT INTO people (firstname, lastname) VALUES (?,?)"
	statement, _ := database.Prepare(QueryInsert)
	statement.Exec(firstname, lastname)
}
func SelectDb() []people {
	data := make([]people, 0)
	database := StartUpDatabase()
	defer database.Close()
	QuerySelect := "SELECT id, firstname, lastname FROM people;"
	rows, _ := database.Query(QuerySelect)
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		obj := people{Id: id, First: firstname, Last: lastname}
		data = append(data, obj)
	}
	return data
}
func SelectDbOneRow(idRecord string) (people, error) {
	database := StartUpDatabase()
	defer database.Close()
	QueryOneRow := "SELECT id, firstname, lastname FROM people WHERE id = ? ;"
	err := database.QueryRow(QueryOneRow, idRecord).Scan(&id, &firstname, &lastname)
	if err != nil {
		return people{}, err
	}
	return people{Id: id, First: firstname, Last: lastname}, nil
}

// server requests
func getPeople(context *gin.Context) {
	data := SelectDb()
	context.IndentedJSON(http.StatusOK, data)
}
func getPerson(context *gin.Context) {
	id := context.Param("id")
	data, err := SelectDbOneRow(id)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": "No result found Person doesn't exist !!"})
	} else {
		context.IndentedJSON(http.StatusOK, data)
	}
}
func addPerson(context *gin.Context) {
	var newPerson people
	if err := context.BindJSON(&newPerson); err != nil {
		return
	}
	firstname, lastname = newPerson.First, newPerson.Last
	InsertDb(firstname, lastname)
	context.IndentedJSON(http.StatusCreated, newPerson)
}

func runServer() {
	server := gin.Default()
	server.GET("/people", getPeople)
	server.GET("/people/:id", getPerson)
	server.POST("/people/create", addPerson)
	server.Run(":8099")
}
func main() {
	runServer()
}
