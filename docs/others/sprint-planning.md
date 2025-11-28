# Sprint Planning

## Summary

This document concerns the planning of each sprint, try to prioritize and setup
the sprint goal, also referred to as an iteration by GitHub Projects.

Planning is done by choosing a set of features to be specified, and to be implemented
from the requirements file. This may or may not be optimal, but it's the best
compromise I have for trying to uphold the software engineering process to a solo
endeavour.

Tasks are split here, but they only contain the brief description of the task, not
the details. For details, look into the issue that was raised for GitHub projects.

## Iteration 1 (Nov 24, 2025 - Dec 7, 2025)

Main Goals:

- Setup a full list of requirements for the project.
- Setup a Foundation (Hello World empty projects) for both backend and frontend.
- Setup a working DevOps pipeline.
- Setup testing if needed (but might not be used)
- Design the database that should be enough for now that adheres to all requirements,
  as well as implementation details for OAuth2.
- Make the following Engineering Specs:
  - Tech Stack used Specification
  - Authentication and Authorization Specification
  - Database Design Specification
- Setup a Docker Compose stack for fast spin-up for development purposes.
- Setup the full authentication flow with basic email/password login using PKCE.
  - OAuth2 Clients and Authentication
  - Access Token Grants
  - Refresh and Refresh Token Rotation
  - Logout and Invalidation

Out of Scope (Optional if every main goal has been achieved):

- Middlewares to integrate with OAuth2 implementation.
- Test-driven Development with `go test` or `vitest`.

Tasks split:

1. SPEC: Rewrite the list of requirements for the project.
2. DES: Design the basic database system that will be used to be bootstrapped with.
3. SPEC: Write a formal specification on each library used that will satisfy the
   project requirements.
4. CI: Setup a GitHub Actions pipeline that works with build-image-deploy flow.
5. SPEC: Write up the final sprint planning for sprint 1.
6. SPEC: Write up the authentication flow engineering specification.
7. SPEC: Design and setup the project's architectures.
8. BE: Setup an empty minimal Golang backend.
9. FE: Setup an empty minimal Vue frontend.
10. BE: Implement GORM models for the authentication flow.
11. BE: Implement register/login flow with JWT pair and PKCE.
12. CHORE: Setup a Docker Compose setup for easy setting-up a local development
    environment.
