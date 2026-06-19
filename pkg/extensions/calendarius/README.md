# Calendarius

The calendar module (happenings, slots, scheduling).

The module's persisted extension id is `calendarius` (see `const4calendarius/module_id.go`),
so its data lives under `/spaces/{spaceID}/ext/calendarius/...`.

## Why `calendarius`?

The name follows the Sneat house style: an English root plus a Latinate suffix that ends in `-us`
(the "us / we" family pun). `calendarius` keeps the descriptive **calendar** root and matches the
existing `-ius` coinages **yardius** (Yard → Yardius) and **companius** (Company → Companius).

The two previous names were retired because each broke the convention:

- `calendarium` (Go backend) — grammatically fine Latin, but the lone `-ium`/`-m` outlier in the
  ecosystem; it does not end in `-us`.
- `schedulus` (Angular frontend) — ends in `-us`, but loses the *calendar* root.

Brand: [calendarius.app](https://calendarius.app). See
`backstage/marketing/brand-strategy/module-naming-and-domains.md` for the full naming rationale.
