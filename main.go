package main

import(
	"net/http"
	"log"
)

func main(){
	

	ConnectRedis()
	go SubscribeToRedis()
	Routes()



	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080",nil)
}