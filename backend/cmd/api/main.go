package main

import (
	"fmt"
	"log"
	"time"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main(){
	host := "localhost"
	port := 5432
	user := "admin"
	password := "lemmaballs"
	dbname := "opentrace_map"

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Println("Connecting to the database")

	var db *sqlx.DB
	var err error

	for i:=0; i<5; i++{
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil{
			break
		}
		fmt.Println("DB is not ready yet, retrying in 2 seconds")
		time.Sleep(2*time.Second)

	}

	if err != nil{
		log.Fatalf("Failed to connect after retries: %v", err)
	}

	err= db.Ping()
	if err != nil{
		log.Fatal("Database is connected by unreachable:", err)
	}

	fmt.Println("Success! Go is now talking to the PostGIS container.")

	defer db.Close()
}
