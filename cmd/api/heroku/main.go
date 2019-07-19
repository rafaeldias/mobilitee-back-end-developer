package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/pkg/device"
	"rafaeldias/mobilitee-back-end-developer/cmd/api"
)

func main() {
	port := os.Getenv("PORT")
	dbConn := os.getenv("DATABASE_URL")

	db, err := gorm.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}

	api.RestfulDevice(httprouter.New(), device.New(db))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}
