// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows
// +build windows

package eventlog

import (
	"errors"

	"github.com/google/deck"
	"golang.org/x/sys/windows/svc/eventlog"
)

// EventLog is a log deck backend that passes logs through to Windows Event Log.
type EventLog struct {
	handle *eventlog.Log
}

// Init initializes the EventLog backend for use in a deck.
func Init(source string) (*EventLog, error) {
	hndl, err := eventlog.Open(source)
	if err != nil {
		return nil, err
	}
	return &EventLog{
		handle: hndl,
	}, nil
}

// Close closes the EventLog backend.
func (e *EventLog) Close() error {
	return e.handle.Close()
}

type message struct {
	parent  *EventLog
	level   deck.Level
	message string
	eventID uint32
}

// New creates a new EventLog message.
func (e *EventLog) New(lvl deck.Level, msg string) deck.Composer {
	return &message{parent: e, level: lvl, message: msg, eventID: 1}
}

// Compose composes the Event prior to writing.
func (m *message) Compose(s *deck.AttribStore) error {
	id, ok := s.Load("EventID")
	if !ok {
		return errors.New("invalid EventID")
	}
	m.eventID = id.(uint32)
	return nil
}

// Write flushes the a stored log to Event Log.
func (m *message) Write() error {
	switch m.level {
	case deck.INFO:
		m.parent.handle.Info(m.eventID, m.message)
	case deck.WARNING:
		m.parent.handle.Warning(m.eventID, m.message)
	case deck.ERROR:
		m.parent.handle.Error(m.eventID, m.message)
	case deck.FATAL:
		m.parent.handle.Error(m.eventID, m.message)
	default:
		m.parent.handle.Info(m.eventID, m.message)
	}
	return nil
}
