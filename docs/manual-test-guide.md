# Phase 1 Manual Test Guide

This guide covers the implemented Phase 1 localhost workflow: authentication, patients, receipts, settings, CSRF protection, database migration, and the Tailwind runtime asset.

## Test Setup

1. Open PowerShell.
2. Start the application:

   ```powershell
   cd D:\Dev\clinic-go\clinic-hcc-app
   go run .\cmd\server
   ```

3. Open <http://127.0.0.1:8080>.
4. Use a dedicated test database. The default path is `clinic-hcc-app\data\clinic.db`.
5. Record defects with the test case ID, steps, actual result, and browser/server error.

To stop the application, press `Ctrl+C`.

## Authentication

### AUTH-01: First-run password setup

1. Start the app with a new database.
2. Visit `/`.
3. Confirm the app redirects to `/setup`.
4. Submit a password shorter than 12 characters.
5. Confirm setup is rejected.
6. Submit two different passwords.
7. Confirm setup is rejected.
8. Submit matching passwords of at least 12 characters.

**Expected:** Password setup succeeds and redirects to `/login`.

### AUTH-02: Login and protected routes

1. Log in with the configured password.
2. Visit `/`, `/patients`, `/receipts`, and `/settings`.
3. Sign out.
4. Visit one of the protected routes again.

**Expected:** Authenticated pages load; after sign-out, protected routes redirect to `/login`.

### AUTH-03: Password reset

1. Stop the app.
2. Run:

   ```powershell
   go run .\cmd\server reset-password
   ```

3. Enter a new password when prompted.
4. Restart the app and log in with the new password.

**Expected:** The old password is rejected and the new password works.

## Patient Workflow

### PAT-01: Create a patient

1. Open `/patients/new`.
2. Enter:
   - Name: `Chan Tai Man`
   - HKID: `A123456(8)`
   - Gender: `Other`
3. Save the patient.

**Expected:** The patient appears in the patient list with the canonical HKID `A123456(8)`.

### PAT-02: HKID normalization

1. Open `/patients/new`.
2. Enter the same HKID as `a 123-456(8)`.
3. Submit the form.

**Expected:** The value is normalized to `A123456(8)`.

### PAT-03: HKID validation and uniqueness

Try each input separately:

- `A123456(7)` — invalid check digit.
- `A12345(8)` — invalid length.
- `A123456(8)` for an existing patient — duplicate.

**Expected:** Each invalid or duplicate patient is rejected with an error.

### PAT-04: Edit patient

1. Edit an existing patient.
2. Change the name or gender.
3. Save.
4. Reopen the patient.

**Expected:** Name and gender change; the HKID field is read-only and unchanged.

## Receipt Workflow

### RCP-01: Create a multi-item draft

1. Open `/receipts/new`.
2. Select a patient and enter a visit date.
3. Add two or more line items, for example:
   - `Tui Na`, quantity `1`, unit price `60000`
   - `Acupuncture`, quantity `2`, unit price `35000`
4. Save the draft.

**Expected:** The draft is saved without a receipt number. All line items, quantities, prices, and their entry order are preserved.

### RCP-02: Draft validation

1. Try saving with no valid line item.
2. Try a blank description.
3. Try quantity `0` or a negative quantity.
4. Try a negative unit price.
5. Try a percentage discount above `100`.
6. Try a fixed discount greater than the subtotal.

**Expected:** The draft is rejected and no invalid receipt is persisted.

### RCP-03: Edit and remove draft items

1. Open a saved draft.
2. Add a line item.
3. Edit an existing line item.
4. Remove one line item.
5. Save and reopen the draft.

**Expected:** The changes persist and at least one valid line item remains.

### RCP-04: Totals

Use:

- Item 1: quantity `1`, unit price `60000`
- Item 2: quantity `2`, unit price `35000`
- Percentage discount: `10`

**Expected:**

- Subtotal: `130000` cents
- Discount: `13000` cents
- Grand total: `117000` cents

### RCP-05: Finalize and immutability

1. Open a valid draft.
2. Finalize it.
3. Confirm a receipt number is assigned.
4. Open the finalized receipt.
5. Try `/receipts/{id}/edit`.
6. Try deleting or updating the finalized receipt.

**Expected:** Finalization succeeds once; the receipt number is present; finalized receipts remain viewable but cannot be edited or deleted.

### RCP-06: Receipt list filters

1. Create receipts with today’s date and another date.
2. Open `/receipts`.
3. Open `/receipts?date=today`.

**Expected:** The filtered page shows only today’s receipts.

## Settings

### SET-01: Update clinic settings

1. Open `/settings`.
2. Change the clinic name, address, phone, or practitioner.
3. Save.
4. Return to the dashboard and receipt pages.

**Expected:** The updated clinic name appears in the shared layout and the values persist after restarting the app.

## Security

### SEC-01: CSRF protection

1. Log in.
2. Use browser developer tools, PowerShell, or another HTTP client.
3. Send a state-changing POST request without `csrf_token` and without the `X-CSRF-Token` header.

**Expected:** The server returns `403 Forbidden` and does not change data.

### SEC-02: Loopback binding

1. Inspect the startup log.
2. Confirm the address is `127.0.0.1:<port>`.
3. Confirm the app is not configured to bind to `0.0.0.0` by default.

**Expected:** The default server is local-only.

## Frontend Asset

### UI-01: Tailwind artifact

1. Confirm `static\css\styles.css` exists.
2. Start the app and open a styled page.
3. In browser developer tools, confirm the stylesheet request to `/static/css/styles.css` succeeds.
4. Optional rebuild:

   ```powershell
   npm run build:css
   ```

**Expected:** The stylesheet builds successfully and the application renders with Tailwind styling.

## Database Migration

### DB-01: Fresh database

1. Start the app with a new database path.
2. Confirm startup completes.
3. Confirm patients, receipts, receipt items, settings, and `schema_migrations` tables exist.

**Expected:** The schema is created automatically.

### DB-02: Legacy database copy

1. Make a backup copy of an existing database before testing.
2. Run the app against the copy.
3. Confirm startup completes without manual schema recreation.
4. Confirm existing finalized receipt numbers remain unchanged.
5. Confirm migrated receipt items retain deterministic order.

**Expected:** The legacy database is upgraded in place and existing data remains available.

## Resetting Test Data

Only do this when the database contains test data:

1. Stop the app.
2. Remove the test database and its SQLite sidecar files:

   ```powershell
   Remove-Item .\data\clinic.db, .\data\clinic.db-shm, .\data\clinic.db-wal -ErrorAction SilentlyContinue
   ```

3. Restart the app and repeat `AUTH-01`.

Do not run this against a production database.

## Intentionally Not Covered in Phase 1

- Automatic or manual receipt archival
- Retention purge jobs
- Reports and exports
- Automatic backups
- Multi-user roles
- Payment, insurance, and appointment integrations
