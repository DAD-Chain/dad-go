/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package mailbox

import (
	"github.com/dad-go/eventbus/internal/queue/mpsc"
)

// UnboundedLockfree returns a producer which creates an unbounded, lock-free mailbox.
// This mailbox is cheaper to allocate, but has a slower throughput than the plain Unbounded mailbox.
func UnboundedLockfree(mailboxStats ...Statistics) Producer {
	return func(invoker MessageInvoker, dispatcher Dispatcher) Inbound {
		return &defaultMailbox{
			userMailbox:   mpsc.New(),
			systemMailbox: mpsc.New(),
			invoker:       invoker,
			mailboxStats:  mailboxStats,
			dispatcher:    dispatcher,
		}
	}
}
