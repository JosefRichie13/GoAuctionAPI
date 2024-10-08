package main

import (
	"database/sql"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for bidItem(). It requires 1 JSON Key itemID
type BidItemParameters struct {
	ItemID   string `form:"itemID" binding:"required"`
	BidPrice int    `json:"bidPrice" binding:"required,numeric,gte=0"`
	BidEmail string `json:"buyerEmail" binding:"required,email"`
}

// Starts the auction for an item.
func bidItem(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, BidItemParameters
	var bidItemParameters BidItemParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&bidItemParameters) != nil {
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
	queryToCheckExistingItem := `SELECT ID, I_STATUS, B_PRICE FROM AUCTION WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingItem, bidItemParameters.ItemID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type CheckItem struct {
		checkID     string
		checkStatus string
		checkPrice  int
	}

	var checkItem CheckItem
	result.Scan(&checkItem.checkID, &checkItem.checkStatus, &checkItem.checkPrice)

	// If the length of the ID is 0, means there is no item by that ID, we reject with a 404
	if len(checkItem.checkID) == 0 {
		c.JSON(404, gin.H{"status": "No Item with ID, " + bidItemParameters.ItemID + " exists"})
		return
	}

	// If the status of the item is not OPEN, means the item is not up for auction, we reject with a 403
	if (checkItem.checkStatus) != "OPEN" {
		c.JSON(403, gin.H{"status": "Item with ID, " + bidItemParameters.ItemID + " is currently not in auction"})
		return
	}

	// If the current price is greater than or equal to the bid price, means the bid is lower, we reject with a 403
	if checkItem.checkPrice >= bidItemParameters.BidPrice {
		c.JSON(403, gin.H{"status": "Your bid is lower or equal to the current bid of " + strconv.Itoa(checkItem.checkPrice)})
		return
	}

	// Finally we update the highest price and email of the bidder
	queryToUpdateABid := `UPDATE AUCTION SET B_EMAIL = $1, B_PRICE = $2 WHERE ID = $3;`
	db.QueryRow(queryToUpdateABid, bidItemParameters.BidEmail, bidItemParameters.BidPrice, bidItemParameters.ItemID)
	c.JSON(200, gin.H{"status": "Your bid of " + strconv.Itoa(bidItemParameters.BidPrice) + " is accepted.", "currentBid": bidItemParameters.BidPrice})

}

// Defining JSON body for getCurrentBidDetails(). It requires 1 Query Parameter itemID.
type GetCurrentBidDetailsParams struct {
	ItemID string `form:"itemID" binding:"required"`
}

// Returns a specific Bid Details
func getCurrentBidDetails(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, GetCurrentBidDetailsParams
	var getCurrentBidDetailsParams GetCurrentBidDetailsParams

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&getCurrentBidDetailsParams) != nil {
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

	// Query the DB and result is held into the variable, result
	queryToCheckCurrentBid := `SELECT * FROM AUCTION WHERE ID = $1;`
	result := db.QueryRow(queryToCheckCurrentBid, getCurrentBidDetailsParams.ItemID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type GetBidDetails struct {
		ID       string
		I_Name   string
		I_Desc   string
		I_Status string
		S_Email  string
		S_Price  float32
		S_Media  string
		B_Email  string
		B_Price  float32
	}

	// Creating an instance of the struct, GetBidDetails
	var getBidDetails GetBidDetails

	// Scan the results into the struct
	result.Scan(&getBidDetails.ID, &getBidDetails.I_Name, &getBidDetails.I_Desc, &getBidDetails.I_Status, &getBidDetails.S_Email, &getBidDetails.S_Price, &getBidDetails.S_Media, &getBidDetails.B_Email, &getBidDetails.B_Price)

	// If the length of the ID is 0, means there is no item by that ID, we reject with a 404
	if len(getBidDetails.ID) == 0 {
		c.JSON(404, gin.H{"status": "No Item with ID, " + getCurrentBidDetailsParams.ItemID + " exists"})
		return
	}

	// If the status of the item is UNLISTED, means the item is not up for auction, we reject with a 403
	if (getBidDetails.I_Status) == "UNLISTED" {
		c.JSON(403, gin.H{"status": "Item with ID, " + getBidDetails.ID + " is not up for auction", "sellerEmail": getBidDetails.S_Email, "itemStatus": getBidDetails.I_Status})
		return
	}

	// If the status of the item is FINISHED, means the auction is over, we reject with a 403
	if (getBidDetails.I_Status) == "FINISHED" {
		c.JSON(403, gin.H{"status": "Item with ID, " + getBidDetails.ID + " is already sold", "buyerEmail": getBidDetails.B_Email, "itemStatus": getBidDetails.I_Status})
		return
	}

	// Finally we return all the specific bid details
	c.JSON(200, gin.H{"itemID": getBidDetails.ID, "itemName": getBidDetails.I_Name, "itemDescription": getBidDetails.I_Desc, "sellerEmail": getBidDetails.S_Email,
		"sellingPrice": getBidDetails.S_Price, "itemMedia": getBidDetails.S_Media, "currentBuyerEmail": getBidDetails.B_Email, "currentBuyingPrice": getBidDetails.B_Price})

}
