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
	"github.com/google/deck"
)

// Init initializes the replay backend for use in a deck.
func Init() *Replay {
	return &Replay{}
}

// Replay is a log deck backend that records log messages, allowing them to be replayed later.
type Replay struct{}

// Close closes the replay backend.
func (d *Replay) Close() error { return nil }

type message struct {
	level deck.Level
}

// New creates a new replay message.
func (d *Replay) New(lvl deck.Level, msg string) deck.Composer {
	return &message{}
}

// Write records a new message to the replay backend.
func (m *message) Write() error {
	return nil
}

// Compose satisfies the composer interface. In the replay backend it has no other purpose.
func (m *message) Compose(s *deck.AttribStore) error {
	return nil
}
