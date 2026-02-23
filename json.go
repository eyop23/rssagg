package main

import (
	"net/http"
	"log"
	"encoding/json"

)
func respondWithError(w http.ResponseWriter,code int, msg string){
if code > 499 {
	log.Println("error with 5XX",msg)
}
type errorResponse struct {
	Error string `json:"error"`
}
respondWithJSON(w,code,errorResponse{
	Error:msg,
})

}

func respondWithJSON(w http.ResponseWriter,code int, payload interface{}){
  data,err := json.Marshal(payload)
  if err != nil {
	log.Printf("failed to marshal JSON response: %v",err)
	w.WriteHeader(500)
	return
  }
  w.Header().Add("Content-Type","application/json")
  w.WriteHeader(code)
  w.Write(data)
}