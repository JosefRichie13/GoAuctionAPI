# Backend for an Auction App using Gin Gonic, SQLite and Go

This repo has the code for an Auction App Backend. <br><br>

The basic flow of the App is 

* User creates an auction item, receives an ID in return
* User updates the item's media
* User then starts the auction
* Anyone can join in the auction to raise their bids
* User ends the auction
* User can delete the auction item <br><br>

The below REST API endpoints are exposed

* POST /addAnAuctionItem -- Creates an Auction Item

* PUT /updateAnAuctionItem -- Updates an Auction Item
  
* PUT /addMediaToItem -- Adds 1 media (photo/video) to an Auction Item
  
* PUT /startAuction -- Starts an auction
  
* PUT /endAuction -- Ends an auction
  
* PUT /resetAuction -- Resets an auction
  
* PUT /bidItem -- Bid's on an auction item
  
* GET /getCurrentBidDetails -- Returns a specific auction item's bid details
  
* GET /getCurrentBidDetailsPage -- Returns a specific auction item's bid details in an HTML page
  
* GET /getUnlistedItems -- Returns all the auction items that are not up for an auction
 
* GET /getOpenItems -- Returns all the auction items that are currently in an auction

* GET /getSoldItems -- Returns all the auction items that are sold in an auction

* DELETE /deleteItem -- Deletes a specific auction item <br><br>

The entire suite of endpoints with payloads are available in this HAR, [GoAuctionAPI.har](GoAuctionAPI.har)
