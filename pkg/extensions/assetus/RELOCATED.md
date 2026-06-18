# ⚠️ Assetus has a new canonical home: `sneat-co/assetus`

Active Assetus development now happens in the dedicated repository
**[`github.com/sneat-co/assetus`](https://github.com/sneat-co/assetus)** (Go backend
space module + web frontend in one repo, mirroring the `listus` layout).

That repo implements the approved **Assetus MVP** — the ownership *system of record*:
a Space-owned `Asset` with closed-enum `Category` / `Condition` / ownership-lifecycle
`Status` / `Visibility`, owner-type derivation from the owning Space type, an
append-only per-asset history, ownership transfer (Space → Space), and soft-archive /
hard-delete. Persisted under `/spaces/{spaceID}/ext/assetus/...` on Firestore via dalgo.

Spec (source of truth):
[`sneat-co/backstage` → `spec/features/assetus-mvp`](https://github.com/sneat-co/backstage/tree/main/spec/features/assetus-mvp).

## Why this code is still here

The MVP in the new repo is a **clean, narrower model** than this legacy package. The
legacy capabilities in this directory — vehicles, dwellings/real-estate, sport gear,
documents, mileage, liabilities, service providers, possession types — are **deliberately
deferred** by the MVP Feature's *Not Doing* section and are **not yet ported** to the new
repo. To honour "no functionality lost", this legacy code is **left in place** for now.

## Follow-up before this directory is removed

1. Port the still-wanted legacy capabilities into `sneat-co/assetus` (or explicitly
   retire them).
2. Resolve the cross-module coupling: `pkg/extensions/brandus/dbo4brands/make_test.go`
   imports `const4assetus.AssetCategory`, and `pkg/extensions/standard_extensions.go`
   registers `assetus.Extension()` — both must be updated when this package is removed.
3. Delete `pkg/extensions/assetus/` and verify `go build ./... && go test ./...`.

Until then, **do not build new Assetus features in this package** — build them in
`sneat-co/assetus`.
