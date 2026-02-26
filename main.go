package main
 import (
	"log"
	"os"
  "time"
	"net/http"
  "database/sql"
  db "github.com/eyop23/rssagg/internal/database"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
  _ "github.com/lib/pq"


 )

 type apiConfig struct {
  DB *db.Queries
 }
func main(){
  godotenv.Load()
  Port:=os.Getenv("PORT")
  if Port == ""{
    log.Fatal("Port is missing")
  }
  db_url:=os.Getenv("DB_URL")
  if db_url == ""{
    log.Fatal("db url is missing")
  }
  conn,err := sql.Open("postgres",db_url)
  if err != nil{
    log.Fatal("can't connect to database:",err)
  }
  if err = conn.Ping(); err != nil {
    log.Fatal("database not reachable:", err)
  }
  db:=db.New(conn)
  apiCfg:=apiConfig{
    DB:db,
  }

  go startScraping(db,10,time.Minute)

  router := chi.NewRouter()

  router.Use(cors.Handler(cors.Options{
	AllowedOrigins: []string{"http://*","https://*"},
	AllowedMethods: []string{"GET","POST","DELETE","PATCH","PUT"},
	AllowedHeaders: []string{"*"},
	AllowCredentials: false,
	MaxAge: 200,
  }))

  v1Router:=chi.NewRouter();

  v1Router.Get("/ready",routerHandler)
  v1Router.Get("/err",handlerErr)
  v1Router.Post("/users",apiCfg.handlerCreateUser)
  v1Router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))
  v1Router.Get("/posts",apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

  v1Router.Post("/feeds",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
  v1Router.Get("/feeds",apiCfg.handlerGetFeeds)
  v1Router.Post("/feed_follows",apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
  v1Router.Get("/feed_follows",apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
  v1Router.Delete("/feed_follows/{feed_follow_id}",apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

  router.Mount("/v1",v1Router)
  

  server := &http.Server{
	Handler : router,
	Addr:":" + Port,
  }
  log.Printf("server starting on port %v",Port)
  err = server.ListenAndServe();
  if err != nil {
	log.Fatal(err)
  }
}