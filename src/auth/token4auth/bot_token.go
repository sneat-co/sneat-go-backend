package token4auth

import "fmt"

func IssueBotToken(userID, botPlatformID, botID string) string {
	issuer := GetBotIssuer(botPlatformID, botID)
	return IssueToken(userID, issuer)
}

func GetBotIssuer(botPlatformID, botID string) string {
	return fmt.Sprintf("%s:%s", botPlatformID, botID)
}
