---
title: To-do Checklist
parent: Other Specifications
---

# 🍒 Cherry Auctions: Official Project Checklist

## 🛠️ Phase 1: Technical & System Requirements

- [x] **DevOps & Source Management**
  - [x] GitHub repository with frequent, meaningful commits (No bulk uploads).
  - [x] RESTful API implementation for backend endpoints.
  - [x] Swagger Documentation for all API endpoints.
  - [x] Logging and monitoring setup (Grafana, ELK, or similar).
- [ ] **Security & Authentication**
  - [x] JWT Implementation: AccessToken and RefreshTokens.
  - [x] Password Hashing: No plain text storage.
  - [x] reCaptcha integration for registration.
  - [x] OTP Verification system for registration and forgot password.
  - [ ] OAuth2 Options (Encouraged): Google, Facebook, X, GitHub.
- [x] **Frontend Standards**
  - [x] SPA Architecture with a dedicated Router service.
  - [x] Form processing and validation (e.g., Angular Forms).
  - [x] Consistent Design System across all screens.

---

## 👥 Phase 2: Role-Based Functionality

### 1. Guest (Public)

- [ ] **Menu System**
  - [x] At least two levels of categories (e.g., Electronics -> Laptops).
- [x] **Home Page Displays**
  - [x] "Top 5 Ending Soon Auctions".
  - [x] "Top 5 Auctions (most number of bids)".
  - [x] "Top 5 Highest Auctions".
- [ ] **Product Catalog & Search**
  - [ ] Paginated catalog view with category filters.
  - [x] Fuzzy search with full-text search.
  - [ ] Paginated search results with sorting (Price Asc/Desc, Ending Time).
  - [x] Eye-catching display for recently posted products.
- [ ] **Product Views**
  - [x] Catalog Attributes: Thumbnails, names, current bid, highest bidder, BIN price, posted date, remaining time, bids count.
  - [x] Product Details: Hero/extra photos, seller/bidder info & reviews, timestamps, description.
  - [x] Relative time display (if ending in ≤ 3 days).
  - [x] Questions & Answers section.
  - [ ] 5 Similar products in the same category.
- [x] **Registration**
  - [x] Fields: Full name, address, email.

### 2. Bidder

- [ ] **Watch List**
  - [ ] Save items from catalog or detail page.
  - [ ] View "My favorites" in profile.
- [ ] **Bidding System**
  - [ ] Bid restriction: Minimum 80% rating check (Likes/Total).
  - [ ] New account bidding permission (if seller allowed).
  - [ ] Minimum bid "step up" logic (Current + Step).
  - [ ] Confirmation dialog for bids.
- [ ] **History & Profile**
  - [ ] Masked bid history (Timestamp, Name, Price) in card or tabular view.
  - [ ] Profile management: Change email, name, address, password.
  - [ ] View personal rating details/comments from others.
  - [ ] Lists: "My current bids" and "My won bids" (with seller review/comment).
- [ ] **Permissions**
  - [ ] Request selling privileges after 7 days of usage.

### 3. Seller

- [ ] **Product Posting**
  - [ ] Required: Name, 3+ photos, start bid, step, description.
  - [ ] Optional: Buy-it-now price.
  - [ ] WYSIWYG text input for descriptions.
  - [ ] Option to enable auto-extend (last 5 minutes).
- [ ] **Product Management**
  - [ ] Edit: Append P.S. remarks to existing description.
  - [ ] Deny Bidder: Transfer bid from denied user to second highest.
  - [ ] Answer questions in QnA section.
- [ ] **Sales & Reviews**
  - [ ] View current and won products.
  - [ ] Review winner (+1/-1 and comment).
  - [ ] Stop transaction (Auto -1 to bidder with "The Top Bid did not follow through").

### 4. Administrator

- [ ] **Management (List, Details, Add, Delete, Update)**
  - [x] Category Management: Block delete if products exist in category.
  - [ ] Product Management: Remove any product.
  - [ ] User Management: Approve/Reject seller upgrade requests.
- [ ] **Global Controls**
  - [ ] Configure global auto-extend settings (When and how much).
- [ ] **Dashboard**
  - [ ] Analytics/Charts: User count, selling requests, financials, bids count.

---

## ⚙️ Phase 3: Core Systems

### 1. Mailing System (Email Notifications)

- [ ] Successful bid: To bidder, seller, and previous top bidder.
- [ ] Rejected bid: To the bidder.
- [ ] Expired auction (no bidders): To the seller.
- [ ] Ended auction: To the seller and winner.
- [ ] Question asked: To the seller.
- [ ] Question answered: To asker and all bidders.

### 2. Auto-Bidding System

- [ ] Seller sets maximum price.
- [ ] Optimization: Check max prices and bid once (avoid back-and-forth).
- [ ] Tie-break: Earliest bidder wins if max prices are identical.
- [ ] Only one bidding type implemented (Manual OR Auto).

### 3. Payment & Transaction Flow

- [ ] Post-auction visibility: "Complete Transaction" for Winner/Seller; "Auction Ended" for others.
- [ ] Payment Integration: Stripe, PayPal, Cards, or Bank.
- [ ] Logistics: Winner sends address; Seller confirms with invoice/parcel proof.
- [ ] Completion: Winner confirms delivery; Mark as "Done" to unlock reviews.
- [ ] Cancellation: Seller cancels for non-payment (24hr window) results in -1 for winner.
- [ ] Communication: Direct chat between seller and winner.

---

## 📊 Phase 4: Data & Testing

- [ ] **Data Fakers**
  - [ ] 20+ products across 4-5 categories.
  - [ ] Full descriptions, images, and bidding history (5+ per product).
