package main

import (
	"log"
	"time"

	"github.com/JaskiratAnand/go-social/internal/db"
	"github.com/JaskiratAnand/go-social/internal/env"
	"github.com/JaskiratAnand/go-social/internal/store"
)

func main() {
	start := time.Now()

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

	log.Println("running database seed script...")
	err = db.Seed(store)
	if err != nil {
		log.Println("error: ", err)
	}

	duration := time.Since(start)
	log.Printf("Database seeding completed (%vms)\n", duration.Milliseconds())
}
