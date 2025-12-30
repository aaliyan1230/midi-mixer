package ui

import (
	"fmt"
	"math"
	"strings"

	"midi-mixer/mixer"

	"github.com/charmbracelet/lipgloss"
)

const (
	FaderHeight    = 10 // Number of rows for fader display
	WaveformWidth  = 80
	WaveformHeight = 8
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

	// Channel name - truncate if too long
	name := ch.Name
	if len(name) > 6 {
		name = name[:6]
	}
	parts = append(parts, ChannelNameStyle.Render(name))
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

// Waveform block characters for different amplitudes
var waveBlocks = []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// RenderWaveform renders a stereo waveform oscilloscope
func RenderWaveform(leftWave, rightWave []float64) string {
	if len(leftWave) == 0 || len(rightWave) == 0 {
		return ""
	}

	width := WaveformWidth
	height := WaveformHeight

	// Downsample waveform to fit width
	step := len(leftWave) / width
	if step < 1 {
		step = 1
	}

	// Build waveform display
	var lines []string

	// Header
	headerStyle := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	lines = append(lines, headerStyle.Render("┌─ WAVEFORM ─────────────────────────────────────────────────────────────────┐"))

	// Create display buffer
	display := make([][]string, height)
	for i := range display {
		display[i] = make([]string, width)
		for j := range display[i] {
			display[i][j] = " "
		}
	}

	// Left channel (top half) - cyan
	leftStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#06B6D4"))
	// Right channel (bottom half) - magenta
	rightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D946EF"))

	halfHeight := height / 2

	for x := 0; x < width && x*step < len(leftWave); x++ {
		// Get sample values
		lSample := leftWave[x*step]
		rSample := rightWave[x*step]

		// Scale to display coordinates
		lY := int((1 - lSample) * float64(halfHeight-1))
		rY := halfHeight + int((1-rSample)*float64(halfHeight-1))

		if lY < 0 {
			lY = 0
		}
		if lY >= halfHeight {
			lY = halfHeight - 1
		}
		if rY < halfHeight {
			rY = halfHeight
		}
		if rY >= height {
			rY = height - 1
		}

		// Draw points
		display[lY][x] = "L"
		display[rY][x] = "R"
	}

	// Render display to strings
	for y := 0; y < height; y++ {
		var line strings.Builder
		line.WriteString("│")
		for x := 0; x < width; x++ {
			char := display[y][x]
			switch char {
			case "L":
				line.WriteString(leftStyle.Render("█"))
			case "R":
				line.WriteString(rightStyle.Render("█"))
			default:
				if y == halfHeight-1 || y == halfHeight {
					line.WriteString(lipgloss.NewStyle().Foreground(ColorSurface).Render("─"))
				} else {
					line.WriteString(" ")
				}
			}
		}
		line.WriteString("│")
		lines = append(lines, line.String())
	}

	// Footer with labels
	footerStyle := lipgloss.NewStyle().Foreground(ColorTextDim)
	lines = append(lines, footerStyle.Render("└─ ")+leftStyle.Render("■ LEFT")+footerStyle.Render("  ")+rightStyle.Render("■ RIGHT")+footerStyle.Render(" ──────────────────────────────────────────────────────────┘"))

	return strings.Join(lines, "\n")
}

// RenderVUMeter renders a horizontal VU meter
func RenderVUMeter(leftWave, rightWave []float64) string {
	// Calculate RMS levels
	var leftRMS, rightRMS float64
	for i := range leftWave {
		leftRMS += leftWave[i] * leftWave[i]
		rightRMS += rightWave[i] * rightWave[i]
	}
	if len(leftWave) > 0 {
		leftRMS = math.Sqrt(leftRMS / float64(len(leftWave)))
		rightRMS = math.Sqrt(rightRMS / float64(len(rightWave)))
	}

	width := 40
	leftBars := int(leftRMS * float64(width) * 2)
	rightBars := int(rightRMS * float64(width) * 2)
	if leftBars > width {
		leftBars = width
	}
	if rightBars > width {
		rightBars = width
	}

	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	dimStyle := lipgloss.NewStyle().Foreground(ColorSurface)

	renderBar := func(level int) string {
		var bar strings.Builder
		for i := 0; i < width; i++ {
			if i < level {
				if i < width*6/10 {
					bar.WriteString(greenStyle.Render("█"))
				} else if i < width*8/10 {
					bar.WriteString(yellowStyle.Render("█"))
				} else {
					bar.WriteString(redStyle.Render("█"))
				}
			} else {
				bar.WriteString(dimStyle.Render("░"))
			}
		}
		return bar.String()
	}

	leftLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("#06B6D4")).Render("L ")
	rightLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("#D946EF")).Render("R ")

	return leftLabel + renderBar(leftBars) + "\n" + rightLabel + renderBar(rightBars)
}
