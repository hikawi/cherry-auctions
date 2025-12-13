---
title: ES01 - Authentication and Authorization
parent: Product and Engineering Specifications
---

# Engineering Spec 01: Authentication and Authorization Flow

## Changelog

Current version: `v2`.

Changes in v2:

- Remapped PKCE authentication into a simpler implementation of JWT key pair to
  keep it simple.

## Summary

This engineering specification concerns the full flow of authentication and
authorization for the app, is subjected to changes at any time during the
iterative process.

To recap, authentication section will concern the process of registering and
logging in, and how to properly save the user credentials as well as what to do
to adhere to OWASP best practices. Authorization section will concern the
process of validating and granting permissions to permitted resources after
logging in succeeded.

## User Flow

High-level user flow for the requirement to use "JWT Access Key / Secret Key"
as a mechanism.

![User Flow Sequence Diagram](../images/authn-authz-user-flow-v1.png)

## Authentication and Authorization

Overview:

- All endpoints must adhere to standards provided by REST API Standards, such as
  GET endpoints are always idempotent and should never have any side effects or
  business operations on the server side.
- All endpoints are named properly after the resources they revolve around.
- Error codes should be returned, with an optional error description, but don't
  prioritize error description as it should integrate with I18n.

### Registration

Each user is registered with a unique email address, with a name that is optional,
as every authentication should be done via the email address (whether retrieved
from Google or inputted directly).

- Endpoint `/v1/auth/register`.
- Accepts `application/json`.
- Must return an `access_token` and a `refresh_token`, both inside HTTP-only cookies
  to protect against XSS or similar attacks.

This format can be changed to a PKCE flow as outlined in the [legacy authentication
authorization document](../legacy/01-authentication-spec.md), but there is currently
no plan to support other devices, such as native, other REST clients or servers.

### Logging in

Each user can login to the system with their email address with the password, or
use one of the OAuth providers.

- Endpoint `/v1/auth/login`.
- Accepts `application/json`.
- Must return an `access_token` and a `refresh_token`, both inside HTTP-only cookies
  to protect against XSS or similar attacks.

### Error Codes

Error codes should be returned in an array form, with a similar structure:

```json
{
  "status": 400,
  "errors": [
    {
      "id": "invalid-username",
      "description": "Optional"
    }
  ]
}
```

This is to have the frontend translate properly on the client-side, instead of
relying on the backend to translate it.
