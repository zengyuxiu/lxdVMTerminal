package main

import (
	"github.com/lxc/lxd/shared/logger"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func StartServer() {
	router := mux.NewRouter()

	router.HandleFunc("/instance/{name}", func(w http.ResponseWriter, r *http.Request) {

		d, err := InitLxdInstanceServer("127.0.0.1")
		if err != nil {
			logger.Error(err.Error())
			return
		}
		server := *d
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		socket := make(chan string)
		vars := mux.Vars(r)
		go vga(server, vars["name"], socket)
		spiceSocket := <-socket
		w.WriteHeader(302)
		_, err = w.Write([]byte(spiceSocket))
		if err != nil {
			logger.Error(err.Error())
			return
		}
	})

	srv := &http.Server{
		Handler: router,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
