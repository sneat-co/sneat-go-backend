package const4debtus

import "errors"

var (
	ErrJsonCountMismatch = errors.New("json slice length is different to length of corresponding count property")
)
