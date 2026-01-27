package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"pomelo/styles"
)

type Header struct {
	logo    string
	desc    string
	version string
	width   int
}

func NewHeader(logo, desc, version string, width int) Header {
	return Header{
		logo:    logo,
		desc:    desc,
		version: version,
		width:   width,
	}
}

func (h *Header) SetWidth(width int) {
	h.width = width
}

func (h Header) View() string {
	meta := lipgloss.JoinVertical(lipgloss.Right, h.desc, h.version)
	metaWidth := lipgloss.Width(meta)
	logoWidth := lipgloss.Width(h.logo)
	spaceWidth := max(h.width-metaWidth-logoWidth-6, 1)
	space := lipgloss.JoinVertical(lipgloss.Left, strings.Repeat(" ", spaceWidth), strings.Repeat(" ", spaceWidth))

	return styles.Header.Render(lipgloss.JoinHorizontal(lipgloss.Top, h.logo, space, meta))
}
