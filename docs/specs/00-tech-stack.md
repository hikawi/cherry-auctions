---
title: ES00 - Tech Stack
parent: Product and Engineering Specifications
---

# Engineering Spec 00: Project Stack

## Changelog

Current Version: `v1`

- `v1`:
  - Initial version.

## Summary

This specification should outline and go into detail about this project's chosen
tech stack, or related systems that might be of use throughout the development
process.

## Technology Stack

Since this is a solo endeavour, I would be the one to select the tech stack myself.
So, in accordance with what I have registered with the professor, this project will
concern the following:

- Frontend System:
  - Vue 3 (as registered with the professor at the time of course signup, and the
    requirement for "SPA")
  - Vue Router (as needed for the technical requirement "Use a router service")
  - Pinia (central state management per the technical requirement "Use a state
    management system and have Flux data flow")
  - TailwindCSS (as per the requirement to "Use a CSS Framework like Bootstrap,
    Tailwind or Material Libraries")
  - Zod (?) (as per the requirement to "Use a form validation library")
- Backend System:
  - Golang with GIN (as registered)
  - Swaggo and Gin Swagger (to generate Swagger API Docs, as per the requirement
    "Have Swagger docs on all endpoints")
  - GIN's built-in go-validator (to validate requests, as per the requirement to
    "have proper validation at needed endpoints")
  - Go's JWT v5 Implementation (to use JWTs, as per the requirement to use "JWT
    Access Token Refresh Token pair for authentication and authorization")
  - GORM: ORM solution for Golang.
- Infrastructure:
  - OpenSearch: Searching service to search logs, as per required by the "ELK or
    similar stack" requirement.
  - OpenSearch Dashboards: The dashboard for accessing OpenSearch in a web-based
    UI, per the "ELK or similar stack" requirement.
  - FluentBit: The log aggregator and sender for sending logs to OpenSearch, per
    the "ELK or similar stack" requirement.
  - Docker & Docker Compose: Containerization and Isolation tool for setting up
    the local development infrastructure, as well as the deployment environment.
  - Postgres 18: Chosen as the database of choice.
  - MinIO (?): self-hosted S3-compliant storage for images, videos and assets.
- Third-party Integrations:
  - Google reCaptcha: per the requirement "Have a reCaptcha for register page".
  - OAuth2 Integration: per the optional requirement "Have Login with Google or
    Login with GitHub options"
  - Payment Integration: per the requirement "Have a payment gateway with Stripe
    or Paypal or similar platforms".
  - PurelyMail: the smtp server host per the requirement to "Have a mailing system".
  - BunnyCDN (?): volume-based and budget CDN to integrate with S3-compliant storage.
- Other Tools:
  - Figma: for designing the webpages.
  - PlantUML: for drawing UML diagrams, such as the database design diagram.
  - Vercel: for hosting the frontend on their Edge Network.

Dependencies, libraries or technologies marked as `(?)` meant not sure, and has a
high chance of being omitted from the project or changed entirely.

## CI/CD

CI/CD is carried out by GitHub Actions, going through the following pipelines:

1. Commit Lint
2. Code Quality Lint
3. Branch Testing (T.B.D., as this may delay the actual product)
4. Product Building
5. Image Building & Pushing
6. Deployment

If the previous step fails, the pipeline stops and should block the merge option,
as configured in GitHub's branch protection rules.

Due to this being a solo endeavour, a staging environment might be too daunting to
setup, so only a production environment is setup at: `cherryauctions.luny.dev`,
reliant on `s3.luny.dev`.
