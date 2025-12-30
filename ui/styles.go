package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary    = lipgloss.Color("#7C3AED") // Purple
	ColorSecondary  = lipgloss.Color("#10B981") // Green
	ColorAccent     = lipgloss.Color("#F59E0B") // Amber
	ColorMuted      = lipgloss.Color("#EF4444") // Red
	ColorSolo       = lipgloss.Color("#3B82F6") // Blue
	ColorBackground = lipgloss.Color("#1F2937") // Dark gray
	ColorSurface    = lipgloss.Color("#374151") // Medium gray
	ColorText       = lipgloss.Color("#F9FAFB") // Light gray
	ColorTextDim    = lipgloss.Color("#9CA3AF") // Dimmed text
	ColorFader      = lipgloss.Color("#4ADE80") // Bright green
	ColorFaderBg    = lipgloss.Color("#374151") // Fader background
)

// Styles
var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Background(ColorBackground).
			Foreground(ColorText)

	// Title bar
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	// Channel strip container
	ChannelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSurface).
			Padding(1).
			Width(10).
			Align(lipgloss.Center)

	// Selected channel
	SelectedChannelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(1).
				Width(10).
				Align(lipgloss.Center)

	// Channel name
	ChannelNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorText).
				Align(lipgloss.Center)

	// Fader track (background)
	FaderTrackStyle = lipgloss.NewStyle().
			Foreground(ColorFaderBg)

	// Fader fill (active part)
	FaderFillStyle = lipgloss.NewStyle().
			Foreground(ColorFader)

	// Value display
	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Align(lipgloss.Center)

	// Mute button
	MuteActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBackground).
			Background(ColorMuted).
			Padding(0, 1)

	MuteInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextDim).
				Padding(0, 1)

	// Solo button
	SoloActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBackground).
			Background(ColorSolo).
			Padding(0, 1)

	SoloInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextDim).
				Padding(0, 1)

	// Pan indicator
	PanStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)

	// Help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			MarginTop(1)

	// Status bar
	StatusStyle = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			MarginTop(1)

	// Device selector styles
	DeviceListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSurface).
			Padding(1).
			Width(50)

	DeviceItemStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Padding(0, 2)

	DeviceSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorBackground).
				Background(ColorPrimary).
				Padding(0, 2)

	// Master fader
	MasterStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorAccent).
			Padding(1).
			Width(12).
			Align(lipgloss.Center)
)
