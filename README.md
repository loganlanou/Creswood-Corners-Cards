# Creswood Corners Cards

Go-based trading card storefront modeled around the structure of `logans3dcreations.com`, but adapted for football card sales, live selling, customer accounts, and admin operations.

## Stack

- Go `net/http`
- Server-rendered `html/template`
- In-memory preview data store
- Cookie sessions with local account registration/login
- Local demo checkout flow with Stripe-ready config placeholders

## What Is Implemented

- Homepage with reference-inspired hero, featured products, and service/value sections
- Shop page with searchable product grid
- Product detail pages
- Live selling page and active live banner
- Local cart stored in browser `localStorage`
- Demo checkout that creates tracked in-memory orders for preview testing
- Account signup, signin, and session-based auth
- Admin dashboard for:
  - live session settings
  - product creation
  - registered account visibility
  - order payment / fulfillment / shipping updates

## Current Auth / Payments State

- Local auth is fully implemented in Go.
- Clerk environment variables are included, but Clerk itself is not yet wired into this Go app.
- Stripe environment variables are included, but checkout is currently local demo mode unless Stripe integration is added next.

That means the application is fully runnable on localhost now, but Clerk and Stripe remain integration points rather than completed live integrations.

## Run Locally

1. Review and copy values from `.env.example` into your environment if needed.

2. Start the server from the repo root:

```bash
env GOCACHE=/tmp/go-build go run ./cmd/server
```

3. Open:

```text
http://localhost:3000
```

For verification in this session, I launched isolated localhost instances on `127.0.0.1:3013`.

## Default Bootstrap Admin

The app creates a bootstrap admin user on first run:

- Email: value of `ADMIN_BOOTSTRAP_EMAIL`
- Password: value of `ADMIN_BOOTSTRAP_PASSWORD`

Default example values are:

```text
owner@example.com
changeme123
```

Change those before real usage.

## Build Verification

Validated with:

```bash
env GOCACHE=/tmp/go-build go build ./...
```

Functional localhost checks completed against the Go app:

- `/` returned `200` and rendered the live banner plus featured products
- `/shop` rendered catalog HTML
- `/live` rendered the live-selling page
- `/admin` redirected to `/sign-in` when unauthenticated
- bootstrap admin sign-in succeeded
- authenticated `/admin` rendered inventory, stream control, registered accounts, and orders
- demo checkout returned a success redirect and created an order

## Vercel Preparation

The repo now includes [vercel.json](/home/loganlanou/Creswood-Corners-Cards/vercel.json) and [api/index.go](/home/loganlanou/Creswood-Corners-Cards/api/index.go) so Vercel can route all requests through the Go handler.

Use these import settings:

- Application Preset: `Other`
- Root Directory: `./`
- Build Command: leave empty
- Output Directory: leave empty

Before deploying:

1. Set these environment variables in Vercel:

```text
APP_URL
SESSION_SECRET
ADMIN_BOOTSTRAP_EMAIL
ADMIN_BOOTSTRAP_PASSWORD
ADMIN_EMAILS
DEMO_CHECKOUT
```

2. Decide on data persistence:

- Current preview data is in memory and resets when the function instance resets.
- Before a real public launch, move the Go app to a persistent database such as Postgres or Turso/libSQL.

3. If you want real auth and payments before launch:

- replace local auth with Clerk integration
- replace demo checkout with Stripe checkout + webhook processing

## Notes

- The old Next.js implementation has been removed so Vercel does not detect or build the wrong app.
- The Go application is the only active implementation now.
- If you want, the next pass can do two focused follow-ups:
  1. wire real Clerk and Stripe integrations into the Go app
  2. move data persistence from local SQLite to a production database
