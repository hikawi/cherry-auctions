---
title: Sprint Retrospective
parent: Other Specifications
---

# Sprint Retrospective

## Summary

This document aims to provide insights, statistics, tasks completed throughout
the iterations, and the solo author's work issues to get better for each next
sprint.

Retrospective should be on each sprint's second week's Sunday.

## Iteration 1 (Nov 24, 2025 - Dec 7, 2025)

What was finished (check list against previous iteration planning):

- [x] Setup a full list of requirements for the project.
- [x] Setup a Foundation (Hello World empty projects) for both backend and frontend.
- [x] Setup a working DevOps pipeline.
- [x] Setup testing if needed (but might not be used)
- [x] Design the database that should be enough for now that adheres to all requirements,
      as well as implementation details for OAuth2.
- [x] Make the following Engineering Specs:
  - [x] Tech Stack used Specification
  - [x] Authentication and Authorization Specification
  - [x] Database Design Specification
- [x] Setup a Docker Compose stack for fast spin-up for development purposes.
- [ ] Setup the full authentication flow with basic email/password login using PKCE.
  - [ ] OAuth2 Clients and Authentication
  - [ ] Access Token Grants
  - [ ] Refresh and Refresh Token Rotation
  - [ ] Logout and Invalidation

What was finished (listed):

1. Full DevOps flow for linting, checking and deploying both frontend and backend.
2. GHA integration for documentation with Jekyll.
3. Setup frontend with Vue 3, Vite, Vitest and Playwright.
4. Setup backend with Go, GORM and Swagger.
5. Setup infrastructure for local development with `docker compose`.
6. Finished a crude login page (Frontend).
7. Finished authentication-related routes (Backend).
8. Wrote documentation on technical stack used.
9. Wrote documentation on PKCE Oauth2-compliant authentication flow (archived).
10. Wrote documentation on the database design.
11. Design the database crudely on the current project requirements.
12. Setup GORM models and auto migration on backend side.
13. Setup the design system for the project.
14. Setup TailwindCSS and Vue I18n for frontend.
15. All changes are deployed to production after rigorous linting and checking.

What was not finished:

1. Full PKCE Authentication flow with Oauth2-compliance (deemed too complicated
   for a school project).
2. Proper linking frontend and backend for logging in and registering in a way
   that is visible and demonstrateable.
3. Full flow authentication and authorization flow specification.

Issues:

1. Playwright was not working on Arch Linux, see [Common Issues](./common-issues.md).
