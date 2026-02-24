package domain

type AggregateRoot interface {
	ID() string
	Version() int
	PullEvents() []Event
}

type BaseAggregate struct {
	IDValue      string
	VersionValue int
	events       []Event
}

func (a *BaseAggregate) ID() string          { return a.IDValue }
func (a *BaseAggregate) Version() int        { return a.VersionValue }
func (a *BaseAggregate) Raise(e Event)       { a.events = append(a.events, e) }
func (a *BaseAggregate) PullEvents() []Event { out := a.events; a.events = nil; return out }
