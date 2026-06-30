package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"siakad-pro/internal/modules/kelas/domain"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=siakad_pro port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	var absensiList []domain.Absensi
	db.WithContext(ctx).Find(&absensiList)
	fmt.Println("Total absensi:", len(absensiList))
	for _, a := range absensiList {
		fmt.Printf("Absensi: ID=%s, PertemuanID=%s, MhsID=%s, Status=%s\n", a.ID, a.PertemuanID, a.MahasiswaID, a.StatusKehadiran)
	}
}
