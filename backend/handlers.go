package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"go-server/shared"
)

func handleRoot(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w,"welcome to root")
}

// func handleViewPosts(w http.ResponseWriter,r *http.Request){
// 	fmt.Fprintf(w,"%v",FinalHTML)
// }


func handleViewPosts(posts []shared.Post) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		log.Println("hit posts")
		tmpl, err:= template.ParseFiles("templates/post.html")
		if err!=nil{
			http.Error(w,"Error while parsing",http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer

		// Execute the template and write the output to the buffer
		err = tmpl.Execute(&buf, posts)
		if err != nil {
			http.Error(w, "Error while rendering template", http.StatusInternalServerError)
			return
		}

		// Print the rendered HTML to the console
		// fmt.Println(buf.String())
		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write(buf.Bytes())
		if err != nil {
			http.Error(w, "Error while writing response", http.StatusInternalServerError)
		}
	}	
}


// func closureFuncCheck(name string) http.HandlerFunc{ // we use this since to return a handler func as it allows us to take in arguments
// 	return func(w http.ResponseWriter,r *http.Request){
// 		fmt.Fprintf(w,"we took the argument !! %v",name)
// 	}
// }