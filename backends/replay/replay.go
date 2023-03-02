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
		recorder: Buffer{},
	}
}

// Buffer aggregates message entries as they're written to the deck. Each Replay instance keeps
// an internal Buffer, and returns copies of the Buffer to the user when queried.
type Buffer []Entry

// Len returns the length of the buffer.
func (b Buffer) Len() int { return len(b) }

// Entry models a log entry as it's written to the buffer. It tracks the log message but also other
// metadata that we may want to recall later, like the deck Level.
type Entry struct {
	Level   deck.Level
	Message string
}

// Replay is a log deck backend that records log messages, allowing them to be replayed later.
type Replay struct {
	mu       sync.Mutex
	recorder Buffer
}

func (r *Replay) append(entry Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recorder = append(r.recorder, entry)
}

func (r *Replay) byLevel(lvl deck.Level) Buffer {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := Buffer{}
	for _, a := range r.recorder {
		if a.Level == lvl {
			out = append(out, a)
		}
	}
	return out
}

// Debug returns all messages recorded to the debug level.
func (r *Replay) Debug() Buffer {
	return r.byLevel(deck.DEBUG)
}

// Error returns all messages recorded to the error level.
func (r *Replay) Error() Buffer {
	return r.byLevel(deck.ERROR)
}

// Fatal returns all messages recorded to the fatal level.
func (r *Replay) Fatal() Buffer {
	return r.byLevel(deck.FATAL)
}

// Info returns all messages recorded to the info level.
func (r *Replay) Info() Buffer {
	return r.byLevel(deck.INFO)
}

// Warning returns all messages recorded to the warning level.
func (r *Replay) Warning() Buffer {
	return r.byLevel(deck.WARNING)
}

// Reset resets the replay deck to its initial state.
func (r *Replay) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recorder = Buffer{}
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
		m.parent.append(Entry{deck.DEBUG, m.message})
	case deck.INFO:
		m.parent.append(Entry{deck.INFO, m.message})
	case deck.WARNING:
		m.parent.append(Entry{deck.WARNING, m.message})
	case deck.ERROR:
		m.parent.append(Entry{deck.ERROR, m.message})
	case deck.FATAL:
		m.parent.append(Entry{deck.FATAL, m.message})
	default:
		m.parent.append(Entry{DEFAULT, m.message})
	}
	return nil
}

// Compose satisfies the composer interface. In the replay backend it has no other purpose.
func (m *message) Compose(s *deck.AttribStore) error {
	return nil
}
