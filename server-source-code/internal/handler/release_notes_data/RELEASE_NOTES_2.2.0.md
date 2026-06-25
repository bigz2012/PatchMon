# PatchMon V2.2.0 Release Notes

Integrates vetted community pull requests from upstream on top of 2.1.0. No
database schema changes (still migration 000041).

## New Features

- **Alpine Linux live patching via APK** (PR #863) — Alpine hosts can now be
  patched through PatchMon like apt/dnf/pacman hosts.
- **Bulk patch action for selected hosts** (PR #797) — patch multiple hosts at
  once from the Hosts page.
- **Host Security Update Status chart** (PR #801) — new dashboard chart showing
  hosts with security updates vs regular updates vs up to date.
- **Arch reboot-required detection** (PR #554) — agents on Arch detect the
  latest installed kernel via pacman for reboot-required status.

## Bug Fixes

- **dnf makecache hang** (PR #854) — pass `-y` so GPG key import prompts don't
  hang the patch run.
- **Debian deb13.x kernel comparison** (PR #818) — correct reboot-required
  detection for Debian 13 kernel package versions.
- **RHEL with subscription** (PR #835) — handle the dnf "Updating Subscription"
  output so update parsing isn't thrown off.
- **Compliance scanner unavailable** (PR #862) — reject compliance scans up
  front when the scanner isn't available instead of letting them fail mid-run.
- **REDIS_USER in queue client** (PR #785) — the asynq queue now honours
  `REDIS_USER` for Redis ACL setups.
- **Agent enrollment** — better curl error reporting (PR #617), quoted Makefile
  paths (PR #827).

## Docs & Dependencies

- Docs fixes (PR #823, PR #792).
- Bump `gopsutil` to v4.26.5 in the agent (PR #812).

All changes credit their original upstream authors.
