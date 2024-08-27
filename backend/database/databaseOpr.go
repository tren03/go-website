package database

import (
	"database/sql"
	"fmt"
	"go-server/shared"
	"log"

	_ "github.com/lib/pq"
)

// func putAllPostsToConsole(db *sql.DB) {
// 	rows, err := db.Query("SELECT * FROM posts")
// 	if err != nil {
// 		log.Fatalln("Could not fetch rows:", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var article Post
// 		err := rows.Scan(&article.ID, &article.Title, &article.Url, &article.Body, &article.Date)
// 		if err != nil {
// 			log.Fatalln("Error scanning row:", err)
// 		}
// 		fmt.Printf("ID: %d\nTitle: %s\nURL: %s\nBody: %s\nDate: %s\n\n",
// 			article.ID, article.Title, article.Url, article.Body, article.Date)
// 	}

// 	// Check for any error after finishing the iteration
// 	if err = rows.Err(); err != nil {
// 		log.Fatalln("Error iterating over rows:", err)
// 	}
// }

// gets all the posts from the database and returns an Array of posts containt
func PutAllToPostsToArray(db *sql.DB) []shared.Post{
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatalln("Could not fetch rows:", err)
	}
	defer rows.Close()
	allPosts:= []shared.Post{}

	for rows.Next() {
		var article shared.Post
		err := rows.Scan(&article.ID, &article.Title, &article.Url, &article.Body, &article.Date)
		if err != nil {
			log.Fatalln("Error scanning row:", err)
		}

		allPosts = append(allPosts, article)
	}

	// for _,val := range allPosts{
	// 	fmt.Println("works")
	// 	fmt.Println(val)
	// }


	// Check for any error after finishing the iteration
	if err = rows.Err(); err != nil {
		log.Fatalln("Error iterating over rows:", err)
	}
	return allPosts
}

// Creates the table needed to store posts
// func createTable(db *sql.DB){
// 	sqlStatement := `
//     CREATE TABLE IF NOT EXISTS posts (
//         ID SERIAL PRIMARY KEY,
//         Title VARCHAR(255) NOT NULL,
//         Url TEXT NOT NULL,
//         Body TEXT,
//         Date TIMESTAMP NOT NULL
//     );`

//     // Execute the SQL statement
//     _, err := db.Exec(sqlStatement)
//     if err != nil {
//         log.Fatalln("Error creating table:", err)
//         return
//     }

//     fmt.Println("Table created successfully or already exists")

// }


func AddPost(db *sql.DB, p shared.Post) {
	sqlStatement := `
	INSERT INTO posts (title, url, body, date)
	VALUES ($1, $2, $3, $4)`

	// Execute the statement
	_, err := db.Exec(sqlStatement, p.Title, p.Url, p.Body, p.Date)
	if err != nil {
		log.Fatalln("insert error", err)
	} else {
		fmt.Println("Record inserted successfully")
	}
}

func DelPost(db *sql.DB, id int) {
    sqlStatement := `
    DELETE FROM posts
    WHERE ID = $1;`

    // Execute the SQL statement
    _, err := db.Exec(sqlStatement, id)
    if err != nil {
		log.Fatalln("insert error", err)
        
    }
}
    

