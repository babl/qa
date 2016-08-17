package http

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func StartHttpServer(listen string,
	HandlerRequestHistory func(w http.ResponseWriter, r *http.Request),
	HandlerRequestDetails func(w http.ResponseWriter, r *http.Request)) {
	const dir = "./httpserver/static"
	r := mux.NewRouter()

	// REST API
	r.HandleFunc("/api/request/history", HandlerRequestHistory).Methods("GET").Queries("blocksize", "{blocksize}")
	r.HandleFunc("/api/request/history", HandlerRequestHistory).Methods("GET")
	r.HandleFunc("/api/request/details/{requestid:[0-9]+}", HandlerRequestDetails).Methods("GET")

	// Static files and assets
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler:      r,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
