package server

import (
	"Rest-api-module/internal/config"
	"Rest-api-module/internal/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var cfg = config.GetDbConfig()

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

type Server struct {
}

func (s *Server) Start() {
	r := mux.NewRouter()

	r.HandleFunc("/api/contacts", GetContacts).Methods("GET")
	r.HandleFunc("/api/contacts/{id}", GetContactByUuid).Methods("GET")
	r.HandleFunc("/api/contacts", CreateContact).Methods("POST")
	r.HandleFunc("/api/contacts/{id}", UpdateContact).Methods("PUT")
	// r.HandleFunc("/api/contacts/{id}", ParticallyUpdateContact).Methods("PATCH")
	r.HandleFunc("/api/contacts/{id}", RemoveContact).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	db, _ := sql.Open("postgres", psqlInfo)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM "Contact"`)
	if err != nil {
		panic(err)
	}
	contacts := []domain.Contact{}
	for rows.Next() {
		p := domain.Contact{}
		if err := rows.Scan(&p.Id, &p.Name, &p.Email); err != nil {
			fmt.Println(err)
			continue
		}
		contacts = append(contacts, p)
	}
	json.NewEncoder(w).Encode(contacts)
}

func GetContactByUuid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	db, _ := sql.Open("postgres", psqlInfo)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	row := db.QueryRow(`Select * from "Contact" where "Id" = $1`, params["id"])
	contacts := []domain.Contact{}
	p := domain.Contact{}
	if err := row.Scan(&p.Id, &p.Name, &p.Email); err != nil {
		panic(err)
	}
	contacts = append(contacts, p)

	json.NewEncoder(w).Encode(contacts)

}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	cont := domain.Contact{}
	json.NewDecoder(body).Decode(&cont)

	db, _ := sql.Open("postgres", psqlInfo)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := db.Exec(`insert into "Contact" ("Name","Email") values ($1,$2)`, cont.Name, cont.Email); err != nil {
		panic(err)
	}

}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("postgres", psqlInfo)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)

	body := r.Body
	defer body.Close()

	cont := domain.Contact{}
	json.NewDecoder(body).Decode(&cont)

	_, err := db.Exec(`update "Contact" set "Name" = $1,"Email" = $2 where "Id" = $3`, cont.Name, cont.Email, params["id"])
	if err != nil {
		panic(err)
	}
}

func RemoveContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db, _ := sql.Open("postgres", psqlInfo)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`delete from "Contact" where "Id" = $1`, params["id"])
}
