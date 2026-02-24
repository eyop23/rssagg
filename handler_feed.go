package main

import (
	"net/http"
	"time"
	"fmt"

	"encoding/json"
	db "github.com/eyop23/rssagg/internal/database"


	"github.com/google/uuid"
)

func(config *apiConfig) handlerCreateFeed(w http.ResponseWriter,r *http.Request,user db.User){

	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url`
	}
	decoder:=json.NewDecoder(r.Body)
	params:=parameters{};
	err:=decoder.Decode(&params)
	if err != nil {
		respondWithJSON(w,400,fmt.Sprintf("error parsing json %v",err))
		return
	}
    feed,err:=config.DB.CreateFeed(r.Context(),db.CreateFeedParams{
		ID:uuid.New(),
		Name:params.Name,
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC(),
		Url:params.Url,
		UserID:user.ID,
	})
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("failed to create user %v",err))
	}
	respondWithJSON(w,201,feed)
}

func(config *apiConfig) handlerGetFeeds(w http.ResponseWriter,r *http.Request){
    feeds,err:=config.DB.GetFeeds(r.Context())
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("failed to load feeds %v",err))
	}
	respondWithJSON(w,201,feeds)
}



