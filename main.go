package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type User struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Phone  string `json:"phone,omitempty"`
}

func main() {
	db, err := sql.Open("postgres", "user=postgres password=123123123a dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	users, err := query(db)
	if err != nil {
		log.Fatal(err)
	}

	output := struct {
		StatusCode int           `json:"status_code"`
		Data       []interface{} `json:"data"`
	}{
		StatusCode: 200,
		Data:       users,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}

func query(db *sql.DB) ([]interface{}, error) {
	rows, err := db.Query("SELECT user_id, name, age, phone FROM public.user_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []interface{}{}
	for rows.Next() {
		var u User
		var phone sql.NullString
		err := rows.Scan(&u.UserID, &u.Name, &u.Age, &phone)
		if err != nil {
			return nil, err
		}

		if phone.Valid {
			phoneWithoutNewLine := strings.Replace(phone.String, "\n", "", -1)
			users = append(users, User{
				UserID: u.UserID,
				Name:   u.Name,
				Age:    u.Age,
				Phone:  phoneWithoutNewLine,
			})
		} else {
			users = append(users, struct {
				UserID int    `json:"user_id"`
				Name   string `json:"name"`
				Age    int    `json:"age"`
			}{
				UserID: u.UserID,
				Name:   u.Name,
				Age:    u.Age,
			})
		}
	}
	return users, nil
}
