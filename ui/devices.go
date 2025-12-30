package ui

import (
	"fmt"
	"strings"

	"gitlab.com/gomidi/midi/v2/drivers"

	"midi-mixer/midi"
)

// DeviceSelector handles device selection UI
type DeviceSelector struct {
	InputPorts     []drivers.In
	OutputPorts    []drivers.Out
	SelectedInput  int
	SelectedOutput int
	FocusInput     bool // true = input list focused, false = output list
}

// NewDeviceSelector creates a new device selector
func NewDeviceSelector() *DeviceSelector {
	return &DeviceSelector{
		InputPorts:     midi.GetInputPorts(),
		OutputPorts:    midi.GetOutputPorts(),
		SelectedInput:  -1,
		SelectedOutput: -1,
		FocusInput:     true,
	}
}

// Refresh reloads available MIDI ports
func (d *DeviceSelector) Refresh() {
	d.InputPorts = midi.GetInputPorts()
	d.OutputPorts = midi.GetOutputPorts()
}

// MoveUp moves selection up in current list
func (d *DeviceSelector) MoveUp() {
	if d.FocusInput {
		if d.SelectedInput > 0 {
			d.SelectedInput--
		} else if d.SelectedInput == -1 && len(d.InputPorts) > 0 {
			d.SelectedInput = 0
		}
	} else {
		if d.SelectedOutput > 0 {
			d.SelectedOutput--
		} else if d.SelectedOutput == -1 && len(d.OutputPorts) > 0 {
			d.SelectedOutput = 0
		}
	}
}

// MoveDown moves selection down in current list
func (d *DeviceSelector) MoveDown() {
	if d.FocusInput {
		if d.SelectedInput < len(d.InputPorts)-1 {
			d.SelectedInput++
		}
	} else {
		if d.SelectedOutput < len(d.OutputPorts)-1 {
			d.SelectedOutput++
		}
	}
}

// ToggleFocus switches between input and output lists
func (d *DeviceSelector) ToggleFocus() {
	d.FocusInput = !d.FocusInput
}

// GetSelectedInput returns the selected input port or nil
func (d *DeviceSelector) GetSelectedInput() drivers.In {
	if d.SelectedInput >= 0 && d.SelectedInput < len(d.InputPorts) {
		return d.InputPorts[d.SelectedInput]
	}
	return nil
}

// GetSelectedOutput returns the selected output port or nil
func (d *DeviceSelector) GetSelectedOutput() drivers.Out {
	if d.SelectedOutput >= 0 && d.SelectedOutput < len(d.OutputPorts) {
		return d.OutputPorts[d.SelectedOutput]
	}
	return nil
}

// RenderDeviceSelector renders the device selection view
func RenderDeviceSelector(d *DeviceSelector) string {
	var sections []string

	// Title
	sections = append(sections, TitleStyle.Render("🎹 MIDI Device Selection"))
	sections = append(sections, "")

	// Input ports list
	inputTitle := "Input Ports"
	if d.FocusInput {
		inputTitle = "▸ Input Ports"
	}
	sections = append(sections, ChannelNameStyle.Render(inputTitle))

	if len(d.InputPorts) == 0 {
		sections = append(sections, DeviceItemStyle.Render("  No input devices found"))
	} else {
		for i, port := range d.InputPorts {
			name := port.String()
			if i == d.SelectedInput {
				if d.FocusInput {
					sections = append(sections, DeviceSelectedStyle.Render(fmt.Sprintf("● %s", name)))
				} else {
					sections = append(sections, DeviceItemStyle.Render(fmt.Sprintf("● %s", name)))
				}
			} else {
				sections = append(sections, DeviceItemStyle.Render(fmt.Sprintf("  %s", name)))
			}
		}
	}

	sections = append(sections, "")

	// Output ports list
	outputTitle := "Output Ports"
	if !d.FocusInput {
		outputTitle = "▸ Output Ports"
	}
	sections = append(sections, ChannelNameStyle.Render(outputTitle))

	if len(d.OutputPorts) == 0 {
		sections = append(sections, DeviceItemStyle.Render("  No output devices found"))
	} else {
		for i, port := range d.OutputPorts {
			name := port.String()
			if i == d.SelectedOutput {
				if !d.FocusInput {
					sections = append(sections, DeviceSelectedStyle.Render(fmt.Sprintf("● %s", name)))
				} else {
					sections = append(sections, DeviceItemStyle.Render(fmt.Sprintf("● %s", name)))
				}
			} else {
				sections = append(sections, DeviceItemStyle.Render(fmt.Sprintf("  %s", name)))
			}
		}
	}

	sections = append(sections, "")
	sections = append(sections, HelpStyle.Render("↑/↓: Select  Tab: Switch List  Enter: Connect  R: Refresh  Esc: Cancel"))

	content := strings.Join(sections, "\n")
	return DeviceListStyle.Render(content)
}
