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

		d, err := InitLxdInstanceServer()
		if err != nil {
			logger.Error(err.Error())
		}

		socket := make(chan string)
		vars := mux.Vars(r)
		go vga(*d, vars["name"], socket)
		spice_socket := <-socket

		w.WriteHeader(302)
		w.Write([]byte(spice_socket))
	})

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
