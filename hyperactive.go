// Steve Phillips / elimisteve
// 2013.11.10

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/222Labs/help"
	"github.com/gorilla/mux"

	"./types"
)

func init() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Dump JSON when this process is killed
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)
	go func() {
		log.Printf("Got this signal: %v\n", <-stop)
		err := types.DumpDB()
		if err != nil {
			log.Printf("Error from DumpDB: %v\n", err)
		}
		os.Exit(0)
	}()
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/services", GetServices).Methods("GET")
	router.HandleFunc("/services", UpdateServices).Methods("PUT")
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

	hs.CreatedBy = r.RemoteAddr
	hs.ModifiedBy = r.RemoteAddr

	// Store to DB
	if err = hs.Save(); err != nil {
		statusCode := 500
		if err == types.ErrServiceDuplicate {
			statusCode = 400
		}
		help.WriteError(w, err.Error(), statusCode)
		return
	}

	help.WriteJSON(w, hs)
}

func UpdateServices(w http.ResponseWriter, r *http.Request) {
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

	hs.ModifiedBy = r.RemoteAddr

	// Update in DB
	if err = hs.Update(); err != nil {
		statusCode := 500
		if err == types.ErrServiceNotFound {
			statusCode = 404
		}
		help.WriteError(w, err.Error(), statusCode)
		return
	}

	help.WriteJSON(w, hs)
}
