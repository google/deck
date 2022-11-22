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

// Package syslog provides a deck backend for syslog.
package syslog

import (
	"log/syslog"

	"github.com/google/deck"
)

// Priority values passed through from the underlying syslog package.
const (
	LOG_EMERG   = syslog.LOG_EMERG
	LOG_ALERT   = syslog.LOG_ALERT
	LOG_CRIT    = syslog.LOG_CRIT
	LOG_ERR     = syslog.LOG_ERR
	LOG_WARNING = syslog.LOG_WARNING
	LOG_NOTICE  = syslog.LOG_NOTICE
	LOG_INFO    = syslog.LOG_INFO
	LOG_DEBUG   = syslog.LOG_DEBUG

	LOG_KERN     = syslog.LOG_KERN
	LOG_USER     = syslog.LOG_USER
	LOG_MAIL     = syslog.LOG_MAIL
	LOG_DAEMON   = syslog.LOG_DAEMON
	LOG_AUTH     = syslog.LOG_AUTH
	LOG_SYSLOG   = syslog.LOG_SYSLOG
	LOG_LPR      = syslog.LOG_LPR
	LOG_NEWS     = syslog.LOG_NEWS
	LOG_UUCP     = syslog.LOG_UUCP
	LOG_CRON     = syslog.LOG_CRON
	LOG_AUTHPRIV = syslog.LOG_AUTHPRIV
	LOG_FTP      = syslog.LOG_FTP
	LOG_LOCAL0   = syslog.LOG_LOCAL0
	LOG_LOCAL1   = syslog.LOG_LOCAL1
	LOG_LOCAL2   = syslog.LOG_LOCAL2
	LOG_LOCAL3   = syslog.LOG_LOCAL3
	LOG_LOCAL4   = syslog.LOG_LOCAL4
	LOG_LOCAL5   = syslog.LOG_LOCAL5
	LOG_LOCAL6   = syslog.LOG_LOCAL6
	LOG_LOCAL7   = syslog.LOG_LOCAL7
)

// Init initializes the Syslog backend for use in a deck.
func Init(tag string, facility syslog.Priority) (*Syslog, error) {
	handle, err := syslog.New(facility, tag)
	if err != nil {
		return nil, err
	}
	return &Syslog{
		handle: handle,
	}, nil
}

// Syslog is a log deck backend that passes logs through to the syslog package.
type Syslog struct {
	handle *syslog.Writer
}

// Close closes the syslog backend.
func (s *Syslog) Close() error {
	return s.handle.Close()
}

type message struct {
	parent  *Syslog
	level   deck.Level
	message string
}

// New creates a new Syslog message.
func (s *Syslog) New(lvl deck.Level, msg string) deck.Composer {
	return &message{parent: s, level: lvl, message: msg}
}

// Write flushes the a stored log to the IO writers.
func (m *message) Write() error {
	switch m.level {
	case deck.DEBUG:
		m.parent.handle.Debug(m.message)
	case deck.INFO:
		m.parent.handle.Info(m.message)
	case deck.WARNING:
		m.parent.handle.Warning(m.message)
	case deck.ERROR:
		m.parent.handle.Err(m.message)
	case deck.FATAL:
		m.parent.handle.Crit(m.message)
	default:
		m.parent.handle.Info(m.message)
	}
	return nil
}

// Compose composes the message prior to writing.
func (m *message) Compose(s *deck.AttribStore) error {
	return nil
}
