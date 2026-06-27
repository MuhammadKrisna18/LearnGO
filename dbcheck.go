package main

import (
	"fmt"
	"log"

	"modular-monolith/config"
	"modular-monolith/internal/shared/database"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	type Kelas struct {
		ID             string
		Name           string
		Hari           string
		JamMulai       string
		JamSelesai     string
	}

	var k Kelas
	err = db.Raw("SELECT id, name, hari, jam_mulai, jam_selesai FROM kelas LIMIT 1").Scan(&k).Error
	if err != nil {
		fmt.Println("DB ERROR:", err)
	} else {
		fmt.Printf("Row: %+v\n", k)
	}
}
