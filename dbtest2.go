package main

import (
	"fmt"
	"log"

	"modular-monolith/config"
	"modular-monolith/internal/modules/programstudi/domain"
	"modular-monolith/internal/shared/database"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var prodi domain.ProgramStudi
	err = db.First(&prodi).Error
	if err != nil {
		fmt.Println("No prodi found:", err)
		return
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
		ID:             "dummy-id-124",
		Name:           "IF-TEST2",
		Capacity:       30,
		Hari:           "Senin",
		JamMulai:       "07:00",
		JamSelesai:     "09:00",
		ProgramStudiID: prodi.ID, 
	}

	err = db.Create(&k).Error
	if err != nil {
		fmt.Println("INSERT ERROR:", err)
	} else {
		fmt.Println("INSERT SUCCESS")
		db.Where("id = ?", "dummy-id-124").Delete(&Kelas{})
	}
}
