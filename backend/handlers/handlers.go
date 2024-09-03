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
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"time"
)

// handles /
func HandleRoot(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,"../frontend/index.html")
}

// handles /LoginView
func HandleLoginView(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,"../frontend/login.html")
}



// handles /adminView - need to send admin.html as a response
func HandleAdminView(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,"../frontend/admin.html")

}

// handles /posts which gets hit by / for all posts go templates
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


type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func createJWT(username string)(string,error){
	expirationTime := time.Now().Add(1 * time.Minute)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtKey := []byte(os.Getenv("JWT_KEY"))
	
	claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
        return "", err
    }

    return tokenString, nil
}


// verify jwt on client request
func verifyJWT(tokenStr string) (*Claims, error) {
    claims := &Claims{}
	jwtKey := []byte(os.Getenv("JWT_KEY"))
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        // Ensure that the signing method used is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}


// handles /login which gets hit by loginView, currently figuring out how to auth myself to add delete posts, thinking about JWT maybe
type LoginReq struct{
	Username string `json:"username"`
	Password string `json:"password"`
}
func HandleLogin(w http.ResponseWriter,r *http.Request){

	bodyBytes,err:=io.ReadAll(r.Body)	
	if err!= nil{
		log.Println("error while reading login response", err)
	}

	var login_details LoginReq 
	err =json.Unmarshal(bodyBytes,&login_details)
	if err!=nil{
		log.Println("problem parsing the login response",err)
	}

	fmt.Printf("Username: %s\n", login_details.Username)
	fmt.Printf("Password: %s\n", login_details.Password)

	if login_details.Username==string(os.Getenv("USER_NAME")) && login_details.Password==string(os.Getenv("PASSWORD")){
		
		token, err := createJWT("vish")
    	if err != nil {
        	http.Error(w, "Failed to create token", http.StatusInternalServerError)
        return
    	}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    token,
			Expires:  time.Now().Add(5 * time.Minute),
			HttpOnly: true, // Prevents JavaScript from accessing the cookie
			Secure:   true, // Ensures the cookie is sent only over HTTPS
			SameSite: http.SameSiteStrictMode, // Helps protect against CSRF attacks
		})
		w.Header().Set("Content-Type", "application/json")
    	// w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))	
	}else{
		log.Println("invalid login")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
	}


}


// handles /admin
func Admin(w http.ResponseWriter,r *http.Request){
    // Get the token from the request cookie
    cookie, err := r.Cookie("jwt_token")
    if err != nil {
        if err == http.ErrNoCookie {
            http.Error(w, "No token found", http.StatusUnauthorized)
			log.Println("No token")
            return
        }
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Verify the token
    claims, err := verifyJWT(cookie.Value)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Unauth")
        return
    }

    // Check the expiration
    if time.Now().After(claims.ExpiresAt.Time) {
        http.Error(w, "Token has expired", http.StatusUnauthorized)
		log.Println("Token expired")
        return
    }

    // Proceed with the request, now that the token is verified
    
	adminHTMLcontent,err:= os.ReadFile("../frontend/index.html")
	if(err!=nil){
		http.Error(w, "Failed to fetch admin", http.StatusInternalServerError)	
		return 	
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(adminHTMLcontent))
	log.Println("success")
}


//assembles content to sent to admin page
// func AdminContentHandler(){

// }