package main

import (
	"encoding/json"
	"fmt"
	"github.com/mailru/dbr"
	_ "github.com/mailru/go-clickhouse"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Rotation struct {
	Name string `db:"name" json:"name"`
	Duration int `db:"duration" json:"duration"`
	Memory int `db:"memory" json:"memory"`
	Origin float64 `db:"origin" json:"origin"`
	StartTime int `db:"start_time" json:"start_time"`
	EndTime int `db:"end_time" json:"end_time"`
}

var connection *dbr.Connection

func rotationListHandler(w http.ResponseWriter, r *http.Request) {
	var items []Rotation

	sess := connection.NewSession(nil)
	query := sess.Select("name", "duration", "memory", "origin", "start_time", "end_time").From("rotation")
	//query.Where(dbr.Eq("country_code", "RU"))

	if _, err := query.Load(&items); err != nil {
		log.Fatal(err)
	}

	res, _ := json.Marshal(items)
	fmt.Fprint(w, res)
}

func rotationAddHandler(w http.ResponseWriter, r *http.Request) {
	var newRotation Rotation

	sess := connection.NewSession(nil)

	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&newRotation); if errDecode != nil {
		fmt.Printf("Error decode %s", errDecode)
	}

	_, err := sess.InsertInto("rotation").
		Columns("name", "duration", "memory", "origin", "start_time", "end_time").
		Record(newRotation).
		Exec()

	if err != nil {
		log.Fatal(err)
	}

	res, _ := json.Marshal(newRotation)
	fmt.Fprint(w, res)
}

func main() {
	var err error
	connection, err = dbr.Open("clickhouse", "http://127.0.0.1:8123/default", nil)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/rotations", rotationListHandler).Methods("GET")
	router.HandleFunc("/api/v1/rotations", rotationAddHandler).Methods("POST")
	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8000", nil)
}