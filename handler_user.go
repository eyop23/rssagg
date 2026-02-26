package main

import (
	"net/http"
	"time"
	"fmt"
	"encoding/json"
	db "github.com/eyop23/rssagg/internal/database"


	"github.com/google/uuid"
)

func(config *apiConfig) handlerCreateUser(w http.ResponseWriter,r *http.Request){

	type parameters struct {
		Name string `json:"name"`
	}
	decoder:=json.NewDecoder(r.Body)
	params:=parameters{};
	err:=decoder.Decode(&params)
	if err != nil {
		respondWithJSON(w,400,fmt.Sprintf("error parsing json %v",err))
		return
	}
    user,err:=config.DB.CreateUser(r.Context(),db.CreateUserParams{
		ID:uuid.New(),
		Name:params.Name,
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC(),
	})
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("failed to create user %v",err))
	}
	respondWithJSON(w,201,user)
}

func (config *apiConfig) handlerGetUser(w http.ResponseWriter,r *http.Request,user db.User){
	
	respondWithJSON(w,200,user)
}

func (config *apiConfig) handlerGetPostsForUser(w http.ResponseWriter,r *http.Request,user db.User){
	 posts,err := config.DB.GetPostsForUser(r.Context(),db.GetPostsForUser{
		UserID:user.ID,
		Limit:10
	 })

	 if err != nil {
		respondWithJSON(w,400,fmt.Sprintf("error finding posts for user %v",err))
		return
	}		
	respondWithJSON(w,200,posts)
}
