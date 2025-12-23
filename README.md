# Cherry Auctions

![wakatime](https://wakatime.com/badge/user/16f7181f-8df5-44c7-b2b4-4fa68b0c2dfe/project/9c6372ce-b005-493b-95b1-60ee50480cc0.svg)

## Overview

This is a project assignment for the course of Advanced Web App Development, or
AWAD for short. There are a list of expected requirements, although no restrictions
on technologies, but the selected technologies must be able to carry out the same
functions.

## Project Overview

**Disclaimer**: This is subjected to be changed at any time during the iterative
process.

This project is a monorepo where each service is contained inside a folder, and
should have no knowledge of what is outside it. All services are laid out flat
in the repository's root directory as I do not want to nest services into a
`backend` folder, in the case of this spiraling out to microservices.

For deployment specifications, see [`specs/03-architecture-design`](./docs/specs/03-architecture-design.md).

- Frontend uses an SPA using Vue and Vite.
- Backend current has these following services (and might have more due to the
  event-based design of an auctions system):
  - CherryAuctions Service: The main backend that holds all the resources, and
    is the one that the frontend should talk to the most.
- Project Documentation (not Swagger) is setup with Jekyll.

All services that need it will have a Docker image built and published for easy
deployment.

Relevant Links:

- Frontend at [cherry-auctions.luny.dev](https://cherry-auctions.luny.dev)
- Backend at [api.cherry-auctions.luny.dev](https://api.cherry-auctions.luny.dev)
- Project Specifications at [docs.cherry-auctions.luny.dev](https://docs.cherry-auctions.luny.dev)

## Cherry Auctions – Services Overview

| Service        | Image                    | Ports                                | Description / Usage                                                                                                            |
| -------------- | ------------------------ | ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------ |
| **Mailpit**    | `axllent/mailpit:latest` | `1025` (SMTP)<br>`8025` (Web UI)     | Local SMTP server for development. Captures outgoing emails so they can be viewed in a browser without sending real emails.    |
| **RustFS**     | `rustfs/rustfs:latest`   | `9000` (S3 API)<br>`9001` (Admin UI) | S3-compatible object storage used for file uploads such as images and attachments. Acts as a local replacement for AWS S3.     |
| **PostgreSQL** | `postgres:18-alpine`     | `5432` (Database)                    | Primary relational database storing application data such as users, auctions, bids, and transactions.                          |
| **Frontend**   | `node:24-alpine`         | `5173` (Dev Server)                  | Frontend development server using PNPM with hot reloading (e.g. Vite). Serves the web UI during development.                   |
| **Backend**    | `cosmtrek/air`           | `3000` (Host) → `80` (Container)     | Backend API server with live reload. Handles authentication, business logic, database access, file uploads, and email sending. |

## Getting Started

This monorepo is orchestrated by Docker Compose. Although you can `cd` into each
directory and spin up the service natively yourself (by using Go compiler or
NodeJS for example), it's recommended that you setup Docker & Compose as it should
setup everything for you.

Running Docker Compose, you may see messages like `environment variable not set,
defaulting to blank string`. This means you should populate an `.env` file for the
current directory.

```bash
docker compose up --build
```

Otherwise, you can install the following for this project:

- OpenSearch
- OpenSearch Dashboards (expected to be on port 5673)
- Fluentbit (expected to be on port 2020)
- Go compiler, at least version 1.25.4.
- NodeJS LTS Krypton (v24) and NPM v11 or PNPM v10.
- Postgres (or Postgis, or any Postgres-compatible database), expected to be on
  port 5432.
- S3-compliant database (I expected to use RustFS, migrating away from MinIO),
  but you can use any other like Ceph or Bunny Storage.
- SMTP server.
- Air by Cosmtrek for hot-reloading Go modules.

**Note**: Some stuff won't be installed by Docker, as it's unnecessary to do so.
Playwright Frameworks for Frontend Testing for example will not be pulled in by
Docker. Refer to their documentation to get yourself setup with E2E.

## Documentation

I try to document my project as clearly as possible, even if this is a solo endeavour.
I try to adhere to the software engineering process: Requirements ->
Specifications -> Design -> Implementation -> Testing, but still stay iterative
for Agile.

Some may call this stupid, waste of time, but how would I have any grounds to
say I want to become a Project Manager or a Technical Lead if I can't write proper
specifications, right?

## License

- **Source code** (including all services under `backend`, `frontend`, and
  related scripts) is licensed under the **Apache License 2.0**.
- **Documentation** (all files under `docs/`) is licensed under the **Creative
  Commons Attribution-ShareAlike 4.0 International License (CC BY-SA 4.0)**.

See the respective `LICENSE` files for full details.
