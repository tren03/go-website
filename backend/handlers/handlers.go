package handlers

import (
	"encoding/json"
	"fmt"
	"go-server/shared"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"database/sql"
	"go-server/database"
	_ "github.com/lib/pq"
)
func HandleLanding(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "/root/frontend/landing.html")
}

func HandleAllPosts(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/root/frontend/index.html")
}

// handles /LoginView
func HandleLoginView(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/root/frontend/login.html")
}

// Handles about page
func HandleAbout(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/root/frontend/about.html")
}

// Handles about page
func HandleContact(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/root/frontend/contact.html")
}

// handles condition rendering of admin based on verified jwt
func HandleAdminView(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Unauthorized: Invalid token")
		return
	}

	// Check the expiration
	if time.Now().After(claims.ExpiresAt.Time) {
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		log.Println("Token expired")
		return
	}

	http.ServeFile(w, r, "/root/frontend/newadmin.html")
}

// Alternative to HandleViewPosts in which database connection is handled here only /posts
func HandleAutoViewPosts(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load("/root/backend/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	connStr := os.Getenv("CONN_STR")


	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("Could not connect to Database:", err)
	}
	defer db.Close()

	posts := database.PutAllToPostsToArray(db)

	// to render most recent post first
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	json_posts_data, err := json.Marshal(posts)
	if err != nil {
		log.Printf("error converting struct to json: %v\n", err)
		http.Error(w, "Failed to convert posts to JSON", http.StatusInternalServerError)
		return
	}

	// fmt.Println(string(json_posts_data))
	log.Println("hit posts")

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json_posts_data)
	if err != nil {
		http.Error(w, "Error while writing response", http.StatusInternalServerError)
	}

	// tmpl, err := template.ParseFiles("/root/backend/templates/post.html")
	// if err != nil {
	// 	log.Println("Error while parsing", err)
	// 	http.Error(w, "Error while parsing", http.StatusInternalServerError)
	// 	return
	// }

	// var buf bytes.Buffer

	// // Execute the template and write the output to the buffer
	// err = tmpl.Execute(&buf, posts)
	// if err != nil {
	// 	http.Error(w, "Error while rendering template", http.StatusInternalServerError)
	// 	return
	// }

	// // Print the rendered HTML to the console
	// // fmt.Println(buf.String())
	// w.Header().Set("Content-Type", "text/html")
	// _, err = w.Write(buf.Bytes())
	// if err != nil {
	// 	http.Error(w, "Error while writing response", http.StatusInternalServerError)
	// }

}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func createJWT(username string) (string, error) {
	expirationTime := time.Now().Add(500 * time.Minute)

	err := godotenv.Load("/root/backend/.env")
	if err != nil {
		log.Println("Error loading .env file")
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

	err := godotenv.Load("/root/backend/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

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
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading login response", err)
	}

	var login_details LoginReq
	err = json.Unmarshal(bodyBytes, &login_details)
	if err != nil {
		log.Println("problem parsing the login response", err)
	}

	fmt.Printf("Username: %s\n", login_details.Username)
	fmt.Printf("Password: %s\n", login_details.Password)

	fmt.Printf("exptected Username: %s\n", string(os.Getenv("USER_NAME")))
	fmt.Printf("expected Password: %s\n", string(os.Getenv("PASSWORD")))

	if login_details.Username == string(os.Getenv("USER_NAME")) && login_details.Password == string(os.Getenv("PASSWORD")) {

		token, err := createJWT("vish")
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}
		log.Println("created jwt")

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    token,
			Expires:  time.Now().Add(5 * time.Minute),
			HttpOnly: true,                    // Prevents JavaScript from accessing the cookie
			Secure:   true,                    // Ensures the cookie is sent only over HTTPS
			SameSite: http.SameSiteStrictMode, // Helps protect against CSRF attacks
		})
		w.Header().Set("Content-Type", "application/json")
	} else {
		log.Println("invalid login")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

}

func HandleAddPost(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Unauthorized: Invalid token")
		return
	}

	// Check the expiration
	if time.Now().After(claims.ExpiresAt.Time) {
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		log.Println("Token expired")
		return
	}

	// If we enter this section, it means we have verified our jwt and can perform authoratative operations

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading login response", err)
	}

	var p shared.Post
	p.Date = time.Now()
	err = json.Unmarshal(bodyBytes, &p)
	if err != nil {
		log.Println("problem parsing the login response", err)
	}

	fmt.Println(p)

	err = godotenv.Load("/root/backend/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	connStr := os.Getenv("CONN_STR")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("Could not connect to Database:", err)
	}
	defer db.Close() // This will be executed when AddPost() exits
	database.AddPost(db, p)

}

func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Unauthorized: Invalid token")
		return
	}

	// Check the expiration
	if time.Now().After(claims.ExpiresAt.Time) {
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		log.Println("Token expired")
		return
	}

	// If we enter this section, it means we have verified our jwt and can perform authoratative operations

	id_str := chi.URLParam(r, "id")
	err = godotenv.Load("/root/backend/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	connStr := os.Getenv("CONN_STR")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("Could not connect to Database:", err)
	}
	defer db.Close() // This will be executed when AddPost() exits

	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	database.DelPost(db, id)

}


