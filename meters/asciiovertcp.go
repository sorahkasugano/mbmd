package meters

import (
	"time"

	"github.com/grid-x/modbus"
)

// ASCIIOverTCP is an ASCII encoder over a TCP modbus connection
type ASCIIOverTCP struct {
	Client  modbus.Client
	Handler *modbus.ASCIIOverTCPClientHandler
	prevID  uint8
}

// NewASCIIOverTCPClientHandler creates a TCP modbus handler
func NewASCIIOverTCPClientHandler(device string) *modbus.ASCIIOverTCPClientHandler {
	handler := modbus.NewASCIIOverTCPClientHandler(device)

	// set default timings
	handler.ProtocolRecoveryTimeout = 10 * time.Second // not used
	handler.LinkRecoveryTimeout = 15 * time.Second     // not used

	return handler
}

var _ Connection = (*ASCIIOverTCP)(nil)

// NewASCIIOverTCP creates a TCP modbus client
func NewASCIIOverTCP(address string) *ASCIIOverTCP {
	handler := NewASCIIOverTCPClientHandler(address)
	client := modbus.NewClient(handler)

	b := &ASCIIOverTCP{
		Client:  client,
		Handler: handler,
	}

	return b
}

// String returns the bus connection address (TCP)
func (b *ASCIIOverTCP) String() string {
	return b.Handler.Address
}

// ModbusClient returns the TCP modbus client
func (b *ASCIIOverTCP) ModbusClient() modbus.Client {
	return b.Client
}

// Logger sets a logging instance for physical bus operations
func (b *ASCIIOverTCP) Logger(l Logger) {
	b.Handler.Logger = l
}

// Slave sets the modbus device id for the following operations
func (b *ASCIIOverTCP) Slave(deviceID uint8) {
	// Some devices like SDM need to have a little pause between querying different device ids
	if b.prevID != 0 && deviceID != b.prevID {
		time.Sleep(time.Duration(100) * time.Millisecond)
		b.prevID = deviceID
	}

	b.Handler.SetSlave(deviceID)
}

// Timeout sets the modbus timeout
func (b *ASCIIOverTCP) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

// ConnectDelay sets the the initial delay after connecting before starting communication
func (b *ASCIIOverTCP) ConnectDelay(delay time.Duration) {
	b.Handler.ConnectDelay = delay
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *ASCIIOverTCP) Close() {
	b.Handler.Close()
}

// Clone clones the modbus connection.
func (b *ASCIIOverTCP) Clone(deviceID byte) Connection {
	handler := b.Handler.Clone()
	handler.SetSlave(deviceID)

	return &ASCIIOverTCP{
		Client:  modbus.NewClient(handler),
		Handler: handler,
	}
}
