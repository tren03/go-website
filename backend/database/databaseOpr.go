package database

import (
	"database/sql"
	"fmt"
	"go-server/shared"
	"log"
	"sync"
	"time"

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
func PutAllToPostsToArray(db *sql.DB) []shared.Post {

//	log.Println("CONCURRENT TESTING")
//	buf := make(chan shared.Post)
//
//    var wg sync.WaitGroup
//	start_single := time.Now()
//	c := 0
//	row := db.QueryRow("SELECT COUNT(*) FROM posts")
//	if err := row.Scan(&c); err != nil {
//		if err == sql.ErrNoRows {
//			log.Println("idk some error")
//		}
//	}
//    log.Println("getting count = ",c)
//
//	concurrentPosts := []shared.Post{}
//	go func() {
//		for item := range buf {
//            wg.Done()
//			concurrentPosts = append(concurrentPosts, item)
//		}
//	}()
//
//	//fetching all at once with go routines
//	for i := 0; i < c; i++ {
//		go func() {
//            wg.Add(1)
//			var article shared.Post
//			row := db.QueryRow("SELECT * FROM posts LIMIT 1 OFFSET %d", i)
//			if err := row.Scan(&article.ID, &article.Title, &article.Url, &article.Body, &article.Date); err != nil {
//				if err == sql.ErrNoRows {
//					log.Println("idk some error")
//				}
//			}
//			buf <- article
//		}()
//	}
//    wg.Wait()
//	time_taken_conc := time.Since(start_single)
//	log.Println("time taken concurrent go ", time_taken_conc)
//    log.Println(concurrentPosts)
//	log.Println("CONCURRENT TESTING DONE, resuming normal operations")

	start := time.Now()

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Println("Could not fetch rows:", err)
	}
	defer rows.Close()

	allPosts := []shared.Post{}

	for rows.Next() {
		var article shared.Post
		err := rows.Scan(&article.ID, &article.Title, &article.Url, &article.Body, &article.Date)
		if err != nil {
			log.Fatalln("Error scanning row:", err)
		}

		allPosts = append(allPosts, article)
	}

	// Check for any error after finishing the iteration
	if err = rows.Err(); err != nil {
		log.Fatalln("Error iterating over rows:", err)
	}
	time_taken := time.Since(start)

	log.Println("THE DATABASE FETCH OPERATION TOOK ", time_taken)
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
