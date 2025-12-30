package mixer

import (
	"midi-mixer/audio"
	"midi-mixer/midi"
)

// Channel represents a single mixer channel strip
type Channel struct {
	ID     int
	Name   string
	Volume uint8 // 0-127, mapped to CC7
	Pan    uint8 // 0-127 (64=center), mapped to CC10
	Mute   bool
	Solo   bool
}

// NewChannel creates a new mixer channel with default values
func NewChannel(id int, name string) Channel {
	return Channel{
		ID:     id,
		Name:   name,
		Volume: 100, // ~79% default
		Pan:    64,  // Center
		Mute:   false,
		Solo:   false,
	}
}

// State holds the complete mixer state
type State struct {
	Channels      []Channel
	MasterVolume  uint8
	SelectedIndex int
	MidiHandler   *midi.Handler
	AudioEngine   *audio.Engine
	InputPortIdx  int
	OutputPortIdx int
}

// NewState creates a new mixer state with 8 channels
func NewState(numChannels int) *State {
	channels := make([]Channel, numChannels)
	for i := 0; i < numChannels; i++ {
		channels[i] = NewChannel(i, channelName(i))
	}

	// Mute FX channel by default (it can sound harsh)
	if numChannels > audio.ChFX {
		channels[audio.ChFX].Mute = true
	}

	// Initialize audio engine
	audioEngine, _ := audio.NewEngine(numChannels)

	state := &State{
		Channels:      channels,
		MasterVolume:  100,
		SelectedIndex: 0,
		MidiHandler:   midi.NewHandler(),
		AudioEngine:   audioEngine,
		InputPortIdx:  -1,
		OutputPortIdx: -1,
	}

	// Sync initial state to audio engine
	if audioEngine != nil {
		for i, ch := range channels {
			audioEngine.SetChannelVolume(i, ch.Volume)
			audioEngine.SetChannelPan(i, ch.Pan)
			audioEngine.SetChannelMute(i, ch.Mute)
		}
		audioEngine.SetMasterVolume(state.MasterVolume)
	}

	return state
}

// channelName returns default name for channel index
func channelName(idx int) string {
	if idx < len(audio.ChannelNames) {
		return audio.ChannelNames[idx]
	}
	return string(rune('1' + idx))
}

// SelectedChannel returns the currently selected channel
func (s *State) SelectedChannel() *Channel {
	if s.SelectedIndex >= 0 && s.SelectedIndex < len(s.Channels) {
		return &s.Channels[s.SelectedIndex]
	}
	return nil
}

// SelectNext moves selection to the next channel
func (s *State) SelectNext() {
	if s.SelectedIndex < len(s.Channels)-1 {
		s.SelectedIndex++
	}
}

// SelectPrev moves selection to the previous channel
func (s *State) SelectPrev() {
	if s.SelectedIndex > 0 {
		s.SelectedIndex--
	}
}

// AdjustVolume changes the selected channel's volume
func (s *State) AdjustVolume(delta int) {
	ch := s.SelectedChannel()
	if ch == nil {
		return
	}

	newVal := int(ch.Volume) + delta
	if newVal < 0 {
		newVal = 0
	} else if newVal > 127 {
		newVal = 127
	}
	ch.Volume = uint8(newVal)

	// Update audio engine
	if s.AudioEngine != nil {
		s.AudioEngine.SetChannelVolume(ch.ID, ch.Volume)
	}

	// Send MIDI CC if not muted
	if !ch.Mute && s.MidiHandler != nil {
		s.MidiHandler.SendCC(uint8(ch.ID), midi.CCVolume, ch.Volume)
	}
}

// AdjustPan changes the selected channel's pan
func (s *State) AdjustPan(delta int) {
	ch := s.SelectedChannel()
	if ch == nil {
		return
	}

	newVal := int(ch.Pan) + delta
	if newVal < 0 {
		newVal = 0
	} else if newVal > 127 {
		newVal = 127
	}
	ch.Pan = uint8(newVal)

	// Update audio engine
	if s.AudioEngine != nil {
		s.AudioEngine.SetChannelPan(ch.ID, ch.Pan)
	}

	// Send MIDI CC
	if s.MidiHandler != nil {
		s.MidiHandler.SendCC(uint8(ch.ID), midi.CCPan, ch.Pan)
	}
}

