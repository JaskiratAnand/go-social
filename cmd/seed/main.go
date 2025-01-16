package main

import (
	"log"

	"github.com/JaskiratAnand/go-social/internal/db"
	"github.com/JaskiratAnand/go-social/internal/env"
	"github.com/JaskiratAnand/go-social/internal/store"
)

func main() {
	conn, err := db.New(
		env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/go-social?sslmode=disable"),
		env.GetInt("DB_MAX_OPEN_CONNS", 30),
		env.GetInt("DB_MAX_IDLE_CONNS", 30),
		env.GetString("DB_MAX_IDLE_TIME", "15m"),
	)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	store := store.New(conn)

	err = db.Seed(store)
	if err != nil {
		log.Println("error: ", err)
	}
}
