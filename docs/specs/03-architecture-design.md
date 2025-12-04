# Engineering Spec 03: Architecture Design

## Summary

This document concerns the logical overview of the application's architecture,
including networking.

In scope topics:

- Deployment Tactics
- Project Structure
- Networking and DNS

## Project Structure

To keep it simple for the course, this uses a simple monolith Go application,
but split to be modular enough to be easy to maintain, as well as fast enough
to iterate with.

Project Stack:

- Frontend (SPA)
- Backend (Single Binary)
- Database (Local PostgreSQL)
- S3 Storage (Local RustFS)

## Deployment Tactics

Frontend is deployed on Vercel to utilize their Edge Network and built-in CDN for
best latency to any user, while also seemingly to mitigate the latency from the
user to the backend hosted on [Netcup](https://netcup.com)
