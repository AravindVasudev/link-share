package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Index Template
var t, _ = template.ParseFiles("./templates/index.html")

// Mongo Conn
var dbCon *mgo.Session
var coll *mgo.Collection

// Add URL JSON Struct
type addURL struct {
	URL   string `json:"url"`
	Group string `json:"group"`
}

// URL Document
type urlDoc struct {
	Group string   `json:"group"`
	URLs  []string `json:"urls"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbURL := os.Getenv("DB_URL")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var err error
	dbCon, err = mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	defer dbCon.Close()

	coll = dbCon.DB("links").C("urls")

	router := http.NewServeMux()

	// Route Handlers
	router.Handle("/assests/",
		http.StripPrefix("/assests", http.FileServer(http.Dir("assests/"))))

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/add", addHandler)
	router.HandleFunc("/urls", urlsHandler)

	// Server Setup
	server := http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	panic(server.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, nil)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	// read body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	// Unmarshal
	var url addURL
	err = json.Unmarshal(b, &url)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Upsert url
	_, err = coll.Upsert(
		bson.M{"group": url.Group},
		bson.M{"$push": bson.M{"urls": url.URL}},
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "1")
}

func urlsHandler(w http.ResponseWriter, r *http.Request) {
	var output []urlDoc

	err := coll.Find(nil).All(&output)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(data))
}
