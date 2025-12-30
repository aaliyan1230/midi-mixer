package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"midi-mixer/midi"
	"midi-mixer/mixer"
	"midi-mixer/ui"
)

// View represents the current screen
type View int

const (
	ViewMixer View = iota
	ViewDevices
)

// Model is the main application model
type Model struct {
	state          *mixer.State
	deviceSelector *ui.DeviceSelector
	currentView    View
	width          int
	height         int
	err            error
}

// MidiMsg is sent when a MIDI CC message is received
type MidiMsg midi.CCMessage

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		listenForMidi(m.state.MidiHandler),
	)
}

// listenForMidi creates a command that listens for MIDI messages
func listenForMidi(handler *midi.Handler) tea.Cmd {
	return func() tea.Msg {
		msg := <-handler.CCChannel()
		return MidiMsg(msg)
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case MidiMsg:
		m.handleMidiCC(midi.CCMessage(msg))
		return m, listenForMidi(m.state.MidiHandler)

	case error:
		m.err = msg
		return m, nil
	}

	return m, nil
}

// handleKey processes keyboard input
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewMixer:
		return m.handleMixerKeys(msg)
	case ViewDevices:
		return m.handleDeviceKeys(msg)
	}
	return m, nil
}

// handleMixerKeys handles keyboard input in mixer view
func (m Model) handleMixerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.state.Close()
		return m, tea.Quit

	case "left", "h":
		m.state.SelectPrev()

	case "right", "l":
		m.state.SelectNext()

	case "up", "k":
		m.state.AdjustVolume(5)

	case "down", "j":
		m.state.AdjustVolume(-5)

	case "shift+up", "K":
		m.state.AdjustVolume(1)

	case "shift+down", "J":
		m.state.AdjustVolume(-1)

	case "[":
		m.state.AdjustPan(-5)

	case "]":
		m.state.AdjustPan(5)

	case "{":
		m.state.AdjustPan(-1)

	case "}":
		m.state.AdjustPan(1)

	case "m":
		m.state.ToggleMute()

	case "s":
		m.state.ToggleSolo()

	case "d":
		m.deviceSelector = ui.NewDeviceSelector()
		m.currentView = ViewDevices

	case "0":
		// Reset selected channel to defaults
		if ch := m.state.SelectedChannel(); ch != nil {
			ch.Volume = 100
			ch.Pan = 64
			ch.Mute = false
			ch.Solo = false
		}
	}

	return m, nil
}

// handleDeviceKeys handles keyboard input in device selection view
func (m Model) handleDeviceKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.state.Close()
		return m, tea.Quit

	case "esc":
		m.currentView = ViewMixer

	case "up", "k":
		m.deviceSelector.MoveUp()

	case "down", "j":
		m.deviceSelector.MoveDown()

	case "tab":
		m.deviceSelector.ToggleFocus()

	case "r":
		m.deviceSelector.Refresh()

	case "enter":
		inPort := m.deviceSelector.GetSelectedInput()
		outPort := m.deviceSelector.GetSelectedOutput()
		if err := m.state.MidiHandler.Connect(inPort, outPort); err != nil {
			m.err = err
		}
		m.currentView = ViewMixer
	}

	return m, nil
}

// handleMidiCC processes incoming MIDI CC messages
func (m *Model) handleMidiCC(msg midi.CCMessage) {
	// Map MIDI channel to mixer channel
	chIdx := int(msg.Channel)
	if chIdx >= len(m.state.Channels) {
		return
	}

	switch msg.Controller {
	case midi.CCVolume:
		m.state.SetChannelVolume(chIdx, msg.Value)
	case midi.CCPan:
		m.state.SetChannelPan(chIdx, msg.Value)
	}
}

// View renders the current view
func (m Model) View() string {
	var content string

	switch m.currentView {
	case ViewMixer:
		content = m.renderMixerView()
	case ViewDevices:
		content = m.renderDevicesView()
	}

	// Center content
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// renderMixerView renders the main mixer interface
func (m Model) renderMixerView() string {
	var sections []string

	// Title
	title := ui.TitleStyle.Render("🎛️  MIDI MIXER")
	sections = append(sections, title)

	// Error message if any
	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
		sections = append(sections, errStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	// Mixer channels
	sections = append(sections, ui.RenderMixer(m.state))

	// Status bar
	sections = append(sections, ui.RenderStatus(m.state))

	// Help
	sections = append(sections, ui.RenderHelp())

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

// renderDevicesView renders the device selection interface
func (m Model) renderDevicesView() string {
	return ui.RenderDeviceSelector(m.deviceSelector)
}

func main() {
	// Create initial state with 8 channels
	state := mixer.NewState(8)

	// Create model
	model := Model{
		state:       state,
		currentView: ViewMixer,
	}

	// Run the program
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
