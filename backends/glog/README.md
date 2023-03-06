# The glog Backend for Deck

The glog backend is based on Google's github.com/golang/glog logging package.

The logger backend supports all platforms.

## Init

The glog backend takes a single setup parameter, `opts`.

The `opts` parameter allows the caller to supply a `glog.Options` struct to
customize the backend behavior.

## Attributes

### GlogV

The underlying glog package implements its own verbosity levels which are
independent of deck and controlled via separate, self-defined flags. The `GlogV`
attribute allows the caller to pass V() directly to the glog package without
invoking Deck's own verbosity handling (which affects *all* backends).

### deck.Depth

The glog package will attempt to make use of deck's core `Depth` attribute. This
allows the caller to modify the call depth, which may affect the file names and
line numbers rendered on the output. If Depth isn't specified, glog tries to do
the right thing and set depth to the original call site.

## Usage

```
import (
  github.com/google/deck
  github.com/google/deck/backends/glog
)

...
func main() {
  deck.Add(glog.Init(nil))
}
```