// ToggleMute toggles mute on the selected channel
func (s *State) ToggleMute() {
	ch := s.SelectedChannel()
	if ch == nil {
		return
	}

	ch.Mute = !ch.Mute

	// Update audio engine
	if s.AudioEngine != nil {
		s.AudioEngine.SetChannelMute(ch.ID, ch.Mute)
	}

	// Send volume 0 when muted, restore when unmuted
	if s.MidiHandler != nil {
		if ch.Mute {
			s.MidiHandler.SendCC(uint8(ch.ID), midi.CCVolume, 0)
		} else {
			s.MidiHandler.SendCC(uint8(ch.ID), midi.CCVolume, ch.Volume)
		}
	}
}

// ToggleSolo toggles solo on the selected channel
func (s *State) ToggleSolo() {
	ch := s.SelectedChannel()
	if ch == nil {
		return
	}

	ch.Solo = !ch.Solo

	// Update audio engine for all channels
	if s.AudioEngine != nil {
		for _, c := range s.Channels {
			s.AudioEngine.SetChannelSolo(c.ID, c.Solo)
		}
	}

	s.updateSoloState()
}

// updateSoloState handles solo logic (mutes non-soloed channels when any solo is active)
func (s *State) updateSoloState() {
	// Check if any channel is soloed
	anySolo := false
	for _, ch := range s.Channels {
		if ch.Solo {
			anySolo = true
			break
		}
	}

	if s.MidiHandler == nil {
		return
	}

	// Update MIDI output based on solo state
	for _, ch := range s.Channels {
		var volume uint8 = 0
		if !anySolo {
			// No solo active, use normal mute logic
			if !ch.Mute {
				volume = ch.Volume
			}
		} else {
			// Solo active, only soloed channels are audible
			if ch.Solo && !ch.Mute {
				volume = ch.Volume
			}
		}
		s.MidiHandler.SendCC(uint8(ch.ID), midi.CCVolume, volume)
	}
}

// SetChannelVolume sets volume for a specific channel (used for incoming MIDI)
func (s *State) SetChannelVolume(channelID int, value uint8) {
	if channelID >= 0 && channelID < len(s.Channels) {
		s.Channels[channelID].Volume = value
		if s.AudioEngine != nil {
			s.AudioEngine.SetChannelVolume(channelID, value)
		}
	}
}

// SetChannelPan sets pan for a specific channel (used for incoming MIDI)
func (s *State) SetChannelPan(channelID int, value uint8) {
	if channelID >= 0 && channelID < len(s.Channels) {
		s.Channels[channelID].Pan = value
		if s.AudioEngine != nil {
			s.AudioEngine.SetChannelPan(channelID, value)
		}
	}
}

// AdjustMasterVolume changes the master volume
func (s *State) AdjustMasterVolume(delta int) {
	newVal := int(s.MasterVolume) + delta
	if newVal < 0 {
		newVal = 0
	} else if newVal > 127 {
		newVal = 127
	}
	s.MasterVolume = uint8(newVal)

	// Update audio engine
	if s.AudioEngine != nil {
		s.AudioEngine.SetMasterVolume(s.MasterVolume)
	}
}

// AdjustBPM changes the tempo
func (s *State) AdjustBPM(delta int) {
	if s.AudioEngine != nil {
		newBPM := s.AudioEngine.GetBPM() + delta
		s.AudioEngine.SetBPM(newBPM)
	}
}

// GetBPM returns current BPM
func (s *State) GetBPM() int {
	if s.AudioEngine != nil {
		return s.AudioEngine.GetBPM()
	}
	return audio.DefaultBPM
}

// NextPattern cycles to next beat pattern
func (s *State) NextPattern() {
	if s.AudioEngine != nil {
		s.AudioEngine.NextPattern()
	}
}

// PrevPattern cycles to previous beat pattern
func (s *State) PrevPattern() {
	if s.AudioEngine != nil {
		s.AudioEngine.PrevPattern()
	}
}

// GetPatternIndex returns current pattern index
func (s *State) GetPatternIndex() int {
	if s.AudioEngine != nil {
		return s.AudioEngine.GetPattern()
	}
	return 0
}

// GetCurrentStep returns the current step (0-15)
func (s *State) GetCurrentStep() int {
	if s.AudioEngine != nil {
		return s.AudioEngine.GetCurrentStep()
	}
	return 0
}

// Close cleans up resources
func (s *State) Close() {
	if s.AudioEngine != nil {
		s.AudioEngine.Close()
	}
	if s.MidiHandler != nil {
		s.MidiHandler.Close()
	}
}
