# Calendarius

The calendar module (happenings, slots, scheduling).

The module's persisted extension id is `calendarius` (see `const4calendarius/module_id.go`),
so its data lives under `/spaces/{spaceID}/ext/calendarius/...`.

Naming history: this module previously used split names — `calendarium` on the Go backend and
`schedulus` on the Angular frontend. It was unified to the single name **`calendarius`** across the
whole Sneat ecosystem.
