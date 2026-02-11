package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const dsn = "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable"

func main() {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		process(db, "S1")
	}()

	go func() {
		defer wg.Done()
		process(db, "S2")
	}()

	wg.Wait()

	var result sql.NullString
	_ = db.QueryRow(`SELECT result FROM demo.history WHERE id='H1'`).Scan(&result)
	fmt.Println("Final History Result:", result.String)
}

func process(db *sql.DB, childID string) {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// ✅ LOCK parent dulu
	_, err = tx.Exec(`
		SELECT id
		FROM demo.history
		WHERE id='H1'
		FOR UPDATE
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Update child
	_, err = tx.Exec(`
		UPDATE demo.child
		SET result='APPROVED'
		WHERE id=$1
	`, childID)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)

	// Aggregate check
	rows, err := tx.Query(`
		SELECT result
		FROM demo.child
		WHERE history_id='H1'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	allApproved := true
	for rows.Next() {
		var result sql.NullString
		_ = rows.Scan(&result)
		if result.String != "APPROVED" {
			allApproved = false
		}
	}

	if allApproved {
		_, err = tx.Exec(`
			UPDATE demo.history
			SET result='APPROVED'
			WHERE id='H1'
		`)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("History updated by", childID)
	}

	tx.Commit()
}
