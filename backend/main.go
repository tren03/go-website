package main

import (
	// "database/sql"

	// "fmt"
	// "log"
	"fmt"
	"net/http"
	// "os"

	// "time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"go-server/handlers"
	// "go-server/shared"
	// "github.com/joho/godotenv"
)

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// connStr:=os.Getenv("CONN_STR")

	// db,err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatalln("Could not connect to Database:", err)
	// }

	// createTable(db)
	// database.AddPost(db,shared.Post{
	//     Title: "Sample Title",
	//     Url:   "http://example.com/sample-url",
	//     Body:  "This is a sample body for testing purposes.",
	//     Date:  time.Now(), // Current time
	// })
	// delPost(db,1)
	// putAllToPostsToArray(db)

	// posts := database.PutAllToPostsToArray(db)

	// chi router to router http requests
	router := chi.NewRouter()

	// cors config
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// handles / (landing page)
	router.Get("/", handlers.HandleLanding)

	// handles /allposts to provide most recent data always
	router.Get("/allposts", handlers.HandleAllPosts)

	// handles json response to index.html
	router.Get("/posts", handlers.HandleAutoViewPosts)

	// handles rendering the login page
	router.Get("/loginView", handlers.HandleLoginView)

	// handles /login which is hit by /loginView
	router.Post("/login", handlers.HandleLogin)

	// handles rendering the adminPage
	router.Get("/adminView", handlers.HandleAdminView)

	// server end point that adds post /addPost
	router.Post("/addPost", handlers.HandleAddPost)

	// server end point to delete post
	router.Delete("/deletePost/{id}", handlers.HandleDeletePost)

	// handles /about for rendering the about page
	router.Get("/about", handlers.HandleAbout)

	// handles /about for rendering the about page
	router.Get("/contact", handlers.HandleContact)


	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("/root/frontend/assets/"))))

	fmt.Println("server started at 8080")
	http.ListenAndServe(":8080", router)

}
