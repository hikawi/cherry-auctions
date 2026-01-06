---
title: To-do Checklist
parent: Other Specifications
---

# CherryAuctions: Official SRS Checklist

## 1. Guest Role

### 1.1. Menu System

- [x] Display all categories
- [x] Categories have at least 2 levels

### 1.2. Home Page

- [x] Top 5 ending soon products
- [x] Top 5 products with most bids
- [x] Top 5 products with highest bids

### 1.3. Products Catalog

- [x] Pagination
- [x] From categories

### 1.4. Products Searching

- [x] Full-text Search
- [x] Search by name
- [x] Search by category
- [x] Paginated results
- [x] Ordered results
  - [x] Ending time
  - [x] Price
- [x] New products (within some minutes) have emphasis
- [x] Product Display on Catalog
  - [x] Thumbnail
  - [x] Name
  - [x] Current Price
  - [x] Current Highest Bid
  - [x] BIN price
  - [x] Posted date
  - [x] Remaining time
  - [x] Number of bids

### 1.5. Products Details

- [x] All content of the product
  - [x] Thumbnail
  - [x] Other images
  - [x] Product name
  - [x] Current price
  - [x] BIN price
  - [x] Seller info & rating
  - [x] Highest bidder info & rating
  - [x] Posted date
  - [x] Expired at
  - [x] Relative expired in if less than 3 days
  - [x] Description
- [x] Questions and Answers
- [x] Similar products

### 1.6. Registration

- [x] reCaptcha
- [x] Passwords are hashed
- [x] Information
  - [x] OTP
  - [x] Email
  - [x] Name
  - [x] Address

## 2. Bidder Role

### 2.1. Watch List

- [x] Can add/remove from catalog
- [x] Can add/remove at details page

### 2.2. Bid

- [x] At details page
- [x] Allowed only if product allows it, or rating is 80% or greater.
- [x] System enforces a reasonable bid (current bid + bid step at least)
- [x] Confirmation dialog for bidding

### 2.3. Bidding History

- [x] All bids are masked
- [x] Needs to show price and time, and a masked name

### 2.4. Question

- [x] At details page
- [x] Seller receives email when someone asks a question

### 2.5. Profile Management

- [x] Change email, name
- [ ] Change password (requires old password)
- [ ] View total rating and all ratings
- [x] View favorites / watch list
- [x] View current bids
- [ ] View won auctions
  - [ ] Rate the product's seller (+1) or (-1), with a feedback

### 2.6. Request Seller privileges

- [x] Bidder can request to be a seller
- [x] Admin can approve or reject

## 3. Seller Role

### 3.1. Post an Auction

- [x] Product name
- [x] At least 3 images
- [x] Starting bid
- [x] Bid step
- [x] BIN price (optional)
- [x] Description
  - [x] WYSIWYG
- [x] Auto extends

### 3.2. Editing a Description

- [x] Modify the description of a posted product
- [x] New changes are appended to the old description, but can't replace them

### 3.3. Deny a Bidder

- [ ] At details page
- [ ] A denied bidder can no longer bid on a product
- [ ] If the bidder gets denied, the highest bid is moved to the second highest

### 3.4. Answer a User

- [x] At details page

### 3.5. Profile Management

- [x] View my auctions
- [ ] View expired auctions
  - [ ] Can rate or feedback on the highest bidder
  - [ ] Can cancel the auction and automatically mark as -1 with a description
        "Bidder did not follow through"

## 4. Admin Role

### 4.1. Categories

- [ ] Can't delete a category with a product
- [x] View all
- [x] View details
- [x] Add
- [x] Delete
- [x] Edit

### 4.2. Products

- [ ] Can remove an auction

### 4.3. Users

- [ ] Basic management tools
- [x] View bidders requesting privileges
- [x] Can approve

### 4.4. Dashboard

- [ ] Charts on auctions, pricing, new users, new sellers

## 5. General Role

### 5.1. Login

- [x] Login
- [ ] (Optional) add Oauth

### 5.2. Update Profile

- [x] Name
- [x] Email

### 5.3. Change password

- [ ] Both passwords are hashed with bcrypt or scrypt or similar

### 5.4. Forgot password

- [ ] OTP code

## 6. System

### 6.1. Mailing System

- [x] New bid, to seller, to bidder, to previous bidder
- [ ] When bidder is denied the bid to bidder
- [x] Auction expired to seller
- [x] Auction ended to seller and winner
- [x] User questions, to the bidder
- [x] Seller answers, to the questioner, and all current bidders, or questioned

### 6.2. Automated bidding (Optional)

- [ ] Bidder sets a maximum price
- [ ] The product will kept getting bidded on if possible
- [ ] If 2 bidders have the same price, then the one who bidded first wins.

### 6.3. Payment System

- [ ] When auction ends, the detail page can lead to the "checkout"
- [ ] Other bidders just see: Auction ended
- [ ] Checkout flow:
  - [ ] Through Stripe
  - [ ] Winner sends the address
  - [ ] Seller confirms delivery
  - [ ] Winner confirms received
  - [ ] Transaction over, and the two rate each other
- [ ] The seller can cancel, and automatically mark the winner as -1.
- [ ] Chat interface for completing the product

## 7. Others

These requirements are non-negotiable.

### 7.1. Technical

- [x] Web App with Client-side rendering
- [x] Backend is RESTful
- [x] Backend Requirements
  - [x] Swagger
  - [x] Validation on all routes
  - [x] Logging and monitoring (ELK stack or similar)
  - [x] Security with JWT key pair
- [x] Frontend
  - [x] Client-side router
  - [x] Form processing and Validation
  - [x] State management
  - [x] Same design system for entire website

### 7.2. Data

- [x] At least 20 products, with descriptions and images
- [x] Products need at least some bids

### 7.3. Source Management

- [x] Github (non-negotiable, course fail if not checked)
