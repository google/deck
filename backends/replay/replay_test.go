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

package replay

import (
	"testing"

	"github.com/google/deck"
	"github.com/google/go-cmp/cmp"
)

func TestLevels(t *testing.T) {
	d := deck.New()
	r := Init()
	d.Add(r)
	tests := []struct {
		desc   string
		inputs []string
		want   Buffer
		fIn    func(message ...any)
		fOut   func() Buffer
	}{
		{
			"error messages",
			[]string{"error message one", "another error"},
			Buffer{Entry{deck.ERROR, "error message one"}, Entry{deck.ERROR, "another error"}},
			d.Error,
			r.Error,
		},
		{
			"info messages",
			[]string{"info message one"},
			Buffer{Entry{deck.INFO, "info message one"}},
			d.Info,
			r.Info,
		},
		{
			"warning messages",
			[]string{"warning message one", "warning message two"},
			Buffer{Entry{deck.WARNING, "warning message one"}, Entry{deck.WARNING, "warning message two"}},
			d.Warning,
			r.Warning,
		},
		{
			"an empty set",
			[]string{},
			Buffer{},
			d.Warning,
			r.Warning,
		},
	}
	for _, tt := range tests {
		r.Reset()
		for _, m := range tt.inputs {
			tt.fIn(m)
		}
		out := tt.fOut()
		if out.Len() != len(tt.inputs) {
			t.Errorf("%s: produced unexpected size of results: got %d, want %d", tt.desc, len(out), len(tt.inputs))
		}
		if diff := cmp.Diff(out, tt.want); diff != "" {
			t.Errorf("%s: produced unexpected diff: %s", tt.desc, diff)
		}
	}
}

func TestReset(t *testing.T) {
	d := deck.New()
	r := Init()
	d.Add(r)
	d.Info("message")
	d.Warning("message")
	if r.Info().Len() < 1 || r.Warning().Len() < 1 {
		t.Errorf("failed to record logs as expected")
	}
	r.Reset()
	if r.Info().Len() != 0 || r.Warning().Len() != 0 {
		t.Errorf("Reset() failed to reset logs as expected")
	}
}
