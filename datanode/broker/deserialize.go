package broker

import (
	"context"
	"fmt"

	"github.com/zeta-protocol/zeta/core/events"
	"github.com/zeta-protocol/zeta/libs/proto"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
)

type Deserializer struct {
	source RawEventReceiver
}

func NewDeserializer(source RawEventReceiver) *Deserializer {
	return &Deserializer{
		source: source,
	}
}

func (e *Deserializer) Listen() error {
	return e.source.Listen()
}

func (e *Deserializer) Receive(ctx context.Context) (<-chan events.Event, <-chan error) {
	in, inErr := e.source.Receive(ctx)

	out := make(chan events.Event)
	outErr := make(chan error)

	go func() {
		defer close(out)
		defer close(outErr)

		for {
			select {
			case rawEvent, ok := <-in:
				if !ok {
					return
				}
				event, err := deserializeEvent(rawEvent)
				if err != nil {
					outErr <- err
					return
				}

				// Listen for context cancels, even if we're blocked sending events
				select {
				case out <- event:
				case <-ctx.Done():
					return
				}
			case err, ok := <-inErr:
				if !ok {
					return
				}
				// Listen for context cancels, even if we're blocked sending events
				select {
				case outErr <- err:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, outErr
}

func deserializeEvent(rawEvent []byte) (events.Event, error) {
	busEvent := &eventspb.BusEvent{}

	if err := proto.Unmarshal(rawEvent, busEvent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bus event: %w", err)
	}

	if busEvent.Version != eventspb.Version {
		return nil, fmt.Errorf("mismatched BusEvent version received: %d, want %d", busEvent.Version, eventspb.Version)
	}

	event := toEvent(context.Background(), busEvent)
	if event == nil {
		return nil, fmt.Errorf("Can not convert proto '%s' event to internal event", busEvent.GetType().String())
	}
	return event, nil
}
