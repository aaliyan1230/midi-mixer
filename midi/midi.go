package midi

import (
	"fmt"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

// CCMessage represents a MIDI Control Change message
type CCMessage struct {
	Channel    uint8
	Controller uint8
	Value      uint8
}

// Common MIDI CC numbers for mixer controls
const (
	CCVolume     uint8 = 7
	CCPan        uint8 = 10
	CCExpression uint8 = 11
	CCReverb     uint8 = 91
	CCChorus     uint8 = 93
)

// Handler manages MIDI input/output connections
type Handler struct {
	inPort    drivers.In
	outPort   drivers.Out
	stopFunc  func()
	ccChan    chan CCMessage
	mu        sync.RWMutex
	connected bool
}

// NewHandler creates a new MIDI handler
func NewHandler() *Handler {
	return &Handler{
		ccChan: make(chan CCMessage, 100),
	}
}

// GetInputPorts returns available MIDI input ports
func GetInputPorts() []drivers.In {
	return midi.GetInPorts()
}

// GetOutputPorts returns available MIDI output ports
func GetOutputPorts() []drivers.Out {
	return midi.GetOutPorts()
}

// Connect opens the specified input and output ports
func (h *Handler) Connect(inPort drivers.In, outPort drivers.Out) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connected {
		h.disconnect()
	}

	h.inPort = inPort
	h.outPort = outPort

	// Open output port if specified
	if outPort != nil {
		if err := outPort.Open(); err != nil {
			return fmt.Errorf("failed to open output port: %w", err)
		}
	}

	// Start listening on input port if specified
	if inPort != nil {
		stop, err := midi.ListenTo(inPort, h.handleMIDI, midi.UseSysEx())
		if err != nil {
			if outPort != nil {
				outPort.Close()
			}
			return fmt.Errorf("failed to listen on input port: %w", err)
		}
		h.stopFunc = stop
	}

	h.connected = true
	return nil
}

// handleMIDI processes incoming MIDI messages
func (h *Handler) handleMIDI(msg midi.Message, timestampms int32) {
	var ch, cc, val uint8
	if msg.GetControlChange(&ch, &cc, &val) {
		select {
		case h.ccChan <- CCMessage{Channel: ch, Controller: cc, Value: val}:
		default:
			// Channel full, drop message
		}
	}
}

// CCChannel returns the channel for receiving CC messages
func (h *Handler) CCChannel() <-chan CCMessage {
	return h.ccChan
}

// SendCC sends a Control Change message
func (h *Handler) SendCC(channel, controller, value uint8) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.outPort == nil || !h.connected {
		return nil // No output port, silently ignore
	}

	msg := midi.ControlChange(channel, controller, value)
	return h.outPort.Send(msg)
}

// disconnect closes all ports (must be called with lock held)
func (h *Handler) disconnect() {
	if h.stopFunc != nil {
		h.stopFunc()
		h.stopFunc = nil
	}
	if h.outPort != nil {
		h.outPort.Close()
	}
	h.connected = false
}

// Close closes all MIDI connections
func (h *Handler) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.disconnect()
	close(h.ccChan)
}

// IsConnected returns whether MIDI is connected
func (h *Handler) IsConnected() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.connected
}

// GetInputPortName returns the name of the connected input port
func (h *Handler) GetInputPortName() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.inPort != nil {
		return h.inPort.String()
	}
	return "None"
}

// GetOutputPortName returns the name of the connected output port
func (h *Handler) GetOutputPortName() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.outPort != nil {
		return h.outPort.String()
	}
	return "None"
}
