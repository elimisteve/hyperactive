// Steve Phillips / elimisteve
// 2013.11.10

package main

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/222Labs/help"
	"github.com/gorilla/mux"

	"./types"
)

func main() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	router := mux.NewRouter()

	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/services", GetServices).Methods("GET")
	router.HandleFunc("/services/new", PostServices).Methods("POST")

	http.Handle("/", router)

	// Start HTTP server
	server := SimpleHTTPServer(router, ":9999")
	log.Printf("HTTP server trying to listen on %v...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP listen failed: %v\n", err)
	}
}

func SimpleHTTPServer(handler http.Handler, host string) *http.Server {
	return &http.Server{
		Addr:           host,
		Handler:        handler,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Hyperactive! Check out /services\n"))
}

func GetServices(w http.ResponseWriter, r *http.Request) {
	// Grab all HypeService objects from DB
	services, err := types.ServicesList()
	if err != nil {
		help.WriteError(w, err.Error(), 500)
		return
	}

	help.WriteJSON(w, services)
}

func PostServices(w http.ResponseWriter, r *http.Request) {
	hs := &types.HypeService{}
	err := help.ReadInto(r.Body, hs)
	if err != nil {
		help.WriteError(w, err.Error(), 400)
		return
	}

	if err = hs.Validate(); err != nil {
		help.WriteError(w, err.Error(), 400)
		return
	}

	// Store to DB
	if err = hs.Save(); err != nil {
		help.WriteError(w, err.Error(), 500)
		return
	}

	help.WriteJSON(w, hs)
}
