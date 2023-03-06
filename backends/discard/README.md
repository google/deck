# The Discard Backend for Deck

The discard backend allows log messages to be discarded.

If logs are processed with no backends attached, deck will throw warnings.
Attaching the discard backend satisfies the need to have at least one backend
registered without actually sending log messages anywhere.

The discard backend supports all platforms.

## Init

The discard backend does not take any setup parameters.

## Attributes

The discard backend does not utilize any custom attributes.

## Usage

```
import (
  github.com/google/deck
  github.com/google/deck/backends/discard
)

...
func main() {
  // silence log events
  deck.Add(discard.Init())
}
```
