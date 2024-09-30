package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for updateAnAuctionItem(). It requires 4 JSON key's itemName, itemDesc, sellerEmail, sellerPrice.
type UpdateAnAuctionItemParameters struct {
	ItemID      string `json:"itemID" binding:"required"`
	ItemName    string `json:"itemName" binding:"required"`
	ItemDesc    string `json:"itemDesc" binding:"required"`
	SellerEmail string `json:"sellerEmail" binding:"required,email"`
	SellerPrice int    `json:"sellerPrice" binding:"required,numeric,gte=0"`
}

func updateAnAuctionItem(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, AddAnAuctionItemParameters
	var updateAnAuctionItemParameters UpdateAnAuctionItemParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&updateAnAuctionItemParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./AUCTION.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Check if the ItemID exists in the DB by querying for the ID
	// Result is scanned into the variable, checkResult
	queryToCheckExistingItem := `SELECT ID FROM AUCTION WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingItem, updateAnAuctionItemParameters.ItemID)
	var checkResult string
	result.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is an item by that ID
	// So, update the item
	// Else, its rejected with a 404 as there is no item by that ID
	if len(checkResult) > 0 {

		// Reset event start and end date
		queryToUpdateAnEvent := `UPDATE AUCTION SET I_NAME = $1, I_DESC = $2, S_EMAIL = $3, S_PRICE = $4 WHERE ID = $5;`
		db.QueryRow(queryToUpdateAnEvent, updateAnAuctionItemParameters.ItemName, updateAnAuctionItemParameters.ItemDesc,
			updateAnAuctionItemParameters.SellerEmail, updateAnAuctionItemParameters.SellerPrice, updateAnAuctionItemParameters.ItemID)
		c.JSON(200, gin.H{"status": "Item with, " + updateAnAuctionItemParameters.ItemID + " updated."})

	} else {
		c.JSON(404, gin.H{"status": "No Item with ID, " + updateAnAuctionItemParameters.ItemID + " exists"})
	}

}
