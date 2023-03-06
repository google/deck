# The Replay Backend for Deck

The replay backend allows log messages to be recorded and replayed.

The replay backend supports all platforms.

## Init

The replay backend does not take any setup parameters.

## Attributes

The replay backend does not utilize any custom attributes.

## Details & Features

The replay backend is particularly useful as part of a testing framework, when
the code under test is already instrumented with deck.

### Log Storage & Retrieval

Once attached to a deck, the replay backend will store any newly logged messages
in memory. Messages are indexed by their deck level and can be retrieved from
the backend using the corresponding level functions (`Info()`, `Error()`, etc.).
The `All()` function returns all collected messages in the order they were
recorded.

Each retrieval function provides the results in a Bundle. The Bundle is an
ordered list of log entries as they were processed by the deck. The Bundle also
provides some helper functions to simplify searching the message contents.

Replay always returns a *copy* of the original Bundle each time a retrieval
function is called. The underlying Bundle is kept private to avoid potential
conflicts between new log events and user activity.

### Bundle Helpers

`Bundle.ContainsRE` allows the user to use a regular expression to search for
matching messages in the Bundle. The Bundle iterates over all messages in the
list and attempts to match `re`, returning true if found.

`Bundle.ContainsString` allows the user to search the Bundle for a string. This
leverages `strings.Contains`, so substrings are matched as well. The function
returns true if a match exists in any of the messages.

## Usage

```
import (
  github.com/google/deck
  github.com/google/deck/backends/replay
)

...
func TestSomething(t *testing.T) {
  d := deck.New()
  r := replay.Init()
  d.Add(r)
  ... execute code that creates logs ...
  if !r.Info().Contains("an expected log event") {
    t.Errorf("expected log event not found")
  }
}
```
