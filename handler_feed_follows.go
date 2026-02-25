package main

import (
	"net/http"
	"time"
	"fmt"

	"encoding/json"
	"github.com/eyop23/rssagg/internal/database"


	"github.com/google/uuid"
	"github.com/go-chi/chi"

)

func(config *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter,r *http.Request,user db.User){

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder:=json.NewDecoder(r.Body)
	params:=parameters{};
	err:=decoder.Decode(&params)
	if err != nil {
		respondWithJSON(w,400,fmt.Sprintf("error parsing json %v",err))
		return
	}
    feed_follows,err:=config.DB.CreateFeedFollow(r.Context(),db.CreateFeedFollowParams{
		ID:uuid.New(),
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC(),
		UserID:user.ID,
		FeedID:params.FeedID,
	})
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("failed to create feed follow %v",err))
	}
	respondWithJSON(w,201,feed_follows)
}

func(config *apiConfig) handlerGetFeedFollow(w http.ResponseWriter,r *http.Request,user db.User){

    feed_follows,err:=config.DB.GetFeedFollows(r.Context(),user.ID)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("failed to get feed follow %v",err))
	}
	respondWithJSON(w,201,feed_follows)
}
func(config *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter,r *http.Request,user db.User){

	idString:=chi.URLParam(r,"feed_follow_id")
	feedFollowId,err := uuid.Parse(idString);

	if err != nil {
		respondWithError(w,400,fmt.Sprintf("failed to parse feed follow id %v",err))
		return
	}

    err = config.DB.DeleteFeedFollow(r.Context(),db.DeleteFeedFollowParams{
		ID:feedFollowId,
		UserID:user.ID,
	})
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("failed to remove feed follow %v",err))
	}
	respondWithJSON(w,200,struct{}{})
}






