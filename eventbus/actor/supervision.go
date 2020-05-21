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

import "github.com/dad-go/eventbus/eventstream"

// DeciderFunc is a function which is called by a SupervisorStrategy
type DeciderFunc func(reason interface{}) Directive

//SupervisorStrategy is an interface that decides how to handle failing child actors
type SupervisorStrategy interface {
	HandleFailure(supervisor Supervisor, child *PID, rs *RestartStatistics, reason interface{}, message interface{})
}

//Supervisor is an interface that is used by the SupervisorStrategy to manage child actor lifecycle
type Supervisor interface {
	Children() []*PID
	EscalateFailure(reason interface{}, message interface{})
	RestartChildren(pids ...*PID)
	StopChildren(pids ...*PID)
	ResumeChildren(pids ...*PID)
}

func logFailure(child *PID, reason interface{}, directive Directive) {
	eventstream.Publish(&SupervisorEvent{
		Child:     child,
		Reason:    reason,
		Directive: directive,
	})
}

//DefaultDecider is a decider that will always restart the failing child actor
func DefaultDecider(_ interface{}) Directive {
	return RestartDirective
}

var (
	defaultSupervisionStrategy    = NewOneForOneStrategy(10, 0, DefaultDecider)
	restartingSupervisionStrategy = NewRestartingStrategy()
)

func DefaultSupervisorStrategy() SupervisorStrategy {
	return defaultSupervisionStrategy
}

func RestartingSupervisorStrategy() SupervisorStrategy {
	return restartingSupervisionStrategy
}
