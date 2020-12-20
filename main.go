package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/filipwtf/url-longer/postgres"
	"github.com/filipwtf/url-longer/server"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	dev, err := strconv.ParseBool(os.Getenv("dev"))
	if err != nil {
		log.Fatal(err)
	}
	host := os.Getenv("host")
	port, err := strconv.ParseInt(os.Getenv("port"), 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	dbUsername := os.Getenv("username")
	dbPassword := os.Getenv("password")
	dbName := os.Getenv("name")

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, dbUsername, dbPassword, dbName)

	conn, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	db := postgres.New(conn)
	color.Set(color.FgGreen)
	log.Println("App is starting 😀")
	color.Unset()
	srv := server.NewServer(db, dev)
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), srv)
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalf("error starting server %s\n", err)
		color.Unset()
	}
}
