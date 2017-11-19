package main

import (
	// "fmt"
	"net/http"
	"html/template"
	"time"
)

// Index Template
var t, _ = template.ParseFiles("./templates/index.html")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, nil)
}

func main() {
	router := http.NewServeMux()

	// Serves Static Files
	// router.Handle("/assests/",
	// 	http.StripPrefix("/assests/", http.FileServer(http.Dir("public/"))))

	// Serves index.html
	router.HandleFunc("/", indexHandler)

	// 

	// Server Config
	srvr := http.Server{
		Addr: ":3000",
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	panic(srvr.ListenAndServe())

}
