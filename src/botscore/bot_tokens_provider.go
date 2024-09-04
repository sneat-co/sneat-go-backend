package botscore

import (
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfwconst"
	"os"
	"strings"
)

var ErrBotTokenNotFound = errors.New("bot token not found by bot ID")

func GetBotToken(botPlatformID botsfwconst.Platform, botID string) (string, error) {
	envVarName := strings.ToUpper(string(botPlatformID)) + "_BOT_TOKEN_" + strings.ToUpper(botID)
	token := os.Getenv(envVarName)
	if token == "" {
		return "", fmt.Errorf("%w: botPlatform=%s, botID=%s, envVarName=%s", ErrBotTokenNotFound, botPlatformID, botID, envVarName)
	}
	return token, nil
}
