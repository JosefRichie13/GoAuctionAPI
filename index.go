package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

func main() {

	request := gin.Default()

	request.GET("/", landingPage)
	request.POST("/addAnAuctionItem", addAnAuctionItem)
	request.POST("/updateAnAuctionItem", updateAnAuctionItem)
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
