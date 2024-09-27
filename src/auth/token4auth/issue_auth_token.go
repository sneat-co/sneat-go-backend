package token4auth

import "context"

var IssueAuthToken func(ctx context.Context, userID string, issuer string) (token string, err error)
