# Production changes — 2026-08-24

This document records the server-side deployment and security hardening performed with the application release. It intentionally contains no passwords, tokens, hashes, private keys, or database contents.

## Deployed application revisions

- Storefront: `f1bb1371510476ad975225b539c5785536e0eb9e`
- Admin frontend: `d02387ddfad2bdf3d40c1ba0eb37856b35bffe9c`
- Backend: `8b5db5bc3caad55fd2054770e1b17638d1d30ab6`

The repository commits were mirrored on the server and pushed with `[skip ci]` to avoid an unrelated automatic deployment while the live release was being verified.

## Backup and rollback

- A pre-change backup was created under `/root/lebedinski-backups/2026-08-24T1355MSK-pre-hardening`.
- It contains a compressed PostgreSQL logical backup, system configuration archive, application source/build archive, Docker inspection data, and SHA-256 checksums.
- The previous backend binary and frontend build directories were retained in the backup.
- Stopped pre-change database containers were retained temporarily as additional rollback points; the active container keeps the original data volume.

## SSH and host access

- Installed the administrator ED25519 public key and independently verified key login.
- Disabled SSH password and keyboard-interactive authentication.
- Kept root access available by public key only.
- Reduced SSH authentication attempts and login grace time.
- Installed and enabled fail2ban for SSH with increasing ban durations.
- Enabled UFW with default-deny incoming policy and SSH connection limiting.

## Network exposure

- Kept public access to SSH, HTTP, and HTTPS.
- Restricted Zabbix agent traffic to the configured monitoring-provider addresses.
- Bound the backend API to `127.0.0.1:8080` so it is reachable externally only through nginx.
- Bound PostgreSQL's published Docker port to `127.0.0.1:5434`.
- Configured the active PostgreSQL container with `unless-stopped` restart behavior.

## Authentication and data protection

- Activated backend-enforced admin sessions and protected private API routes.
- Verified that customer/order data is rejected without authentication.
- Rotated the root console, admin, and PostgreSQL credentials.
- Set application environment files to mode `600` and the application owner.
- Removed reviewed server-side credentials from the frontend build inputs and bundles.
- Added explicit production CORS origins; arbitrary origins are rejected.

## Deployment validation

- nginx configuration validation succeeded.
- PostgreSQL reports that it is accepting connections.
- Public storefront, admin frontend, and catalog API return `200`.
- Unauthenticated order and promocode reads and item mutations return `401`.
- A complete admin login, session check, logout, and post-logout rejection cycle passed.
- The API process and its `tmux` session were verified active after deployment.

## Items intentionally not changed

- Payment webhook behavior was explicitly excluded.
- Process supervision remains the existing `tmux` arrangement; migration to a systemd service was outside this release.

## Required provider-side follow-up

Credentials previously exposed to browser builds or logs should be regenerated in the CDEK, YooKassa, DaData, and mail-provider dashboards. The application-side exposure paths are fixed, but provider-issued credentials cannot be revoked from the server alone. Any browser-intended map key should be restricted to the production domains in its provider dashboard.

