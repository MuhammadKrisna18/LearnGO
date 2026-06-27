package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Login
	loginBody := map[string]string{
		"email":    "admin@univ.edu",
		"password": "password123",
	}
	b, _ := json.Marshal(loginBody)
	resp, err := http.Post("http://localhost:8080/api/v1/auth/admin/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Login err:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Login resp:", string(body))

	var loginData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(body, &loginData)
	token := loginData.Data.Token

	if token == "" {
		fmt.Println("No token!")
		return
	}

	// Fetch program studi
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/program-studi", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp2, _ := client.Do(req)
	body2, _ := ioutil.ReadAll(resp2.Body)
	
	var prodiData struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(body2, &prodiData)
	if len(prodiData.Data) == 0 {
		fmt.Println("No prodi found")
		return
	}
	prodiId := prodiData.Data[0].ID

	// Create Kelas
	kelasBody := map[string]interface{}{
		"name":             "IF-999",
		"capacity":         30,
		"hari":             "Senin",
		"jam_mulai":        "07:00",
		"jam_selesai":      "09:00",
		"program_studi_id": prodiId,
	}
	b3, _ := json.Marshal(kelasBody)
	req3, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/kelas", bytes.NewBuffer(b3))
	req3.Header.Set("Authorization", "Bearer "+token)
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := client.Do(req3)
	body3, _ := ioutil.ReadAll(resp3.Body)
	fmt.Println("Create Kelas Status:", resp3.StatusCode)
	fmt.Println("Create Kelas Resp:", string(body3))
}
