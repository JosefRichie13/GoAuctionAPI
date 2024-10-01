package main

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

func main() {

	request := gin.Default()

	// Configuring 8mb for file uploads
	request.MaxMultipartMemory = 8 << 20

	request.GET("/", landingPage)
	request.POST("/addAnAuctionItem", addAnAuctionItem)
	request.PUT("/updateAnAuctionItem", updateAnAuctionItem)
	request.PUT("/addMediaToItem", addMediaToItem)
	request.PUT("/startAuction", startAuction)
	request.PUT("/endAuction", endAuction)
	request.PUT("/resetAuction", resetAuction)
	request.PUT("/bidItem", bidItem)
	request.GET("/getCurrentBidDetails", getCurrentBidDetails)
	request.GET("/getUnlistedItems", getUnlistedItems)
	request.GET("/getOpenItems", getOpenItems)
	request.GET("/getSoldItems", getSoldItems)
	request.DELETE("/deleteItem", deleteItem)
	request.Run(":8083")

}

// Landing page route
func landingPage(c *gin.Context) {
	c.JSON(200, "Welcome to Auction API, currently being built")
}

// Defining JSON body for addAnAuctionItem(). It requires 4 JSON key's itemName, itemDesc, sellerEmail, sellerPrice.
type AddAnAuctionItemParameters struct {
	ItemName    string `json:"itemName" binding:"required"`
	ItemDesc    string `json:"itemDesc" binding:"required"`
	SellerEmail string `json:"sellerEmail" binding:"required,email"`
	SellerPrice int    `json:"sellerPrice" binding:"required,numeric,gte=0"`
}

func addAnAuctionItem(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, AddAnAuctionItemParameters
	var addAnAuctionItemParameters AddAnAuctionItemParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&addAnAuctionItemParameters) != nil {
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

	// Generate Unique ID for Auction Item and use the same for the media name
	generatedID := uniqueIDGenerator()

	// Add the item to the DB
	queryToAddAnItemForAuction := `INSERT INTO AUCTION (ID, I_NAME, I_DESC, I_STATUS, S_EMAIL, S_PRICE, S_MEDIA, B_EMAIL, B_PRICE) Values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	db.QueryRow(queryToAddAnItemForAuction, generatedID, sanitizeString(addAnAuctionItemParameters.ItemName), sanitizeString(addAnAuctionItemParameters.ItemDesc), "UNLISTED",
		addAnAuctionItemParameters.SellerEmail, addAnAuctionItemParameters.SellerPrice, "PATH_NOT_FOUND", "0", 0)
	c.JSON(200, gin.H{"status": "Item added for Auction", "itemID": generatedID})

}

/////////////////////

// Defining JSON body for deleteItem(). It requires 1 Query Parameter itemID.
type DeleteItemParams struct {
	ItemID string `form:"itemID" binding:"required"`
}

// Returns a specific Event Details
func deleteItem(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, DeleteItemParams
	var deleteItemParams DeleteItemParams

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&deleteItemParams) != nil {
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
	queryToCheckItemStatus := `SELECT ID, I_STATUS, S_MEDIA FROM AUCTION WHERE ID = $1;`
	result := db.QueryRow(queryToCheckItemStatus, deleteItemParams.ItemID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type GetBidDetails struct {
		ID       string
		I_Status string
		S_Media  string
	}

	// Creating an instance of the struct, GetBidDetails
	var getBidDetails GetBidDetails

	// Scan the results into the struct
	result.Scan(&getBidDetails.ID, &getBidDetails.I_Status, &getBidDetails.S_Media)

	// If the length of the ID is 0, means there is no item by that ID, we reject with a 404
	if len(getBidDetails.ID) == 0 {
		c.JSON(404, gin.H{"status": "No Item with ID, " + deleteItemParams.ItemID + " exists"})
		return
	}

	// If the status of the item is OPEN, means the item is currently under auction, we cannot delete it, so reject with a 403
	if (getBidDetails.I_Status) == "OPEN" {
		c.JSON(403, gin.H{"status": "Item with ID, " + deleteItemParams.ItemID + " is currently in auction and cannot be deleted"})
		return
	}

	// Finally we delete the item along with its media
	queryToDeleteAnItem := `DELETE FROM AUCTION WHERE ID = $1;`
	db.QueryRow(queryToDeleteAnItem, deleteItemParams.ItemID)

	error := os.RemoveAll(getBidDetails.S_Media)
	if error != nil {
		c.JSON(200, gin.H{"status": "Item with " + deleteItemParams.ItemID + " deleted, but could not find the media to delete."})
		return
	}

	c.JSON(200, gin.H{"status": "Item with " + deleteItemParams.ItemID + " deleted."})

}
