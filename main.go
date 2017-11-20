package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {
	var err error
	dbCon, err = mgo.Dial("127.0.0.1:27017")
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
		Addr:         "127.0.0.1:3000",
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
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "1")
}

func urlsHandler(w http.ResponseWriter, r *http.Request) {
	link := addURL{}
	err := coll.Find(bson.M{}).One(&link)
	if err != nil {
		// panic(err)
		log.Fatal(err)
	}

	fmt.Println(link)

	log.Println("urls")
	fmt.Fprintf(w, "hello")
}
