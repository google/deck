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
	"strings"

	"github.com/google/deck"
	"golang.org/x/sys/windows/svc/eventlog"
)

// EventLog is a log deck backend that passes logs through to Windows Event Log.
type EventLog struct {
	handle *eventlog.Log
}

// Init initializes the EventLog backend for use in a deck.
//
// This version of Init does not automatically register the log source with Event Log. Log sources
// only require registration once, and can be registered outside of the application code, such as
// during software installation.
//
// If you want the backend to attempt registration, and your code is running with Administrator-level
// permissions, use InitWithInstall or InitWithDefaultInstall. Keep in mind that these will make
// (unnecessary) registration attempts each time the Init happens.
func Init(source string) (*EventLog, error) {
	hndl, err := eventlog.Open(source)
	if err != nil {
		return nil, err
	}
	return &EventLog{
		handle: hndl,
	}, nil
}

var (
	allSources uint32 = eventlog.Info | eventlog.Warning | eventlog.Error
)

// InitWithDefaultInstall initializes the EventLog backend for use in a deck while also installing
// the log source in the Windows registry.
//
// This method uses EventCreate.exe as the message file. EventCreate.exe is commonly available on
// most versions of Windows but it does *not* export all possible event IDs! Events with high
// ID numbers will render incorrectly in the event viewer using this message file.
//
// Registration of the source requires Administrator-level permissions.
func InitWithDefaultInstall(source string) (*EventLog, error) {
	if err := eventlog.InstallAsEventCreate(source, allSources); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	}
	return Init(source)
}

// InitWithInstall initializes the EventLog backend for use in a deck while also installing the log
// source in the Windows registry.
//
// This method requires the user to designate the message file. The message file must export all of
// the Event IDs used by the application in order for them to render correctly in Event Viewer.
//
// Registration of the source requires Administrator-level permissions.
func InitWithInstall(source string, messageFile string) (*EventLog, error) {
	if err := eventlog.Install(source, messageFile, true, allSources); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	}
	return Init(source)
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
