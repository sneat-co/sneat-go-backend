# Generic API

Handlers are registered in [endpoints.go](endpoints.go).

There is 3 endpoints:

- POST /api/$generic/[create](create.go)?entity={kind}
- PUT /api/$generic/[update](update.go)?entity={kind:id}
- DELETE /api/$generic/[delete](delete.go)?entity={kind:id}