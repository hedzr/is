package basics

type (
	// Peripheral represents a basic infrastructure which can
	// be initialized and destroyed.
	//
	// For a Peripheral, the host should add it into a list
	// and destroy them while host is shutting down.
	Peripheral interface {
		// Close provides a closer to cleanup the peripheral gracefully
		Close()
	}

	Closable interface {
		// Close provides a closer to cleanup the peripheral gracefully
		Close()
	}

	Closer interface { // = io.Closer
		Close() error // = io.Closer
	}

	// AutoStart identify a peripheral object can be started automatically.
	// see AddPeripheral.
	AutoStart interface {
		AutoStart()
	}
)

// Basic is a base type to simplify your codes since you're using Peripheral type.
type Basic struct {
	peripherals []Peripheral
}

// AddPeripheral adds a Peripheral object into Basic holder/host.
//
// A peripheral represents an external resource such as redis manager which manages the links to remote redis server, etc..
// A peripheral can be auto-started implicit by AddPeripheral while it implements AutoStart interface.
func (s *Basic) AddPeripheral(peripherals ...Peripheral) {
	s.peripherals = append(s.peripherals, peripherals...)
	for _, p := range peripherals {
		if as, ok := p.(AutoStart); ok {
			as.AutoStart()
		}
	}
}

// Close provides a closer to cleanup the peripheral gracefully
func (s *Basic) Close() {
	for _, p := range s.peripherals {
		if p != nil {
			p.Close()
		}
	}
	s.peripherals = nil
}
