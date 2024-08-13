package queues

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/const4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/const4userus"
)

var KnownQueues = []string{

	// Common queues
	QueueChats,
	QueueSupport,
	QueueReminders,
	QueueEmails,
	QueueInvites,

	const4userus.QueueUsers,
	const4contactus.QueueContacts,

	// Debtus module
	const4debtus.QueueDebtus,
	const4debtus.QueueTransfers,
	const4debtus.QueueReceipts,

	// Splitus module
	const4splitus.QueueSplitus,
}
