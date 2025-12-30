package audio

import (
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

const (
	sampleRate   = 44100
	channelCount = 2
	bitDepth     = 2
	waveformSize = 128
	MinBPM       = 60
	MaxBPM       = 200
	DefaultBPM   = 120
)

// BeatPreset contains patterns for all drums
type BeatPreset struct {
	Name        string
	Description string
	Kick        []int
	Snare       []int
	HiHat       []int
	Bass        []int
}

// Sick beat presets for different vibes
var BeatPresets = []BeatPreset{
	{
		Name:        "🔥 Trap Fire",
		Description: "Hard-hitting trap beat with rolling hi-hats",
		Kick:        []int{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		Snare:       []int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		HiHat:       []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		Bass:        []int{1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 1},
	},
	{
		Name:        "🎸 Rock Solid",
		Description: "Classic rock beat - simple but powerful",
		Kick:        []int{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
		Snare:       []int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		HiHat:       []int{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
		Bass:        []int{1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		Name:        "🕺 Disco Funk",
		Description: "Groovy disco vibes with syncopated rhythm",
		Kick:        []int{1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0},
		Snare:       []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0},
		HiHat:       []int{0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
		Bass:        []int{1, 0, 1, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 1},
	},
	{
		Name:        "🌊 Lo-Fi Chill",
		Description: "Relaxed, laid-back beats to study to",
		Kick:        []int{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		Snare:       []int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
		HiHat:       []int{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		Bass:        []int{1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0},
	},
	{
		Name:        "🎹 House Party",
		Description: "Four-on-the-floor house music energy",
		Kick:        []int{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
		Snare:       []int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		HiHat:       []int{0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		Bass:        []int{1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1},
	},
	{
		Name:        "💀 Dubstep Drop",
		Description: "Heavy wobbles and aggressive rhythms",
		Kick:        []int{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
		Snare:       []int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
		HiHat:       []int{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1},
		Bass:        []int{1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0},
	},
	{
		Name:        "🥁 Drum & Bass",
		Description: "Fast-paced jungle rhythms",
		Kick:        []int{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		Snare:       []int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		HiHat:       []int{1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1},
		Bass:        []int{1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0},
	},
	{
		Name:        "🎺 Latin Heat",
		Description: "Salsa-inspired rhythm with clave pattern",
		Kick:        []int{1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0},
		Snare:       []int{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		HiHat:       []int{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
		Bass:        []int{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	},
}

// Channel types
const (
	ChKick = iota
	ChSnare
	ChHiHat
	ChBass
	ChLead1
	ChLead2
	ChPad
	ChFX
)

// Channel names for display
var ChannelNames = []string{"KICK", "SNARE", "HIHAT", "BASS", "LEAD1", "LEAD2", "PAD", "FX"}

type Engine struct {
	ctx          *oto.Context
	player       oto.Player
	mu           sync.RWMutex
	channels     []ChannelState
	master       float64
	running      bool
	samplePos    int64
	waveformL    []float64
	waveformR    []float64
	waveformIdx  int
	waveformMu   sync.RWMutex
	envelopes    []float64
	noisePhase   float64
	bassPhase    float64
	leadPhases   []float64
	padPhases    []float64
	BPM          int
	PatternIndex int
	CurrentStep  int
}

type ChannelState struct {
	Volume    float64
	Pan       float64
	Mute      bool
	Solo      bool
	Frequency float64
}

type audioStream struct {
	engine *Engine
}

func NewEngine(numChannels int) (*Engine, error) {
	ctx, ready, err := oto.NewContext(sampleRate, channelCount, bitDepth)
	if err != nil {
		return nil, err
	}
	<-ready

	channels := make([]ChannelState, numChannels)
	envelopes := make([]float64, numChannels)

	// Initialize channel defaults
	defaults := []struct {
		vol  float64
		freq float64
	}{
		{0.9, 60},   // Kick
		{0.7, 200},  // Snare
		{0.5, 8000}, // HiHat
		{0.8, 55},   // Bass
		{0.6, 440},  // Lead1
		{0.5, 523},  // Lead2
		{0.4, 220},  // Pad
		{0.3, 1000}, // FX
	}

	for i := 0; i < numChannels; i++ {
		vol, freq := 0.7, 440.0
		if i < len(defaults) {
			vol = defaults[i].vol
			freq = defaults[i].freq
		}
		channels[i] = ChannelState{
			Volume:    vol,
			Pan:       0.0,
			Mute:      false,
			Solo:      false,
			Frequency: freq,
		}
	}

	e := &Engine{
		ctx:          ctx,
		channels:     channels,
		master:       0.8,
		running:      true,
		waveformL:    make([]float64, waveformSize),
		waveformR:    make([]float64, waveformSize),
		envelopes:    envelopes,
		leadPhases:   make([]float64, 2),
		padPhases:    make([]float64, 4),
		BPM:          DefaultBPM,
		PatternIndex: 0,
		CurrentStep:  0,
	}

	e.player = ctx.NewPlayer(&audioStream{engine: e})
	e.player.Play()

	return e, nil
}

func (s *audioStream) Read(buf []byte) (int, error) {
	s.engine.mu.RLock()
	running := s.engine.running
	master := s.engine.master
	channels := make([]ChannelState, len(s.engine.channels))
	copy(channels, s.engine.channels)
	envelopes := make([]float64, len(s.engine.envelopes))
	copy(envelopes, s.engine.envelopes)
	s.engine.mu.RUnlock()

	if !running {
		for i := range buf {
			buf[i] = 0
		}
		return len(buf), nil
	}

	anySolo := false
	for _, ch := range channels {
		if ch.Solo {
			anySolo = true
			break
		}
	}

	// Get current BPM and pattern
	s.engine.mu.RLock()
	bpm := s.engine.BPM
	patternIdx := s.engine.PatternIndex
	s.engine.mu.RUnlock()

	if patternIdx >= len(BeatPresets) {
		patternIdx = 0
	}
	pattern := BeatPresets[patternIdx]

	samplesPerBeat := sampleRate * 60 / bpm / 4 // 16th notes

	samples := len(buf) / 4
	for i := 0; i < samples; i++ {
		s.engine.mu.Lock()
		samplePos := s.engine.samplePos
		s.engine.samplePos++

		step := int(samplePos/int64(samplesPerBeat)) % 16
		s.engine.CurrentStep = step
		stepProgress := float64(samplePos%int64(samplesPerBeat)) / float64(samplesPerBeat)

		// Trigger envelopes on beat
		if samplePos%int64(samplesPerBeat) == 0 {
			if pattern.Kick[step] == 1 {
				s.engine.envelopes[ChKick] = 1.0
			}
			if pattern.Snare[step] == 1 {
				s.engine.envelopes[ChSnare] = 1.0
			}
			if pattern.HiHat[step] == 1 {
				s.engine.envelopes[ChHiHat] = 1.0
			}
			if pattern.Bass[step] == 1 {
				s.engine.envelopes[ChBass] = 1.0
			}
		}

		// Decay envelopes
		for j := range s.engine.envelopes {
			s.engine.envelopes[j] *= 0.9997
		}

		var leftSum, rightSum float64

		// Generate each channel
		for chIdx := range channels {
			ch := channels[chIdx]
			if ch.Mute || (anySolo && !ch.Solo) {
				continue
			}

			var sample float64
			env := s.engine.envelopes[chIdx]

			switch chIdx {
			case ChKick:
				// Kick: pitch-dropping sine
				kickFreq := 150*math.Exp(-5*stepProgress) + 40
				s.engine.bassPhase += 2 * math.Pi * kickFreq / sampleRate
				sample = math.Sin(s.engine.bassPhase) * env * 1.2

			case ChSnare:
				// Snare: noise + tone
				s.engine.noisePhase += 0.1
				noise := (rand.Float64()*2 - 1) * 0.6
				tone := math.Sin(s.engine.noisePhase*200) * 0.4
				sample = (noise + tone) * env

			case ChHiHat:
				// HiHat: filtered noise
				noise := rand.Float64()*2 - 1
				sample = noise * env * 0.5 * math.Exp(-10*stepProgress)

			case ChBass:
				// Bass: saw-ish wave
				bassFreq := ch.Frequency
				t := float64(samplePos) / sampleRate
				saw := 2*math.Mod(t*bassFreq, 1) - 1
				sample = saw * env * 0.7

			case ChLead1, ChLead2:
				// Lead: detuned saws
				idx := chIdx - ChLead1
				freq := ch.Frequency * (1 + float64(step%4)*0.02)
				s.engine.leadPhases[idx] += 2 * math.Pi * freq / sampleRate
				sample = math.Sin(s.engine.leadPhases[idx]) * 0.5
				sample += math.Sin(s.engine.leadPhases[idx]*2.01) * 0.25

			case ChPad:
				// Pad: soft chord
				freqs := []float64{ch.Frequency, ch.Frequency * 1.25, ch.Frequency * 1.5, ch.Frequency * 2}
				for pi, f := range freqs {
					s.engine.padPhases[pi] += 2 * math.Pi * f / sampleRate
					sample += math.Sin(s.engine.padPhases[pi]) * 0.15
				}

			case ChFX:
				// FX: filtered sweep
				sweep := math.Sin(float64(samplePos) * 0.0001)
				sample = math.Sin(float64(samplePos)*0.05*(1+sweep*0.5)) * 0.3
			}

			sample *= ch.Volume

			// Panning
			angle := (ch.Pan + 1) * math.Pi / 4
			leftSum += sample * math.Cos(angle)
			rightSum += sample * math.Sin(angle)
		}

		s.engine.mu.Unlock()

		leftSum *= master
		rightSum *= master
		leftSum = softClip(leftSum)
		rightSum = softClip(rightSum)

		// Store waveform for visualization
		s.engine.waveformMu.Lock()
		s.engine.waveformL[s.engine.waveformIdx] = leftSum
		s.engine.waveformR[s.engine.waveformIdx] = rightSum
		s.engine.waveformIdx = (s.engine.waveformIdx + 1) % waveformSize
		s.engine.waveformMu.Unlock()

		leftInt := int16(leftSum * 32767 * 0.7)
		rightInt := int16(rightSum * 32767 * 0.7)

		idx := i * 4
		buf[idx] = byte(leftInt)
		buf[idx+1] = byte(leftInt >> 8)
		buf[idx+2] = byte(rightInt)
		buf[idx+3] = byte(rightInt >> 8)
	}

	return len(buf), nil
}

func softClip(x float64) float64 {
	if x > 1 {
		return 1
	}
	if x < -1 {
		return -1
	}
	return 1.5*x - 0.5*x*x*x
}

// GetWaveform returns current waveform data for visualization
func (e *Engine) GetWaveform() ([]float64, []float64) {
	e.waveformMu.RLock()
	defer e.waveformMu.RUnlock()

	left := make([]float64, waveformSize)
	right := make([]float64, waveformSize)

	// Copy in order starting from current index
	for i := 0; i < waveformSize; i++ {
		idx := (e.waveformIdx + i) % waveformSize
		left[i] = e.waveformL[idx]
		right[i] = e.waveformR[idx]
	}
	return left, right
}

func (e *Engine) SetChannelVolume(channel int, value uint8) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if channel >= 0 && channel < len(e.channels) {
		e.channels[channel].Volume = float64(value) / 127.0
	}
}

func (e *Engine) SetChannelPan(channel int, value uint8) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if channel >= 0 && channel < len(e.channels) {
		e.channels[channel].Pan = (float64(value) - 64) / 64.0
	}
}

func (e *Engine) SetChannelMute(channel int, muted bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if channel >= 0 && channel < len(e.channels) {
		e.channels[channel].Mute = muted
	}
}

func (e *Engine) SetChannelSolo(channel int, solo bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if channel >= 0 && channel < len(e.channels) {
		e.channels[channel].Solo = solo
	}
}

func (e *Engine) SetMasterVolume(value uint8) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.master = float64(value) / 127.0
}

// SetBPM sets the tempo in beats per minute
func (e *Engine) SetBPM(bpm int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if bpm < MinBPM {
		bpm = MinBPM
	}
	if bpm > MaxBPM {
		bpm = MaxBPM
	}
	e.BPM = bpm
}

// GetBPM returns current tempo
func (e *Engine) GetBPM() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.BPM
}

// SetPattern sets the current beat pattern
func (e *Engine) SetPattern(index int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if index >= 0 && index < len(BeatPresets) {
		e.PatternIndex = index
	}
}

// GetPattern returns current pattern index
func (e *Engine) GetPattern() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.PatternIndex
}

// GetCurrentStep returns the current step (0-15)
func (e *Engine) GetCurrentStep() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.CurrentStep
}

// NextPattern cycles to the next pattern
func (e *Engine) NextPattern() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.PatternIndex = (e.PatternIndex + 1) % len(BeatPresets)
}

// PrevPattern cycles to the previous pattern
func (e *Engine) PrevPattern() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.PatternIndex--
	if e.PatternIndex < 0 {
		e.PatternIndex = len(BeatPresets) - 1
	}
}

func (e *Engine) Close() {
	e.mu.Lock()
	e.running = false
	e.mu.Unlock()
	if e.player != nil {
		e.player.Close()
	}
}
