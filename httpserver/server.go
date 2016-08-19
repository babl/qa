package httpserver

import (
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	. "github.com/larskluge/babl-server/utils"
)

func GetVarsBlockSize(r *http.Request, defaultvalue int64) int64 {
	result := defaultvalue
	vars := mux.Vars(r)
	blocksize := vars["blocksize"]
	if blocksize != "" {
		bsize, errParse := strconv.ParseInt(blocksize, 10, 64)
		if errParse == nil {
			result = bsize
		}
	}
	return result
}

func StartHttpServer(listen string,
	HandlerRequestHistory func(w http.ResponseWriter, r *http.Request),
	HandlerRequestDetails func(w http.ResponseWriter, r *http.Request)) {
	pwd, err := os.Getwd()
	Check(err)
	dir := pwd + "/httpserver/static"
	//fmt.Println("WorkingDir: ", pwd)
	//fmt.Println("HttpServer: ", dir)
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
