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
	queryToCheckExistingItem := `SELECT ID, I_STATUS, S_MEDIA FROM AUCTION WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingItem, auctionParameters.ItemID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type CheckItem struct {
		checkID     string
		checkStatus string
		checkMedia  string
	}

	var checkItem CheckItem
	result.Scan(&checkItem.checkID, &checkItem.checkStatus, &checkItem.checkMedia)

	// Checks if the media is present, if not, we won't start the auction
	if checkItem.checkMedia == "PATH_NOT_FOUND" {
		c.JSON(404, gin.H{"status": "No media found for the Item with ID, " + auctionParameters.ItemID + ". Cannot start the auction"})
		return
	}

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
			c.JSON(403, gin.H{"status": "Item with ID, " + auctionParameters.ItemID + " is already sold or not listed for auction or currently in auction"})
		}

	} else {
		c.JSON(404, gin.H{"status": "No Item with ID, " + auctionParameters.ItemID + " exists"})
	}

}

// Ends the auction for an item.
func endAuction(c *gin.Context) {

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

		// We check if the item is OPEN, only then we update it to FINISHED, which means auction is over for that item
		if (checkItem.checkStatus) == "OPEN" {
			queryToStartAuction := `UPDATE AUCTION SET I_STATUS = "FINISHED" WHERE ID = $1;`
			db.QueryRow(queryToStartAuction, auctionParameters.ItemID)
			c.JSON(200, gin.H{"status": "Auction for Item with ID, " + auctionParameters.ItemID + " is over."})
		} else {
			c.JSON(403, gin.H{"status": "Item with ID, " + auctionParameters.ItemID + " is already sold or not listed for auction or currently in auction"})
		}

	} else {
		c.JSON(404, gin.H{"status": "No Item with ID, " + auctionParameters.ItemID + " exists"})
	}

}

// Resets the auction for an item.
func resetAuction(c *gin.Context) {

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

		// We check if the item is FINISHED, only then we update it to UNLISTED, which means auction is reset.
		// We set the buyer_price and buyer_email to 0
		if (checkItem.checkStatus) == "FINISHED" {
			queryToStartAuction := `UPDATE AUCTION SET I_STATUS = "UNLISTED", B_EMAIL = $1, B_PRICE = $2 WHERE ID = $3;`
			db.QueryRow(queryToStartAuction, "0", 0, auctionParameters.ItemID)
			c.JSON(200, gin.H{"status": "Auction for Item with ID, " + auctionParameters.ItemID + " is reset."})
		} else {
			c.JSON(403, gin.H{"status": "Item with ID, " + auctionParameters.ItemID + " is not ready for auction or currently in auction"})
		}

	} else {
		c.JSON(404, gin.H{"status": "No Item with ID, " + auctionParameters.ItemID + " exists"})
	}

}
