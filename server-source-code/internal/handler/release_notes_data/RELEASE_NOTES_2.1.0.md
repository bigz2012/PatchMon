# PatchMon V2.1.0 Release Notes

Feature release adding user-facing preferences and integration improvements.

## New Features

- **Per-user timezone preference (#786):** Each user can now choose their
  timezone under Profile → Profile Information. Dates and times across the app
  render in the selected timezone; leaving it on "Browser default" uses the
  browser's locale. No longer restricted to the server-side `TZ`/`TIMEZONE`
  environment variable.

- **"Reboot required" in the Homepage integration (#829):** The gethomepage
  widget stats endpoint now includes `hosts_needing_reboot`, so you can surface
  how many hosts are pending a reboot after patching directly on your dashboard.

## Notes

- Notifications for new security updates (#811) are already available: enable the
  **Host security updates exceeded** alert (default threshold 1) under alerting
  to be notified when a host has pending security updates.

## Upgrade

This release adds a database migration (`000041`) that adds a nullable
`timezone` column to the `users` table. It is applied automatically on startup
and is backward compatible.
