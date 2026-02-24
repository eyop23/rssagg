package auth

import 
(
	"net/http"
	"errors"
	"strings"
)

func GetAPIKey(headers http.Header) (string,error){
	value := headers.Get("Authorization")

	if value == "" {
       return "",errors.New("No authentication info found")
	}
 vals:= strings.Split(value," ")
 if len(vals) != 2{
	return "",errors.New("malformed auth")
 }
 if vals[0] != "ApiKey"{
	return "",errors.New("malformed first part of auth")
 }
 return vals[1],nil
}