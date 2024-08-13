package const4userus

type AuthProviderCode = string

const (
	TelegramAuthProvider          AuthProviderCode = "telegram"
	GoogleAuthProvider            AuthProviderCode = "google"
	FacebookAuthProvider          AuthProviderCode = "fb"
	FacebookMessengerAuthProvider AuthProviderCode = "fbm"
	EmailAuthProvider             AuthProviderCode = "email"
	ViberAuthProvider             AuthProviderCode = "viber"
	LineAuthProvider              AuthProviderCode = "line"
	WeChatAuthProvider            AuthProviderCode = "wechat"
)

func IsKnownUserAccountProvider(p string) bool {
	switch p {
	case TelegramAuthProvider:
	case GoogleAuthProvider:
	case FacebookAuthProvider:
	case FacebookMessengerAuthProvider:
	case EmailAuthProvider:
	case ViberAuthProvider:
	case LineAuthProvider:
	case WeChatAuthProvider:
	default:
		return false
	}
	return true
}
