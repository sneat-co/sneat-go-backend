package token4auth

import (
	"context"
	"fmt"
)

var IssueBotToken = func(ctx context.Context, userID, botPlatformID, botID string) (string, error) {
	issuer := GetBotIssuer(botPlatformID, botID)
	return IssueAuthToken(ctx, userID, issuer)
}

func GetBotIssuer(botPlatformID, botID string) string {
	return fmt.Sprintf("%s:%s", botPlatformID, botID)
}
