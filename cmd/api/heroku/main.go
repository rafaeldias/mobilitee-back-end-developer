package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/pkg/api"
	"github.com/rafaeldias/mobilitee-back-end-developer/pkg/device"
)

func main() {
	port := os.Getenv("PORT")
	dbConn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	router := httprouter.New()

	api.RestfulDevice(router, device.New(db))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		panic(err)
	}
}
