---
title: Database Design
parent: Product and Engineering Specifications
last_modified_date: 2025-12-28
---

# Engineering Spec 02: Database Design

## Summary

This document concerns the design of the database throughout the project.

In scope:

- Relations and what to save, how to save, data types to use, or where to index
  and where not to.
- Changelogs for each version of the database.

Out of Scope:

- Row Level Security

## Database Design

![Database Design](../images/db-design-v4.png)

## Changelogs

### Version 5 (Current)

- Removed the `role` field in the `users` table.
- Add new table `roles`.

### Version 4

- Removed `oauth_clients` due to PKCE being overkill for this project, I was stupid.

### Version 3

- Added an `oauth_clients` table as the requirement "JWT key pair" does not specify
  which implementation, so I chose to go with PKCE as suggested by Auth0, and this
  requires registering with the backend a public client for the SPA.

### Version 2

- Added a `refresh_tokens` table as per the requirement to use a "JWT key pair"
  that requires refresh tokens to be persisted somewhere.

### Version 1

Initial version. No changes.
