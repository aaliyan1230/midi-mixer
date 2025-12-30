package audio

import (
	"math"
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

const (
	sampleRate   = 44100
	channelCount = 2
	bitDepth     = 2
)

var channelFrequencies = []float64{
	261.63, 293.66, 329.63, 349.23,
	392.00, 440.00, 493.88, 523.25,
}

type Engine struct {
	ctx      *oto.Context
	player   oto.Player
	mu       sync.RWMutex
	channels []ChannelState
	master   float64
	running  bool
	phase    []float64
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
	phase := make([]float64, numChannels)
	for i := 0; i < numChannels; i++ {
		freq := 440.0
		if i < len(channelFrequencies) {
			freq = channelFrequencies[i]
		}
		channels[i] = ChannelState{
			Volume:    0.7,
			Pan:       0.0,
			Mute:      false,
			Solo:      false,
			Frequency: freq,
		}
	}

	e := &Engine{
		ctx:      ctx,
		channels: channels,
		master:   0.7,
		phase:    phase,
		running:  true,
	}

	e.player = ctx.NewPlayer(&audioStream{engine: e})
	e.player.Play()

	return e, nil
}

func (s *audioStream) Read(buf []byte) (int, error) {
	s.engine.mu.RLock()
	defer s.engine.mu.RUnlock()

	if !s.engine.running {
		for i := range buf {
			buf[i] = 0
		}
		return len(buf), nil
	}

	anySolo := false
	for _, ch := range s.engine.channels {
		if ch.Solo {
			anySolo = true
			break
		}
	}

	samples := len(buf) / 4
	for i := 0; i < samples; i++ {
		var leftSum, rightSum float64

		for chIdx, ch := range s.engine.channels {
			if ch.Mute {
				continue
			}
			if anySolo && !ch.Solo {
				continue
			}

			sample := math.Sin(s.engine.phase[chIdx]) * ch.Volume
			angle := (ch.Pan + 1) * math.Pi / 4
			leftGain := math.Cos(angle)
			rightGain := math.Sin(angle)

			leftSum += sample * leftGain
			rightSum += sample * rightGain

			s.engine.phase[chIdx] += 2 * math.Pi * ch.Frequency / sampleRate
			if s.engine.phase[chIdx] > 2*math.Pi {
				s.engine.phase[chIdx] -= 2 * math.Pi
			}
		}

		leftSum *= s.engine.master
		rightSum *= s.engine.master

		leftSum = softClip(leftSum)
		rightSum = softClip(rightSum)

		leftInt := int16(leftSum * 32767 * 0.25)
		rightInt := int16(rightSum * 32767 * 0.25)

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
	return x - (x*x*x)/3
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

func (e *Engine) Close() {
	e.mu.Lock()
	e.running = false
	e.mu.Unlock()
	if e.player != nil {
		e.player.Close()
	}
}
