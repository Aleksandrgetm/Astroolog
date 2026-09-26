# Booking bonus and messenger flow — 2026-09-26

Implemented for the booking form; the question form retains its existing flow.

- Active main readings remain separate radio options. `online-meeting-40` is an optional checkbox after the four readings, with its own API price and `is_addon=true`. It cannot be submitted as a main option.
- Backend validates active service, main option, addon ownership and both quoted prices. It calculates `total_price` from DB prices. `price` remains the base-price snapshot; bonus code/title/price are stored separately.
- `booking-bonus-20260926-v1` is an additive, versioned migration. The historical €70 option remains stored but inactive. Historical bookings and dates are not rewritten. Bootstrap descriptions are replaced only on exact matches; subsequent starts preserve administrative changes.
- Date input and date validation are removed from the booking form. New bookings store NULL for the legacy date column.
- WhatsApp/Telegram selection validates and saves first. Messages use the response's saved snapshot. The success panel confirms saving only. A retry link does not submit again; switching messenger with the same payload reuses its idempotency key. Session storage contains only a payload digest and random key.
- A blank popup is reserved before the POST to preserve browser activation; it is closed on save failure. Navigation to the messenger occurs only after saving. Blocked popup or URL failure keeps the saved state and provides a retry.
- RU/LV/EN strings, service option displays, direction links and privacy copy are updated. Existing routes, canonical URLs and hreflang are unchanged.

## Verification

- TypeScript, Vite and prerender: 42 content pages + 3 localized 404 pages.
- `test:i18n`, `test:whatsapp`, `test:pre-admin`, `test:seo`, `test:directions`.
- Existing interaction and responsive scripts: 210 route/viewport checks, navigation, language switching, retries, HTTP 301/404.
- New `test:booking`: 390/430/768/1024/1440, RU/LV/EN, both messenger buttons, no date payload, totals, saved response reuse, backend error, stale quote, popup failure, URL failure and manipulated total. Messenger navigation is intercepted; no external messages are sent. This does not claim native WhatsApp/Telegram delivery testing.
- `go test ./...`, `go vet ./...`, plus integration tests on an empty database and a copy of the local database.
- Historical booking checksum before/after migration matched. Admin-modified bonus price survives repeated migration; snapshots retain their original prices.
- Local working DB migration applied after a backup at `/tmp/astroolog-before-bonus-20260926.sql`. Integration/browser test submissions use isolated databases, not the working DB.

Production deployment was not performed. The generated build uses the existing test origin `https://example.invalid`; a deployment must rebuild with the real `VITE_SITE_URL` as documented in README.
