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

// Package logger provides a deck backend that leverages Go's core log package.
package logger

import (
	"errors"
	"io"
	"log"

	"github.com/google/deck"
)

// Init initializes the logger backend for use in a deck.
func Init(out io.Writer, flags int) *Logger {
	if flags == 0 {
		flags = log.LstdFlags
	}
	return &Logger{
		debug:   log.New(out, TagDebug, flags),
		info:    log.New(out, TagInfo, flags),
		warning: log.New(out, TagWarning, flags),
		error:   log.New(out, TagError, flags),
		fatal:   log.New(out, TagFatal, flags),
	}
}

var (
	// TagDebug is the tag added to messages logged at the DEBUG level.
	TagDebug = "DEBUG: "
	// TagInfo is the tag added to messages logged at the INFO level.
	TagInfo = "INFO: "
	// TagWarning is the tag added to messages logged at the WARN level.
	TagWarning = "WARN: "
	// TagError is the tag added to messages logged at the ERROR level.
	TagError = "ERROR: "
	// TagFatal is the tag added to messages logged at the FATAL level.
	TagFatal = "FATAL: "
)

// Logger is a log deck backend that passes logs through to Go's core log package.
type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	fatal   *log.Logger
}

// Close closes the Logger backend. The io.Writer passed to Init() is not closed and
// must be closed by the caller.
func (l *Logger) Close() error { return nil }

type message struct {
	level   deck.Level
	message string
	depth   int
	parent  *Logger
}

// New creates a new logger message.
func (l *Logger) New(lvl deck.Level, msg string) deck.Composer {
	return &message{level: lvl, message: msg, parent: l}
}

// Depth affects certain flags including Llongfile and Lshortfile. Adding 3 excludes the frames in
// logger and deck.go, so the user's code locations should be rendered by default.
const depthOffset = 3

// Write flushes a stored log message.
func (m *message) Write() error {
	switch m.level {
	case deck.DEBUG:
		m.parent.debug.Output(m.depth+depthOffset, m.message)
	case deck.INFO:
		m.parent.info.Output(m.depth+depthOffset, m.message)
	case deck.WARNING:
		m.parent.warning.Output(m.depth+depthOffset, m.message)
	case deck.ERROR:
		m.parent.error.Output(m.depth+depthOffset, m.message)
	case deck.FATAL:
		m.parent.fatal.Output(m.depth+depthOffset, m.message)
	default: // any levels that don't map go to info
		m.parent.info.Output(m.depth+depthOffset, m.message)
	}
	return nil
}

// Compose composes the message prior to writing.
func (m *message) Compose(s *deck.AttribStore) error {
	id, ok := s.Load("Depth")
	if !ok {
		return errors.New("invalid Depth")
	}
	m.depth = id.(int)
	return nil
}
