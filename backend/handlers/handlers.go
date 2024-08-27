package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-server/shared"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func HandleRoot(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,"../frontend/index.html")
}

func HandleLoginView(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,"../frontend/login.html")
}

// handles /admin
func HandleAdminView(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,"../frontend/admin.html")
}

// func handleViewPosts(w http.ResponseWriter,r *http.Request){
// 	fmt.Fprintf(w,"%v",FinalHTML)
// }


func HandleViewPosts(posts []shared.Post) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){

		log.Println("hit posts")
		tmpl, err:= template.ParseFiles("templates/post.html")
		if err!=nil{
			log.Println("Error while parsing",err)
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


type LoginReq struct{
	Username string `json:"username"`
	Password string `json:"password"`
}
func HandleLogin(w http.ResponseWriter,r *http.Request){

	bodyBytes,err:=io.ReadAll(r.Body)	
	if err!= nil{
		log.Println("erro while reading login response", err)
	}

	var login_details LoginReq 
	err =json.Unmarshal(bodyBytes,&login_details)
	if err!=nil{
		log.Println("problem parsing the login response",err)
	}

	fmt.Printf("Username: %s\n", login_details.Username)
	fmt.Printf("Password: %s\n", login_details.Password)

	if login_details.Username==string(os.Getenv("USER_NAME")) && login_details.Password==string(os.Getenv("PASSWORD")){
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Login successful"}`))

		log.Println("Login success, redirecting to /admin")
		http.Redirect(w,r,"/admin",http.StatusSeeOther)
		return		
	}else{
		log.Println("invalid login")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
	}


}

func HandleAdmin(w http.ResponseWriter,r *http.Request){
		

}
	



// func closureFuncCheck(name string) http.HandlerFunc{ // we use this since to return a handler func as it allows us to take in arguments
// 	return func(w http.ResponseWriter,r *http.Request){
// 		fmt.Fprintf(w,"we took the argument !! %v",name)
// 	}
// }