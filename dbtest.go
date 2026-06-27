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
		ID             string `gorm:"primaryKey;type:varchar(255)"`
		Name           string `gorm:"type:varchar(255)"`
		Capacity       int
		Hari           string `gorm:"type:varchar(20);not null"`
		JamMulai       string `gorm:"type:varchar(10);not null"`
		JamSelesai     string `gorm:"type:varchar(10);not null"`
		ProgramStudiID string `gorm:"type:varchar(255);not null"`
	}

	k := Kelas{
		ID:             "dummy-id-123",
		Name:           "IF-TEST",
		Capacity:       30,
		Hari:           "Senin",
		JamMulai:       "07:00",
		JamSelesai:     "09:00",
		ProgramStudiID: "0c874dfd-1959-4507-b35f-14909a3dc377", 
	}

	err = db.Create(&k).Error
	if err != nil {
		fmt.Println("INSERT ERROR:", err)
	} else {
		fmt.Println("INSERT SUCCESS")
		db.Where("id = ?", "dummy-id-123").Delete(&Kelas{})
	}
}
