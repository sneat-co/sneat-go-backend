package queues

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/const4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/const4splitus"
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
