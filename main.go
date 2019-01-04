package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mailru/go-clickhouse"
	"net/http"
)

type Rotation struct {
	Name string `db:"name" json:"name"`
	Duration int `db:"duration" json:"duration"`
	Memory int `db:"memory" json:"memory"`
	Origin float64 `db:"origin" json:"origin"`
	StartTime int `db:"start_time" json:"start_time"`
	EndTime int `db:"end_time" json:"end_time"`
}

var connection *gorm.DB

func rotationListHandler(w http.ResponseWriter, r *http.Request) {
	var items []Rotation
	connection.Find(&items)

	res, _ := json.Marshal(items)
	w.Write([]byte(res))
}

func rotationAddHandler(w http.ResponseWriter, r *http.Request) {
	var newRotation Rotation
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&newRotation); if errDecode != nil {
		fmt.Printf("Error decode %s", errDecode)
	}

	connection.Create(&newRotation)

	res, _ := json.Marshal(newRotation)
	w.Write([]byte(res))
}

func main() {
	var err error
	connection, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer connection.Close()

	connection.AutoMigrate(&Rotation{})

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/rotations", rotationListHandler).Methods("GET")
	router.HandleFunc("/api/v1/rotations", rotationAddHandler).Methods("POST")
	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8000", nil)
}