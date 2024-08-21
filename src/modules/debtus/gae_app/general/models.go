package general

type CreatedOn struct {
	CreatedOnPlatform string `firestore:",omitempty"` // e.g. "Telegram"
	CreatedOnID       string `firestore:",omitempty"` // e.g. "DebtsTrackerBot"
}
