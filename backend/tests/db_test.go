package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// test db
// test query
// test api endpoints
//
//	/ landing
//	/allposts posttab
//
// /loginView login page
// /adminView admin page
// /about
// /contact
// /time

// Testing the db connection
func TestDbConn(t *testing.T) {
	err := godotenv.Load("/root/backend/.env")
	if err != nil {
		log.Println("Error loading .env file from godotenv, you may be using docker for the env?")
	}

	connStr := os.Getenv("CONN_STR")

	db, err := sql.Open("postgres", connStr)
    log.Println("this is connstr",connStr)

	if err != nil {
		log.Println("err connecting to database", err)
		t.Failed()
	}
	defer db.Close()
}

func TestDbQuery(t *testing.T) {
	err := godotenv.Load("/root/backend/.env")
	if err != nil {
		log.Println("Error loading .env file from godotenv, you may be using docker for the env?")
	}

	connStr := os.Getenv("CONN_STR")
    log.Println(connStr)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Println("err connecting to database", err)
		t.Failed()
	}
	defer db.Close()

	rows, err := db.Query("SELECT COUNT(*) FROM posts")
	if err != nil {
		log.Println("Could not fetch rows:", err)
        t.Fail()
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Println("Error scanning row:", err)
            t.Fail()
		}
	}

	// for _,val := range allPosts{
	// 	fmt.Println("works")
	// 	fmt.Println(val)
	// }

	// Check for any error after finishing the iteration
	if err = rows.Err(); err != nil {
		log.Fatalln("Error iterating over rows:", err)
        t.Fail()
	}
	log.Println("We got the number of rows, ", count)
}



