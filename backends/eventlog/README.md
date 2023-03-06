# The EventLog Backend for Deck

The eventlog backend allows log messages to be sent to Windows Event Log.

The eventlog backend only supports the Windows platform.

## Init

The eventlog backend takes a `source` setup parameter. `source` corresponds to
the name of the Event Log message source. The source must be registered with
Event Log for events to render properly in Event Viewer.

## Attributes

### eventlog.EventID

The eventlog package exports the `EventID` attribute. This attribute allows
events to be stored with custom Event IDs. The Event ID will default to 1 if not
specified.

## Details & Features

### Source Registration

Event Log requires message sources to be registered in order for events to
render properly in Event Viewer. This is a one-time operation that creates
registry keys. Event Log must be supplied with a Message File which exports any
and all Event IDs the application will use.

It is recommended to perform event registration as part of application setup or
installation. Registration has to be done with administrator privileges, and
only needs to happen once per source.

The eventlog backend provides two helpers which may be used to perform the
source registration: `InitWithInstall` and `InitWithDefaultInstall`. These
helpers have a couple caveats:

1.  They will attempt to register the source every time they're called (every
    time Init happens), which is somewhat wasteful given that registration does
    not need to happen every time an application wants to log to Event Log.
2.  They *must* be called with administrator-level privileges, which means they
    are inappropriate to use in any applications running with user-level
    privileges. In these scenarios, use application installer code to perform
    the registration in advance.

`InitWithInstall` requires the user to provide their own
[message file](https://learn.microsoft.com/en-us/windows/win32/eventlog/message-files).
The message file must export all event IDs used by the application.

`InitWithDefaultInstall` leverages code in the underlying eventlog package to
register %SystemRoot%\System32\EventCreate.exe as the message file. This file is
provided with most modern Windows versions, and is a safe bet for many simple
situations where the user does not want to provide their own message file,
however there is one important limitation: EventCreate.exe *does not* support
the full range of possible Event IDs. To use Event IDs greater than ~1000, you
will need to provide a different message file.

## Usage

```
import (
  github.com/google/deck
  github.com/google/deck/backends/eventlog
)

...

func main() {
  evt, err := eventlog.Init("My Application")
  if err != nil {
      os.Exit(1)
  }
  deck.Add(evt)
  deck.InfoA("A Windows event with Event ID.").With(eventlog.EventID(10)).Go()
}
```
