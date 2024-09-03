package main

import (
	"database/sql"
	"fmt"
	// "fmt"
	"log"
	"net/http"
	"os"

	// "time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"go-server/database"
	"go-server/handlers"

	// "go-server/shared"

	"github.com/joho/godotenv"
)









func main(){


	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	connStr:=os.Getenv("CONN_STR")
	
	
	db,err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("Could not connect to Database:", err)
	}

	// createTable(db)
	// database.AddPost(db,shared.Post{
    //     Title: "Sample Title",
    //     Url:   "http://example.com/sample-url",
    //     Body:  "This is a sample body for testing purposes.",
    //     Date:  time.Now(), // Current time
    // })	
	// delPost(db,1)
	// putAllToPostsToArray(db)
	

	
	
	

	posts := database.PutAllToPostsToArray(db)

	// chi router to router http requests
	router:=chi.NewRouter()	

	// cors config
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	  }))
	

	// handles /
	router.Get("/",handlers.HandleRoot)

	// handles /posts -> this gets hit by the index.html rendered by the root
	router.Get("/posts",handlers.HandleViewPosts(posts))

	router.Get("/loginView",handlers.HandleLoginView)

	// handles /login to auth myself
	router.Post("/login",handlers.HandleLogin)

	router.Get("/adminView",handlers.HandleAdminView)

	// handles /admin to edit posts
	router.Get("/admin",handlers.Admin)

	

	fmt.Println("server started at 8080")
	http.ListenAndServe(":8080", router)




}
