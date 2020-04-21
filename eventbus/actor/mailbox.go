package actor

import "github.com/dad-go/eventbus/mailbox"

var (
	defaultDispatcher = mailbox.NewDefaultDispatcher(300)
)

var defaultMailboxProducer = mailbox.Unbounded()
