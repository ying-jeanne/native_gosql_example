package main

import (
	"database/sql"
	"log"
	"time"

	// the database driver for sqlite3 for example
	_ "github.com/mattn/go-sqlite3"
)

const file string = "grafana.db"

func main() {
	// here we prepare db connection which is type *sql.DB with driver of sqlite3
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}
	// here we actually test that db connection is etablished
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	setUser(db)
	// Query is used for any query that returns rows. Exec is used for queries that doesn't return data.
	getUser(db)

	// apparently SQL doesn't support nullable columns very well, we must have some of the nullable type redefined by ourselves?!
	// COALESCE(other_field, '') as otherField could do the trick
}

type Team struct {
	ID        int
	Name      string
	OrgID     int
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func setUser(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO team(name, org_id, created, updated, email) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec("wangy", 0, time.Now(), time.Now(), "w.x@gmail.com")
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	_, err = db.Exec("DELETE FROM team WHERE id = ?", lastId)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("ID = %d, user deleted", lastId)
	}
}

func getUser(db *sql.DB) {
	// concatenating string is a bad idea here because of SQL injection attacks, instead we can prepare queries and reuse it later on.
	// rows, err := db.Query("select id, name, org_id, created, updated, email from user where id= ?", 1)
	stmt, err := db.Prepare("select id, name, org_id, created, updated, email from user where id= ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// actual query execution that is doing prepare, execute and close a preapred statement, it is 3 round-trips, might triple the number of database interactions
	rows, err := stmt.Query(1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var team Team
		err := rows.Scan(&team.ID, &team.Name, &team.OrgID, &team.CreatedAt, &team.UpdatedAt, &team.Email)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(team)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
