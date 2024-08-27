package main

import (
	"database/sql"
	// "fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"go-server/database"
	"go-server/handlers"
	"go-server/shared"

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

	// post1 := Post{
	// 	ID:    2,
	// 	Title: "Introduction to Go",
	// 	Url:   "https://cpu.land",
	// 	Body:  "Go is an open-source programming language designed for simplicity and efficiency. In this post, we will cover the basics of Go, including its syntax, types, and how to write a simple Go program.",
	// 	Date: time.Now(),
	// }

	//Sample Post 2
	// post2 := Post{
	// 	ID:    2,
	// 	Title: "Understanding HTTP in Go",
	// 	Url:   "http://antirez.com/news/108",
	// 	Body:  "In this post, we will explore how to build HTTP servers in Go. We will look at how to create routes, handle requests, and return responses. By the end, you'll have a good understanding of how to work with HTTP in Go.",
	// 	Date: time.Now(),
	// }


	// createTable(db)
	database.AddPost(db,shared.Post{
        Title: "Sample Title",
        Url:   "http://example.com/sample-url",
        Body:  "This is a sample body for testing purposes.",
        Date:  time.Now(), // Current time
    })	
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

	// handles /posts
	router.Get("/posts",handlers.HandleViewPosts(posts))

	


	http.ListenAndServe(":8080", router)




}
