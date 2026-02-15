package api4logist

import (
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var defaultJsonWithAuthRequired = verify.Request(
	verify.AuthenticationRequired(true),
	verify.MaximumContentLength(verify.DefaultMaxJSONRequestSize),
)
