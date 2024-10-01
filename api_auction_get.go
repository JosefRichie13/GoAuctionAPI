package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Returns all the unlisted items, i.e. items not yet opened for auction
func getUnlistedItems(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./AUCTION.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Query the DB and result is held into the variable, result
	queryToCheckUnlistedItems := `SELECT ID, I_NAME, I_DESC, I_STATUS, S_MEDIA, S_EMAIL, S_PRICE FROM AUCTION WHERE I_STATUS = 'UNLISTED';`
	result, error := db.Query(queryToCheckUnlistedItems)

	// If there's any error when querying, return it
	if error != nil {
		c.JSON(500, gin.H{"status": "Could not execute Query"})
		return
	}
	defer result.Close()

	// Defining a struct to hold the data queried by the query
	type GetItemDetails struct {
		ID       string  `json:"itemId"`
		I_Name   string  `json:"itemName"`
		I_Desc   string  `json:"itemDescription"`
		I_Status string  `json:"itemStatus"`
		S_Media  string  `json:"itemMedia"`
		S_Email  string  `json:"sellerEmail"`
		S_Price  float32 `json:"sellingPrice"`
	}

	// Creating a slice from the struct
	getItemDetails := []GetItemDetails{}

	for result.Next() {

		//Creating a new struct
		GetItemDetails := GetItemDetails{}

		// Scan the results into the struct
		result.Scan(&GetItemDetails.ID, &GetItemDetails.I_Name, &GetItemDetails.I_Desc, &GetItemDetails.I_Status, &GetItemDetails.S_Media, &GetItemDetails.S_Email, &GetItemDetails.S_Price)

		// Append to the slice
		getItemDetails = append(getItemDetails, GetItemDetails)
	}

	c.JSON(200, gin.H{"allUnlistedItems": getItemDetails})
}

// Returns all the open items, i.e. items currently in auction
func getOpenItems(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./AUCTION.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Query the DB and result is held into the variable, result
	queryToCheckOpenItems := `SELECT * FROM AUCTION WHERE I_STATUS = 'OPEN';`
	result, error := db.Query(queryToCheckOpenItems)

	// If there's any error when querying, return it
	if error != nil {
		c.JSON(500, gin.H{"status": "Could not execute Query"})
		return
	}
	defer result.Close()

	// Defining a struct to hold the data queried by the query
	type GetItemDetails struct {
		ID       string  `json:"itemId"`
		I_Name   string  `json:"itemName"`
		I_Desc   string  `json:"itemDescription"`
		I_Status string  `json:"itemStatus"`
		S_Media  string  `json:"itemMedia"`
		S_Email  string  `json:"sellerEmail"`
		S_Price  float32 `json:"sellingPrice"`
		B_Email  string  `json:"buyerEmail"`
		B_Price  float32 `json:"buyingPrice"`
	}

	// Creating a slice from the struct
	getItemDetails := []GetItemDetails{}

	for result.Next() {

		//Creating a new struct
		GetItemDetails := GetItemDetails{}

		// Scan the results into the struct
		result.Scan(&GetItemDetails.ID, &GetItemDetails.I_Name, &GetItemDetails.I_Desc, &GetItemDetails.I_Status, &GetItemDetails.S_Email, &GetItemDetails.S_Price,
			&GetItemDetails.S_Media, &GetItemDetails.B_Email, &GetItemDetails.B_Price)

		// Append to the slice
		getItemDetails = append(getItemDetails, GetItemDetails)
	}

	c.JSON(200, gin.H{"allOpenItems": getItemDetails})
}
