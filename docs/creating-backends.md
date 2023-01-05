# Creating Deck Backends

Deck backends are intended to be plug-and-play. Packages that implement deck
logging can add any combination of backends, with log lines being sent to each
attached backend at commit.

Backends have substantial internal flexibility with how they behave. The only
real requirement for a backend is that it implement the necessary function calls
to integrate with the core deck package. Beyond that, what a given backend
decides to do with incoming log messages can vary based on the whims of the
package's author.

## Architecture

Each backend should implement at least two structs, one for the backend itself
and one for the backend's message type.

This doc uses examples from the eventlog backend.

### Backend

A backend object is typically instantiated only once, generally around the time
the backend is added to the deck. The backend object can hold open connections,
handles, or other state that are shared across the individual messages.

```
type EventLog struct {
    handle *eventlog.Log
}
```

#### Init

Each backend must declare an Init() function, which is called prior to
`deck.Add()`. Note that the Init doesn't have to conform to a specific
interface, meaning it's reasonable for the Init to collect input from the user
for the sake of configuring the backend.

The Init *should* return a copy of the backend struct. It *may* also opt to
return an error, if the backend setup can fail.

```
func Init(source string) (*EventLog, error) {
    hndl, err := eventlog.Open(source)
    if err != nil {
        return nil, err
    }
    return &EventLog{
        handle: hndl,
    }, nil
}
```

#### Close()

Each backend must declare a Close() member to perform clean up on any open
handles or other internal resources.

```
func (e *EventLog) Close() error {
    return e.handle.Close()
}
```

#### New()

The backend must define a New() function which generates new [message](#Message)
objects. New must accept a Level and string (the message), and return a
deck.Composer corresponding to the generated Message.

```
func (e *EventLog) New(lvl deck.Level, msg string) deck.Composer {
    return &message{parent: e, level: lvl, message: msg, eventID: 1}
}
```

### Message

The message struct defines the structure for each individual log message. Every
message (deck.Info, deck.Error, etc) will correspond to a unique instance of
this object. The message includes at least the log string, but may include other
backend-specific metadata.

```
type message struct {
    parent  *EventLog
    level   deck.Level
    message string
    eventID uint32
}
```

#### Compose()

The message must support the Compose() method. Compose is required for the
`A`ttribute extended logging functions: Any time the With() function is used on
a message, the message's internal AttribStore is populated with metadata. When
Compose is called, the message has an opportunity to interrogate the AttribStore
for any relevant key/value sets. This allows backends to modify their behavior
based on the presence of known attributes in the store.

Backends aren't require to perform any specific behavior inside of Compose(),
and can safely return nil if they don't require any extended attributes.

A message's AttribStore is shared between all registered backends. Compose()
should take care to avoid mutating any stored content.

```
func (m *message) Compose(s *deck.AttribStore) error {
    id, ok := s.Load("EventID")
    if !ok {
        return errors.New("invalid EventID")
    }
    m.eventID = id.(uint32)
    return nil
```

#### Write()

Messages must provide the Write() method. Write() signals the message to flush
its content to the final destination.

```
func (m *message) Write() error {
    switch m.level {
    case deck.INFO:
        m.parent.handle.Info(m.eventID, m.message)
    case deck.WARNING:
        m.parent.handle.Warning(m.eventID, m.message)
    ...
    default:
        m.parent.handle.Info(m.eventID, m.message)
    }
    return nil
}
```
