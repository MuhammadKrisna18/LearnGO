package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MataKuliah struct {
	ID             string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string `json:"name"`
	SKS            int    `json:"sks"`
	ProgramStudiID string `json:"program_studi_id"`
}

type PengajuanMataKuliah struct {
	ID           string      `json:"id" gorm:"primaryKey;type:varchar(255)"`
	DosenID      string      `json:"dosen_id"`
	MataKuliahID string      `json:"mata_kuliah_id"`
	MataKuliah   *MataKuliah `json:"mata_kuliah,omitempty" gorm:"foreignKey:MataKuliahID"`
	Status       string      `json:"status"`
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	var list []PengajuanMataKuliah
	err = db.Preload("MataKuliah").Where("status = 'approved'").Order("created_at desc").Find(&list).Error
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	fmt.Println("=== MY-REQUESTS RESPONSE (simulated JSON) ===")
	data, _ := json.MarshalIndent(list, "", "  ")
	fmt.Println(string(data))
}
