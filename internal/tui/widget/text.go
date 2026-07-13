package widget

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"
)

// StringWidth returns the number of terminal cells needed to display s.
func StringWidth(s string) int {
	return uniseg.StringWidth(s)
}

// WordWrap splits text into lines of at most width cells. It breaks at the
// last space that fits; words longer than width are hard-broken. Explicit
// "\n" line breaks are preserved.
func WordWrap(text string, width int) []string {
	if width < 1 {
		return nil
	}
	var lines []string
	for line := range strings.SplitSeq(text, "\n") {
		for StringWidth(line) > width {
			breakAt := -1
			lastSpace := -1
			cells := 0
			pos := 0
			graphemes := uniseg.NewGraphemes(line)
			for graphemes.Next() {
				if cells+graphemes.Width() > width {
					breakAt = pos
					break
				}
				if graphemes.Str() == " " {
					lastSpace = pos
				}
				cells += graphemes.Width()
				pos += len(graphemes.Str())
			}
			if breakAt <= 0 {
				break
			}
			if lastSpace > 0 {
				lines = append(lines, line[:lastSpace])
				line = line[lastSpace+1:]
			} else {
				lines = append(lines, line[:breakAt])
				line = line[breakAt:]
			}
		}
		lines = append(lines, line)
	}
	return lines
}

// DrawText draws text at (x, y), clipped to maxWidth cells, and returns the
// printed width in cells.
func DrawText(screen tcell.Screen, x, y, maxWidth int, text string, style tcell.Style) int {
	return DrawTextStyled(screen, x, y, maxWidth, text, style, nil)
}

// DrawTextStyled draws text like DrawText. If styleAt is non-nil, it is
// called with the byte index of each grapheme cluster and may return a
// modified style for it.
func DrawTextStyled(screen tcell.Screen, x, y, maxWidth int, text string, style tcell.Style, styleAt func(byteIndex int, base tcell.Style) tcell.Style) int {
	printed := 0
	byteIndex := 0
	graphemes := uniseg.NewGraphemes(text)
	for graphemes.Next() {
		width := graphemes.Width()
		if width == 0 {
			byteIndex += len(graphemes.Str())
			continue
		}
		if printed+width > maxWidth {
			break
		}
		st := style
		if styleAt != nil {
			st = styleAt(byteIndex, style)
		}
		runes := graphemes.Runes()
		screen.SetContent(x+printed, y, runes[0], runes[1:], st)
		printed += width
		byteIndex += len(graphemes.Str())
	}
	return printed
}
