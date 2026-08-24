# Change log

## 2026-08-24 — API authorization and secret handling

Commit: `8b5db5bc3caad55fd2054770e1b17638d1d30ab6`

### Added

- Server-side admin login, logout, and session endpoints.
- Opaque, time-limited admin sessions delivered through a Secure, HttpOnly, SameSite=Strict cookie.
- Per-IP login throttling with five failed attempts allowed per fifteen-minute window.
- Backend proxy endpoints for DaData address suggestions and delivery lookups, including request limits, reply limits, rate limiting, and outbound timeouts.
- Automated tests covering authentication, protected routes, session invalidation, and CORS configuration.
- A documented environment-variable template without live credentials.

### Changed

- Protected order reads, order mutations, promocode administration, item mutations, banner uploads, and catalog administration behind admin authorization middleware.
- Kept public catalog reads available without authentication.
- Replaced wildcard CORS behavior with an explicit production-origin allowlist.
- Made the API bind address configurable so production can listen on loopback only.
- Removed sensitive CDEK credential and token-response logging.

### Validation

- `go test ./...` passed.
- `go vet ./...` passed.
- Public catalog requests return `200`; unauthenticated private reads and mutations return `401`.
- Valid admin-origin preflights succeed and unapproved origins return `403`.

### Deliberately unchanged

- Payment webhook behavior was excluded from this release by request.

See `docs/PRODUCTION_CHANGES_2026-08-24.md` for the accompanying server changes.

