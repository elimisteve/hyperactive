// Steve Phillips / elimisteve
// 2013.11.10

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

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
		writeError(w, err)
		return
	}
	// Marshall all HypeService ~objects to JSON
	jsonStr, err := json.Marshal(services)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonStr)
}

func PostServices(w http.ResponseWriter, r *http.Request) {
	// Read POSTed body (should be JSON)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err)
		return
	}
	defer r.Body.Close()

	hs := &types.HypeService{}
	// Unmarshal JSON into TentServer var
	if err := json.Unmarshal(body, hs); err != nil {
		writeError(w, err)
		return
	}

	// Store to DB
	if err = hs.Save(); err != nil {
		writeError(w, err)
		return
	}

	// Marshal to JSON and return to user
	jsonData, err := json.Marshal(hs)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonData)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, err.Error(), 500)
}
