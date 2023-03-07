// Copyright 2023 Google LLC
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

// Package replay provides a deck backend for replaying log messages.
package replay

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/google/deck"
)

var (
	// DEFAULT is a deck level for any messages not caught as coming from one of the standard levels.
	DEFAULT deck.Level = 1000
)

// Init initializes the replay backend for use in a deck.
func Init() *Replay {
	return &Replay{
		recorder: Bundle{},
	}
}

// Bundle aggregates message entries as they're written to the deck. Each Replay instance keeps
// an internal Bundle, and returns copies of the Bundle to the user when queried.
type Bundle []Log

// ContainsString searches the Bundle for a string which contains str. This call uses strings.Contains
// which will match on substrings in addition to full strings.
func (b Bundle) ContainsString(str string) bool {
	for _, a := range b {
		if strings.Contains(a.Message, str) {
			return true
		}
	}
	return false
}

// ContainsRE searches the Bundle for a message which matches the regular expression re.
func (b Bundle) ContainsRE(re *regexp.Regexp) bool {
	for _, a := range b {
		if re.MatchString(a.Message) {
			return true
		}
	}
	return false
}

// Len returns the length of the Bundle.
func (b Bundle) Len() int { return len(b) }

// Log models a log entry as it's written to the Bundle. It tracks the log message but also other
// metadata that we may want to recall later, like the deck Level.
type Log struct {
	Level   deck.Level
	Message string
}

// String stringifies Log objects for nicer printing.
func (e Log) String() string {
	levels := map[deck.Level]string{
		deck.DEBUG:   "DEBUG",
		deck.ERROR:   "ERROR",
		deck.WARNING: "WARNING",
		deck.INFO:    "INFO",
		deck.FATAL:   "FATAL",
		DEFAULT:      "DEFAULT",
	}
	return fmt.Sprintf("%s: %q", levels[e.Level], e.Message)
}

// Replay is a log deck backend that records log messages, allowing them to be replayed later.
type Replay struct {
	mu       sync.Mutex
	recorder Bundle
}

func (r *Replay) append(entry Log) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recorder = append(r.recorder, entry)
}

func (r *Replay) byLevel(lvl deck.Level) Bundle {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := Bundle{}
	for _, a := range r.recorder {
		if a.Level == lvl {
			out = append(out, a)
		}
	}
	return out
}

// All returns all messages recorded to all levels.
func (r *Replay) All() Bundle {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(Bundle, len(r.recorder))
	copy(out, r.recorder)
	return out
}

// Debug returns all messages recorded to the debug level.
func (r *Replay) Debug() Bundle {
	return r.byLevel(deck.DEBUG)
}

// Error returns all messages recorded to the error level.
func (r *Replay) Error() Bundle {
	return r.byLevel(deck.ERROR)
}

// Fatal returns all messages recorded to the fatal level.
func (r *Replay) Fatal() Bundle {
	return r.byLevel(deck.FATAL)
}

// Info returns all messages recorded to the info level.
func (r *Replay) Info() Bundle {
	return r.byLevel(deck.INFO)
}

// Warning returns all messages recorded to the warning level.
func (r *Replay) Warning() Bundle {
	return r.byLevel(deck.WARNING)
}

// Reset resets the replay deck to its initial state.
func (r *Replay) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recorder = Bundle{}
}

// Close closes the replay backend.
func (r *Replay) Close() error { return nil }

type message struct {
	message string
	parent  *Replay
	level   deck.Level
}

// New creates a new replay message.
func (r *Replay) New(lvl deck.Level, msg string) deck.Composer {
	return &message{message: msg, parent: r, level: lvl}
}

// Write records a new message to the replay backend.
func (m *message) Write() error {
	switch m.level {
	case deck.DEBUG:
		m.parent.append(Log{deck.DEBUG, m.message})
	case deck.INFO:
		m.parent.append(Log{deck.INFO, m.message})
	case deck.WARNING:
		m.parent.append(Log{deck.WARNING, m.message})
	case deck.ERROR:
		m.parent.append(Log{deck.ERROR, m.message})
	case deck.FATAL:
		m.parent.append(Log{deck.FATAL, m.message})
	default:
		m.parent.append(Log{DEFAULT, m.message})
	}
	return nil
}

// Compose satisfies the composer interface. In the replay backend it has no other purpose.
func (m *message) Compose(s *deck.AttribStore) error {
	return nil
}
