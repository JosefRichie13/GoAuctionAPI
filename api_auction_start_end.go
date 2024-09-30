package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for startAuction(). It requires 1 JSON Key itemID
type AuctionParameters struct {
	ItemID string `form:"itemID" binding:"required"`
}

// Starts the auction for an item.
func startAuction(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, AuctionParameters
	var auctionParameters AuctionParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&auctionParameters) != nil {
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
	queryToCheckExistingItem := `SELECT ID, I_STATUS FROM AUCTION WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingItem, auctionParameters.ItemID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type CheckItem struct {
		checkID     string
		checkStatus string
	}

	var checkItem CheckItem
	result.Scan(&checkItem.checkID, &checkItem.checkStatus)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is an item by that ID
	// So, update the item
	// Else, its rejected with a 404 as there is no item by that ID
	if len(checkItem.checkID) > 0 {

		// We check if the item is UNLISTED, only then we update it to OPEN, which means its open for auction
		if (checkItem.checkStatus) == "UNLISTED" {
			queryToStartAuction := `UPDATE AUCTION SET I_STATUS = "OPEN" WHERE ID = $1;`
			db.QueryRow(queryToStartAuction, auctionParameters.ItemID)
			c.JSON(200, gin.H{"status": "Item with ID, " + auctionParameters.ItemID + " is open for Auction."})
		} else {
			c.JSON(403, gin.H{"status": "Item with ID, " + auctionParameters.ItemID + " is already sold or not listed for auction"})
		}

	} else {
		c.JSON(404, gin.H{"status": "No Item with ID, " + auctionParameters.ItemID + " exists"})
	}

}
