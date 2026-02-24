package main

import (
	"net/http"
	"fmt"
	db "github.com/eyop23/rssagg/internal/database"
	"github.com/eyop23/rssagg/internal/auth"

)

type authedHandler func(http.ResponseWriter,*http.Request,db.User)

func (config *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter,r *http.Request) {
		apiKey,err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w,403,fmt.Sprintf("Auth error %v",err))
		return
	}
	user,err := config.DB.GetUserByAPIKey(r.Context(),apiKey)

	if err != nil {
		respondWithError(w,404,fmt.Sprintf("no user found %v",err))
		return
	}
	handler(w,r,user)
	}
}