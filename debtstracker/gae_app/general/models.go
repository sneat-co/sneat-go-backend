package general

type CreatedOn struct {
	CreatedOnPlatform string `datastore:",noindex,omitempty"` // e.g. "Telegram"
	CreatedOnID       string `datastore:",noindex,omitempty"` // e.g. "DebtsTrackerBot"
}
