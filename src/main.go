// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package main

import (
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	dsn = "DSN_STRING"
)

// test upstream call
func webServer() {
	http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("SampleApplication"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Hello XRay Go SDK Customer!"))
	})))
	log.Println("SampleApp is listening on port 8000. Visit localhost:8000/ in your browser to generate segments on incoming request!")
	http.ListenAndServe(":8000", nil)
}

// test downstream aws calls
func testAWSCalls() {
	// example of custom segment
	ctx, root := xray.BeginSegment(context.Background(), "AWS SDK Calls")
	defer root.Close(nil)

	awsSess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable,})
	if err != nil {
		log.Fatalf("failed to open aws session")
	}

	// S3 and SQS Clients
	s3Client := s3.New(awsSess)
	sqsClient := sqs.New(awsSess)

	// XRay Setup
	xray.AWS(s3Client.Client)
	xray.AWS(sqsClient.Client)

	// List SQS queues
	if _, err = sqsClient.ListQueuesWithContext(ctx, nil); err != nil {
		log.Println(err)
		return
	}

	// List s3 buckets
	if _, err = s3Client.ListBucketsWithContext(ctx, nil); err != nil {
		log.Println(err)
		return
	}
}

// trace SQL queries
func testSQL() {
	dsnString := os.Getenv(dsn)
	if dsnString == "" {
		log.Println("Set DSN_STRING environment variable as a connection string to the database")
		return
	}

	ctx, root := xray.BeginSegment(context.Background(), "SQL")
	defer root.Close(nil)

	db, _ := xray.SQLContext("mysql", dsnString)
	defer db.Close()

	// Create a table
	if _, err := db.ExecContext(ctx, "CREATE TABLE ID1 (val int)"); err != nil {
		log.Println(err)
		return
	}

	// Transaction
	if err := transaction(ctx, db); err != nil {
		log.Println(err)
		return
	}

	// Drop the table
	if _, err := db.ExecContext(ctx, "DROP TABLE ID1"); err != nil {
		log.Println(err)
		return
	}
}

func transaction(ctx context.Context, db *sql.DB) error {
	// Begin Tx
	tx, _ := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	// Populate the data
	if _, err := db.ExecContext(ctx, "INSERT INTO ID1 (val) VALUES (?)", 1); err != nil {
		return err
	}

	// Query the data
	r, _ := tx.QueryContext(ctx, "SELECT val FROM ID1 WHERE val = ?", 1)
	defer r.Close()

	for r.Next() {
		var val int
		err := r.Scan(&val)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func main() {
	log.Println("SampleApp Starts")
	testAWSCalls() // Outgoing AWS SDK calls
	time.Sleep(1 * time.Second)
	testSQL() // SQL calls
	time.Sleep(1 * time.Second)
	webServer() // Upstream calls
}

