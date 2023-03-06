# The Syslog Backend for Deck

The syslog backend is based on Go's core `log/syslog` package, allowing log
messages to be written to the system log service.

The syslog backend does not support Windows.

## Init

The syslog backend takes two setup parameters, `tag` and `facility`.

The `tag` parameter allows the caller to supply a fixed prefix to log messages.
This is passed directly to the underlying syslog package.

The `facility` parameter allows the caller to specify the syslog facility and
severity. This is passed directly to the underlying syslog package.

## Attributes

The syslog backend does not utilize any custom attributes.

## Usage

```
import (
  github.com/google/deck
  github.com/google/deck/backends/syslog
)

...
func main() {
  sl, err := syslog.Init("MY-APP", syslog.LOG_INFO|syslog.LOG_USER)
  if err != nil {
    os.Exit(1)
  }
  deck.Add(sl)
}
```
