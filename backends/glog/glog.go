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

// Package glog provides a deck backend for github.com/golang/glog.
//
// The glog backend supports two types of logging to the glog package. The
// first is to use deck's standard level functions .Info(), .Error(), etc. These map
// to the same function calls inside glog.
//
// We also support glog's V-style logging by exporting a V() attribute. This attribute
// is separate from deck's core V() attribute and only affects glog. When a glog.V()
// attribute is applied, the V level is passed to glog's V-style info functions.
//
// Example:
//
//	deck.InfoA("a message at level 3").With(glog.V(3)).Go()
//
// Providing the V() attribute overrides the deck level, so this:
//
//	deck.ErrorA("a message with verbosity").With(glog.V(3)).Go()
//
// is the same as this:
//
//	deck.InfoA("a message with verbosity").With(glog.V(3)).Go()
//
// The glog backend will also attempt to honor deck's Depth() attribute, allowing the
// caller to modify the stack frames output by glog messages.
package glog

import (
	"errors"

	log "github.com/golang/glog"
	"github.com/google/deck"
)

// Options allows the caller to customize the behavior of the glog backend.
type Options struct {
	// DebugLevel indicates the V() to use for deck Debug messages. Defaults to 1.
	DebugLevel log.Level
}

// Init initializes the glog backend for use in a deck.
func Init(opts *Options) *GLog {
	if opts == nil {
		opts = &Options{
			DebugLevel: 1,
		}
	}
	return &GLog{opts: opts}
}

// GLog is a log deck backend that passes logs through to the glog package.
type GLog struct {
	opts *Options
}

// Close closes the glog backend.
func (g *GLog) Close() error {
	return nil
}

type message struct {
	parent    *GLog
	level     deck.Level
	glogLevel log.Level
	message   string
	depth     int
}

// New creates a new GLog message.
func (g *GLog) New(lvl deck.Level, msg string) deck.Composer {
	return &message{parent: g, level: lvl, message: msg}
}

// Adding an offset of 3 excludes the frames in glog.go and deck.go, so the user's
// code locations should be rendered by default.
const depthOffset = 3

// Write flushes the stored message to glog.
func (m *message) Write() error {
	if m.glogLevel != 0 {
		log.V(m.glogLevel).InfoDepth(m.depth+depthOffset-1, m.message)
		return nil
	}
	switch m.level {
	case deck.DEBUG:
		log.V(m.parent.opts.DebugLevel).InfoDepth(m.depth+depthOffset, m.message)
	case deck.INFO:
		log.InfoDepth(m.depth+depthOffset, m.message)
	case deck.WARNING:
		log.WarningDepth(m.depth+depthOffset, m.message)
	case deck.ERROR:
		log.ErrorDepth(m.depth+depthOffset, m.message)
	case deck.FATAL:
		log.FatalDepth(m.depth+depthOffset, m.message)
	default:
		log.InfoDepth(m.depth+depthOffset, m.message)
	}
	return nil
}

// Compose composes the message prior to writing.
func (m *message) Compose(s *deck.AttribStore) error {
	id, ok := s.Load("GlogV")
	if !ok {
		return errors.New("invalid GlogV")
	}
	m.glogLevel = id.(log.Level)

	dep, ok := s.Load("Depth")
	if !ok {
		return errors.New("invalid Depth")
	}
	m.depth = dep.(int)
	return nil
}

// V is an attribute that allows use of glog's V-style logging.
//
// This is different from deck's V() attribute, which affects the behavior of all backends!
// glog.V() only modifies the behavior of glog.
func V(level log.Level) func(*deck.AttribStore) {
	return func(a *deck.AttribStore) {
		a.Store("GlogV", level)
	}
}
