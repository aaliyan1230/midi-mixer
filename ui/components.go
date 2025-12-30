package ui

import (
	"fmt"
	"strings"

	"midi-mixer/mixer"

	"github.com/charmbracelet/lipgloss"
)

const (
	FaderHeight = 12 // Number of rows for fader display
)

// RenderFader renders a vertical fader for a value 0-127
func RenderFader(value uint8, height int) string {
	// Calculate filled blocks
	filled := int(float64(value) / 127.0 * float64(height))

	var lines []string
	for i := height - 1; i >= 0; i-- {
		if i < filled {
			lines = append(lines, FaderFillStyle.Render("██"))
		} else {
			lines = append(lines, FaderTrackStyle.Render("░░"))
		}
	}

	return strings.Join(lines, "\n")
}

// RenderPanKnob renders a simple pan indicator
func RenderPanKnob(pan uint8) string {
	// Convert 0-127 to position indicator
	// 0 = full left, 64 = center, 127 = full right
	const width = 7
	pos := int(float64(pan) / 127.0 * float64(width-1))

	indicator := strings.Repeat("─", pos) + "●" + strings.Repeat("─", width-1-pos)

	label := "C"
	if pan < 54 {
		label = fmt.Sprintf("L%d", (64-int(pan))*100/64)
	} else if pan > 74 {
		label = fmt.Sprintf("R%d", (int(pan)-64)*100/63)
	}

	return PanStyle.Render(fmt.Sprintf("[%s]\n %s", indicator, label))
}

// RenderChannel renders a single channel strip
func RenderChannel(ch mixer.Channel, selected bool) string {
	var parts []string

	// Channel name
	parts = append(parts, ChannelNameStyle.Render(fmt.Sprintf("Ch %s", ch.Name)))
	parts = append(parts, "")

	// Fader
	parts = append(parts, RenderFader(ch.Volume, FaderHeight))

	// Volume value
	volPercent := int(float64(ch.Volume) / 127.0 * 100)
	parts = append(parts, ValueStyle.Render(fmt.Sprintf("%3d%%", volPercent)))
	parts = append(parts, "")

	// Pan
	parts = append(parts, RenderPanKnob(ch.Pan))
	parts = append(parts, "")

	// Mute/Solo buttons
	var muteStr, soloStr string
	if ch.Mute {
		muteStr = MuteActiveStyle.Render("M")
	} else {
		muteStr = MuteInactiveStyle.Render("M")
	}
	if ch.Solo {
		soloStr = SoloActiveStyle.Render("S")
	} else {
		soloStr = SoloInactiveStyle.Render("S")
	}
	parts = append(parts, fmt.Sprintf("%s %s", muteStr, soloStr))

	content := strings.Join(parts, "\n")

	if selected {
		return SelectedChannelStyle.Render(content)
	}
	return ChannelStyle.Render(content)
}

// RenderMasterFader renders the master volume fader
func RenderMasterFader(volume uint8) string {
	var parts []string

	parts = append(parts, ChannelNameStyle.Render("MASTER"))
	parts = append(parts, "")
	parts = append(parts, RenderFader(volume, FaderHeight))

	volPercent := int(float64(volume) / 127.0 * 100)
	parts = append(parts, ValueStyle.Render(fmt.Sprintf("%3d%%", volPercent)))

	return MasterStyle.Render(strings.Join(parts, "\n"))
}

// RenderMixer renders the complete mixer view
func RenderMixer(state *mixer.State) string {
	var channelViews []string

	// Render each channel
	for i, ch := range state.Channels {
		channelViews = append(channelViews, RenderChannel(ch, i == state.SelectedIndex))
	}

	// Add master fader
	channelViews = append(channelViews, RenderMasterFader(state.MasterVolume))

	// Join channels horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, channelViews...)
}

// RenderHelp renders the help bar
func RenderHelp() string {
	help := "←/→: Select  ↑/↓: Volume  [/]: Pan  M: Mute  S: Solo  D: Devices  Q: Quit"
	return HelpStyle.Render(help)
}

// RenderStatus renders the status bar with MIDI info
func RenderStatus(state *mixer.State) string {
	inPort := state.MidiHandler.GetInputPortName()
	outPort := state.MidiHandler.GetOutputPortName()

	status := fmt.Sprintf("MIDI In: %s │ MIDI Out: %s", inPort, outPort)
	return StatusStyle.Render(status)
}
