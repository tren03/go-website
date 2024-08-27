package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	// "os"
	// "log"
	// "strconv"
	// "sync"
)

type Post struct{
	ID int `json:"ID"`
	Title string `json:"Title"`
	Url string `json:"Url"`
	Body string `json:"Body"`
	Date  time.Time `json:"Date"`
	
}







func main(){
	// Sample Post 1

	post1 := Post{
		ID:    1,
		Title: "Introduction to Go",
		Url:   "https://cpu.land",
		Body:  "Go is an open-source programming language designed for simplicity and efficiency. In this post, we will cover the basics of Go, including its syntax, types, and how to write a simple Go program.",
		Date: time.Now(),
	}
	//Sample Post 2
	post2 := Post{
		ID:    2,
		Title: "Understanding HTTP in Go",
		Url:   "http://antirez.com/news/108",
		Body:  "In this post, we will explore how to build HTTP servers in Go. We will look at how to create routes, handle requests, and return responses. By the end, you'll have a good understanding of how to work with HTTP in Go.",
		Date: time.Now(),
	}

	posts:=[]Post{post1,post2}

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
	router.Get("/",handleRoot)

	// handles /posts
	router.Get("/posts",handleViewPosts(posts))

	


	http.ListenAndServe(":8080", router)




}