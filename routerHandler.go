package main

import(
	"net/http"
)


func routerHandler(w http.ResponseWriter,r *http.Request){
	respondWithJSON(w,200,struct{}{})
}