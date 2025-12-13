---
title: Project Requirements
parent: Product and Engineering Specifications
---

# Cherry Auctions

## Overview

This is a project assignment for the course of Advanced Web App Development, or
AWAD for short. There are a list of expected requirements, although no restrictions
on technologies, but the selected technologies must be able to carry out the same
functions.

## Roles

### Guest

1. **Menu System**: The guest may view on the website a menu of categories of
   items. Categories must support at least two levels, such as Electronics ->
   Mobile Devices, and Electronics -> Laptops.
1. **Home Page**: The guest may view a home page that can display "Top 5 Ending
   Soon Auctions", "Top 5 Auctions (most number of bids)" and "Top 5 Highest Auctions".
1. **Product Catalog**: The guest may view a catalog of products, can filter
   with categories, and most importantly **paginated**.
1. **Search**: The guest may make a search for categories or products. They may
   search using a **fuzzy search** with a **full-text search**. The results must
   be **paginated** and **sorted** (can be chosen, prices ascending/descending or
   ending time). Recently posted products should also have an eye-catching display.
1. **Product Display**: Products displayed on catalogues should have at least
   the following attributes: thumbnails, names, current bid, current highest bidder,
   buy-it-now price (if available), posted date, remaining time, current bids count.
1. **Product Details**: Hero photos, additional photos, product names, prices,
   current bid, buy-it-now price, seller infos and reviews, highest bidder infos
   and reviews, posted timestamp, ending timestamp, detailed description.
   1. If the ending time is in 3 days or fewer, then make it display relative to
      the user's time.
   2. Questions and Answers between bidders and seller section.
   3. 5 similar products in the same category.
1. **Registration**: Register an account using full names, address and emails.
   1. Verification with OTP is needed.
   2. reCaptcha
   3. Use a hashing algorithm to secure the password. **DO NOT SAVE THE PASSWORD
      IN PLAIN TEXT**.

### Bidder

1. **Watch List**: A user can save an item to their watch list, from the catalog
   or the product detail page.
2. **Bid**: Can bid at the product detail page. The system disallows users with
   less than 80% rating to bid.
   1. Ratings are calculated based on "likes" and "dislikes", e.g. 8 likes and
      2 dislikes means 8/10 total, 80% rating. If the seller allows
      it, new accounts can also bid.
   2. The system sets a minimum bid to "step up". Current bid + The step (% or
      absolute value) set by the seller.
   3. The system must have a confirmation dialog.
3. **History of Bids**: Bids can be viewed as a card-based or tabular view.
   Names should be masked. Must also contain timestamp, name, and price.
4. **QnAs**: The Product Detail view should have a QnA section between bidders
   and seller. When the seller or bidder gets a question or a response, an email
   must be sent.
5. **Profiles**:
   1. The email, name, address, password can be changed.
   2. View my own rating, details of each rating + comments of those rating given
      by others.
   3. My favorites (or marked products in watch list)
   4. My current bids
   5. My won bids (along with a rating from the seller and a comment): For example,
      a successful transaction and the seller is happy. The bidder runs away and
      the seller is mad.
6. **Requesting Selling Privileges**: After 7 days of usage, a bidder can ask
   for selling permissions (or upgrading the role). Administrators can use ratings
   to evaluate and choose to approve or reject.

### Seller

1. **Post a Product**: Must include all information: product name, at least 3
   photos, starting bid, step, buy-it-now price (optional), and description.
   1. The product description must support a WYSIWYG text input.
   2. If there's a new bid when it's 5 minutes left, should it auto-extend
      the timer? The seller can choose whether to auto-extend.
   3. The administrators can choose when should it extend, and how much to
      extend by, and this setting applies to all products.
2. **Edit a Product**: Can edit the description by making P.S. remarks. The
   new changes are appended to the description.
3. **Deny a Bidder**: In the product detail view, the seller can deny any
   bidders. The denied bidders can no longer bid on the product. If the denied
   bidder is the current highest, the bid is transferred to the second highest.
4. **Answer a Question**: Can answer bidders' questions in QnA section.
5. **Profiles (More Advanced)**:
   1. Similar to Bidder's Profiles.
   2. My current products.
   3. Products that are won. Can review the top bidder a (+1) or a (-1), along
      with a feedback comment. Can stop the transaction, which automatically (-1)
      the top bidder, with an automated comment "The Top Bid did not follow through".

