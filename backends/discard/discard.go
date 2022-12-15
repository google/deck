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

// Package discard provides a deck backend that discards all logs.
package discard

import (
	"github.com/google/deck"
)

// Init initializes the discard backend for use in a deck.
func Init() *Discard {
	return &Discard{}
}

// Discard is a log deck backend that simply discards messages. It satisfies the need to have at
// least one backend registered with deck in situations where you want to avoid any actual output.
type Discard struct{}

// Close closes the discard backend.
func (d *Discard) Close() error { return nil }

type message struct{}

// New creates a new discard message.
func (d *Discard) New(lvl deck.Level, msg string) deck.Composer {
	return &message{}
}

// Write satisfies the composer interface. In the discard backend it has no other purpose.
func (m *message) Write() error {
	return nil
}

// Compose satisfies the composer interface. In the discard backend it has no other purpose.
func (m *message) Compose(s *deck.AttribStore) error {
	return nil
}
