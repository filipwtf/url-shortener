package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/filipwtf/url-longer/postgres"
	"github.com/filipwtf/url-longer/server"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	dev := flag.Bool("dev", true, "hides dev routes")
	host := flag.String("host", "localhost", "speicify postgres host")
	port := flag.Int("port", 5432, "specify postgres port")
	dbUsername := flag.String("username", "postgres", "db user username")
	dbPassword := flag.String("password", "postgres", "db user password")
	dbName := flag.String("name", "longer", "database name")
	flag.Parse()

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", *host, *port, *dbUsername, *dbPassword, *dbName)

	conn, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	db := postgres.New(conn)
	color.Set(color.FgGreen)
	log.Println("App is starting 😀")
	color.Unset()
	srv := server.NewServer(db, *dev)
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), srv)
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalf("error starting server %s\n", err)
		color.Unset()
	}
}
