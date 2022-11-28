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

package deck_test

import (
	"os"

	"github.com/google/deck/backends/logger"
	"github.com/google/deck"
)

func init() {
	deck.Add(logger.Init(os.Stdout, 0))
}

func ExampleInfo() {
	deck.Info("This is a simple log line.")
}

func ExampleInfof() {
	deck.Infof("Is this a %s line? %t", "format", true)
}
