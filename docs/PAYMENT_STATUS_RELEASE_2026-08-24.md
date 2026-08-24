# Payment status release — 2026-08-24

This document records the application and production-server changes for the paid-order status fix. It contains no passwords, tokens, private keys, payment identifiers, or customer details.

## Application revisions

- Backend: `773cca22e2acd83099c68391ff5ceb8f128fb195`
- Admin frontend: `3bf63af6f1d54338595d65f1159a02c05a9d1f5f`
- Storefront: unchanged by this release.

The application commits were pushed with `[skip ci]`; the tested builds were deployed deliberately from the matching server checkouts.

## Backend behavior

- YooKassa is now the authoritative payment source: the webhook payment is fetched and verified before local state changes.
- A successful payment is persisted before CDEK work begins.
- Payment and fulfillment are stored separately as `payment_status`, `fulfillment_status`, and `fulfillment_error`.
- The exact `payment_id` row is used for state changes, protecting duplicate cart records.
- Fulfillment is claimed atomically, making repeated webhook delivery idempotent.
- CDEK processing runs outside the webhook response. A CDEK rejection records `Needs Attention` while the order remains `Paid`.
- Raw provider bodies and customer data are not stored in the warning message.

## Admin frontend behavior

- Affected paid orders display a warning indicator beside `Paid`.
- Clicking the status opens a modal that explains payment and delivery states separately.
- The modal provides copy-ID and order-detail actions for the existing support-chat workflow.
- Automatic shipment retry is intentionally absent until a provider-safe retry design can guarantee that no CDEK shipment already exists.

## Production server changes

- Pre-release backup: `/root/lebedinski-backups/2026-08-24T2153MSK-pre-payment-status`.
- The backup contains a PostgreSQL custom-format dump, prior backend binary, prior admin build, backend environment file, nginx API configuration, and SHA-256 checksums.
- Added the three payment/fulfillment columns using an additive PostgreSQL transaction.
- Deployed backend binary: `/home/github/lebedinski/bin/lebedinski-20260824-payment-status`.
- Rebuilt the admin frontend in `/home/github/Lebedinski_Admin/dist`.
- Changed the nginx API upstream from `localhost:8080` to `127.0.0.1:8080` to match the backend's IPv4-only loopback bind and remove intermittent IPv6 `::1` connection failures.

## Data reconciliation

- Rechecked orders 1788 and 1816 directly against YooKassa after deployment.
- Each had exactly one succeeded, paid, unrefunded payment with the expected order description, currency, and amount.
- Updated only those exact payment rows to `Paid / Succeeded / Needs Attention`.
- No CDEK request, email, inventory change, or promocode change was triggered by reconciliation.
- Ambiguous historical duplicate records were deliberately left unchanged.

## Verification

- Backend unit tests and `go vet` passed locally and on the server.
- Admin production builds passed locally and on the server.
- nginx configuration validation passed before reload.
- Public API and admin HTTPS checks returned `200` after deployment.
- The backend process is running from the new binary in the existing `tmux` session.
- A live authenticated browser check confirmed that order 1816 displays `Paid` with the warning indicator and opens the support explanation modal.
- The browser test session was logged out after verification.

## Rollback outline

If rollback is required, restore the previous backend binary, admin build, and nginx configuration from the backup directory. The database columns are additive and are safe for the previous binary to ignore; the database dump is retained for full data rollback if separately approved.
