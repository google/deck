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

// Package eventlog is a deck backend for writing log messages to Windows Event Log.
package eventlog

import (
	"github.com/google/deck"
)

// EventID is an attribute that appends Event Log Event IDs to log messages.
func EventID(id uint32) func(*deck.AttribStore) {
	return func(a *deck.AttribStore) {
		a.Store("EventID", id)
	}
}
