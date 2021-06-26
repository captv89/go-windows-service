package main

import (
	"fmt"
	"strconv"
)

// Create DB and Table If not exist
func createDB() {
	query, err := db.Prepare("CREATE TABLE IF NOT EXISTS vlog (id INTEGER PRIMARY KEY, VideoTitle TEXT, DownloadStatus TEXT, Activity TEXT , ErrorMsg TEXT, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	errorHandler(err)
	query.Exec()
}

// Make Entry in the Video Log Table (vlog)
func makeEntry(title string, status string, action string, e string) {
	query, err := db.Prepare("INSERT INTO vlog (VideoTitle, DownloadStatus, Activity , ErrorMsg) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	query.Exec(title, status, action, e)
}

//Print Contents of the table
func printTable() {
	rows, err := db.Query("SELECT id, VideoTitle, DownloadStatus, Activity, ErrorMsg, Timestamp FROM vlog")
	errorHandler(err)

	var id int
	var VideoTitle, DownloadStatus, ErrorMsg, Activity, Timestamp string

	for rows.Next() {
		rows.Scan(&id, &VideoTitle, &DownloadStatus, &Activity, &ErrorMsg, &Timestamp)
		fmt.Println(strconv.Itoa(id) + ": " + VideoTitle + " " + DownloadStatus + " " + Activity + " " + ErrorMsg + " " + Timestamp)
	}
}
