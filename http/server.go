package http

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func HandlerRequestHistory(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Request History!\n"))
}

func HandlerRequestDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte("Request Details! RequestId=" + vars["requestid"] + "\n"))
}

func StartHttpServer(listen string) {
	const dir = "./http/static"
	r := mux.NewRouter()

	// REST API
	r.HandleFunc("/api/request/logs", HandlerRequestHistory).Methods("GET")
	r.HandleFunc("/api/request/details/{requestid:[0-9]+}", HandlerRequestDetails).Methods("GET")

	// Static files and assets
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler: r,
		Addr:    listen,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
