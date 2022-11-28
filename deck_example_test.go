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

package deck_example_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/google/deck/backends/logger"
	"github.com/google/deck"
)

// The global instance of deck can be configured during init, main, or anywhere else, as long
// as it happens before log messages are used.
func init() {
	f, err := ioutil.TempFile("", "deck_example")
	if err != nil {
		panic(err)
	}

	// logger supports an io.Writer, so we can pass a single writer or a multiwriter
	mw := io.MultiWriter(os.Stdout, f)

	// Add the backends we want to the global deck instance.
	deck.Add(logger.Init(mw, 0))
}

func WriteToDeck(d *deck.Deck) {
	d.Infoln("hello from WriteToDeck")
}

func Example() {
	// The global deck can be used in any package by referencing the top level deck package functions.
	deck.Info("this is the start of the example")
	deck.Errorf("this is an example error: %v", errors.New("oh no"))

	// Custom decks can also be initialized with separate settings and passed around.
	anotherDeck := deck.New()
	defer anotherDeck.Close()
	anotherDeck.Add(logger.Init(os.Stderr, 0))
	WriteToDeck(anotherDeck)
	anotherDeck.Info("this is the end of the example")
}
