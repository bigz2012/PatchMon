# PatchMon V2.0.3 Release Notes

Maintenance release with notification-delivery and agent-installation fixes.

## Bug Fixes

- **SMTP relays without TLS (#817):** Authenticated SMTP over a plaintext
  connection no longer fails with `unencrypted connection`. When TLS is
  explicitly disabled (e.g. a trusted internal relay), credentials are now
  sent over the plaintext connection as configured.

- **Email line length (#845):** Alert emails (notably `host_down`) are now
  encoded as quoted-printable, so no line exceeds the RFC 5322 998-octet
  limit. Strict mail servers no longer reject these messages.

- **Mattermost webhooks (#851):** Mattermost incoming webhooks
  (`/hooks/<token>`) now receive a Slack-compatible payload instead of the
  generic JSON body, fixing the `HTTP 400` rejection.

- **Arch / Manjaro agent install (#850):** The installer now pulls in
  `pacman-contrib` (provides `checkupdates`) and `fakeroot` on pacman-based
  systems, so the first system report succeeds and adding the host completes.
