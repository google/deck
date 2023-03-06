# The Logger Backend for Deck

The logger backend allows log messages to written to any io.Writer.

The logger backend supports all platforms.

## Init

The logger backend takes two setup parameters, `out` and `flags`.

The `out` parameter must be an io.Writer. Log messages will be flushed directly
to the io.Writer specified. An io.MultiWriter may be used here as well. Note
that the logger backend does not close this Writer even if the user calls
logger.Close(); it is up to the user to manage the io.Writer handle.

The `flags` parameter references one or more
[flag constants](https://pkg.go.dev/log#pkg-constants) defined in Go's core
`log` package. These flags can be used to modify the rendering of the log lines
as they're written to the output. The default behavior is to use log.LstdFlags.

## Attributes

### deck.Depth

The logger package will attempt to make use of deck's core `Depth` attribute.
This allows the caller to modify the call depth, which may affect the file names
and line numbers rendered on the output, depending on which log flags are
supplied during setup. If Depth isn't specified, logger tries to do the right
thing and set depth to the original call site.

## Usage

```
import (
  log
  github.com/google/deck
  github.com/google/deck/backends/logger
)

...

func main() {
  // Log to Stdout
  deck.Add(logger.Init(os.Stdout, log.LstdFlags))

  // Log to a file
  lf, err := os.OpenFile("/tmp/my_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
  if err != nil {
      os.Exit(1)
  }
  defer lf.Close()
  deck.Add(logger.Init(lf, log.LstdFlags))
}
```
