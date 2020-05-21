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

package actor

import (
	"math/rand"
	"time"
)

// NewExponentialBackoffStrategy creates a new Supervisor strategy that restarts a faulting child using an exponential
// back off algorithm:
//
//	delay =
func NewExponentialBackoffStrategy(backoffWindow time.Duration, initialBackoff time.Duration) SupervisorStrategy {
	return &exponentialBackoffStrategy{
		backoffWindow:  backoffWindow,
		initialBackoff: initialBackoff,
	}
}

type exponentialBackoffStrategy struct {
	backoffWindow  time.Duration
	initialBackoff time.Duration
}

func (strategy *exponentialBackoffStrategy) HandleFailure(supervisor Supervisor, child *PID, rs *RestartStatistics, reason interface{}, message interface{}) {
	strategy.setFailureCount(rs)

	backoff := rs.FailureCount() * int(strategy.initialBackoff.Nanoseconds())
	noise := rand.Intn(500)
	dur := time.Duration(backoff + noise)
	time.AfterFunc(dur, func() {
		supervisor.RestartChildren(child)
	})
}

func (strategy *exponentialBackoffStrategy) setFailureCount(rs *RestartStatistics) {
	if rs.NumberOfFailures(strategy.backoffWindow) == 0 {
		rs.Reset()
	}

	rs.Fail()
}