> I'm not sure about 5.3's last condition.

### Administrators

> For the following section, when "Manage" is referred, it means:
>
> - View list
> - View Details
> - Add or Create
> - Delete
> - Update
> - More specialized tools specific to the category of management

1. **Category Management**: If there is a product in a category, that category
   can not be deleted.
2. **Product Management**: Can remove any products off the platform.
3. **Users Management**: Can see who is requesting the upgrade from bidder to
   seller. Can approve/reject the requests.
4. **Dashboard**: Analytics, Charts, Summary about the users count, selling
   requests, financials, bids count, etc. You may suggest your own feature
   to track and analyze.

### Shared Functions

1. **Login**: Optional, but encouraged to also provide OAuth2 options such as
   Google, Facebook, X, GitHub.
2. **Forgot Password**: Can change password through an OTP code sent by an email.

## Systems

### Mailing System

With each transaction, or an important action, the system must send an email to
the relevant parties, such as:

- Successful bid: to the bidder, to the seller, and to the previous top bidder.
- Rejected bid: to the bidder.
- Expired auction (no bidders): to the seller.
- Ended auction: to the seller, to the winner.
- Question asked: to the seller.
- Question answered: to all bidders and the one who asked the question.

### Auto-Bidding System

- The seller may set a **maximum price** to pay for a product.
- If a new bid comes in, it may automatically bid again. You may optimize this
  by checking maximum prices between bidders and bid once, instead of back and
  forth.
- If 2 bidders have the same set price, the one that bid first would be set as the
  top bidder.

Example for a product, starting bid at \$10, step is \$0.1.

Manual bidding:

1. A bids $10. A is top bid.
2. B bids $10.1. B is top bid.
3. C bids $10.5. C is top bid.
4. A bids $10.8. A is top bid.
5. C bids $11. C is top bid.

Auto-bidding:

1. A's max bid is \$11. A bids \$10. A is top bid.
2. B's max bid is \$10.8. B bids. A auto-overrides at \$10.8. A is top bid.
3. C's max bid is \$11.5. C bids \$11.1. C is top bid.
4. D's max bid is \$11.5. D bids. C auto-overrides at \$11.5. C is top bid.
5. D's max bid changed to \$11.7. D auto bids at \$11.6. D is top bid.

As you can see, the bidders don't have to interactive continuously to bid. The
auto-bidding should be that the price is only barely enough to win.

> "The system only implements one type. Do not implement both"
>
> This sentence is quite unclear.

## Payments

After an auction, the winner and the seller can see the "Complete Transaction"
button from the product detail page, but other users can only see "Auction Ended".

The process is as followed:

1. The winner makes the transaction. Integrate with Stripe, PayPal, Cards or bank
   here.
2. The winner sends along the address.
3. The seller confirms the transaction has followed through, and sends the package.
   Confirm with the winner with a parcel send invoice, or similar items.
4. The winner confirms the delivery.
5. The transaction is marked as "Done" and reviews can now be given.

In this process, the seller may **cancel**, which automatically sets a (-1) for the
winner. For cases: Seller requests payment in 24 hours, but did not receive, cancel.

The seller and the winner both have access to a chat to communicate during the process.

Afterwards, both can change or modify the reviews at any time.

## Other Requirements

- Front-end: Must be a web application, using an SPA.
  - Use a router service of the SPA. For example, React Router for React.
  - Form Processing & Validation. For example, Angular Forms.
  - Managed State (optional): Redux, Pinia,...
  - **CONFORM TO ONE DESIGN SYSTEM AND CONSISTENT THROUGH ALL SCREENS.**
- Back-end: At least some endpoints must follow a RESTful API standard.
  - Full validation on endpoints.
  - Swagger Documentation.
  - Logging and monitoring with Grafana or ELK or similar technology stacks.
  - Security using JWTs AccessToken and RefreshTokens.
- Complete exactly the functions required. Don't omit. Although, you may add more
  features for UX of each function.
- Data Fakers: Must have at least 20 products spread around in 4-5 categories, with
  full descriptions and images. All products must have a bidding history, each at
  least 5 bids.
- Source Management: All work must be tracked meaningfully through GitHub.
  - No source management, or only few commits such as "Upload backend files" or
    "Upload frontend files" => Non-negotiable fail of this course.
