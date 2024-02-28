package main

import(
	"net/http"
)

func Routes(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/index.html")
    })

	http.HandleFunc("/ws",wsHandler)
}