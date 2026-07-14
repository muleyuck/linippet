package widget

import (
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// listItem is one row of a List.
type listItem struct {
	mainText      string
	secondaryText string // not drawn; carries caller data such as an ID
	matchIndices  []int  // byte indices in mainText to highlight
}

// List displays selectable rows of text with optional per-byte match
// highlighting. Navigation is driven externally via SetCurrentItem; the list
// itself handles no keys.
type List struct {
	*Box
	items             []*listItem
	currentItem       int
	itemOffset        int // number of items scrolled off the top
	selectedLabel     string
	mainTextStyle     tcell.Style
	selectedStyle     tcell.Style
	matchedColor      tcell.Color
	highlightFullLine bool
}

func NewList() *List {
	return &List{
		Box:           NewBox(),
		mainTextStyle: tcell.StyleDefault,
		selectedStyle: tcell.StyleDefault.Reverse(true),
		matchedColor:  tcell.ColorGreen,
	}
}

// SetLabel sets the text displayed in front of the selected item.
func (l *List) SetLabel(label string) *List {
	l.selectedLabel = label
	return l
}

func (l *List) SetMainTextStyle(style tcell.Style) *List {
	l.mainTextStyle = style
	return l
}

func (l *List) SetSelectedStyle(style tcell.Style) *List {
	l.selectedStyle = style
	return l
}

// SetHighlightFullLine makes the selected item's background span the whole
// width of the list.
func (l *List) SetHighlightFullLine(highlight bool) *List {
	l.highlightFullLine = highlight
	return l
}

// AddItem appends an item. matchIndices are byte indices into mainText whose
// grapheme clusters are drawn with the match highlight color.
func (l *List) AddItem(mainText, secondaryText string, matchIndices []int) *List {
	l.items = append(l.items, &listItem{
		mainText:      mainText,
		secondaryText: secondaryText,
		matchIndices:  matchIndices,
	})
	return l
}

func (l *List) Clear() *List {
	l.items = nil
	l.currentItem = 0
	l.itemOffset = 0
	return l
}

func (l *List) GetItemCount() int { return len(l.items) }

// GetItemText returns an item's main and secondary text. Panics if the index
// is out of range.
func (l *List) GetItemText(index int) (main, secondary string) {
	return l.items[index].mainText, l.items[index].secondaryText
}

func (l *List) GetCurrentItem() int {
	return l.currentItem
}

// SetCurrentItem sets the selected item, clamping out-of-range indices.
func (l *List) SetCurrentItem(index int) *List {
	if index >= len(l.items) {
		index = len(l.items) - 1
	}
	if index < 0 {
		index = 0
	}
	l.currentItem = index
	return l
}

func (l *List) Draw(screen tcell.Screen) {
	l.Box.Draw(screen)
	x, y, width, height := l.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}

	// Keep the selected item in view.
	if l.currentItem < l.itemOffset {
		l.itemOffset = l.currentItem
	} else if l.currentItem-l.itemOffset >= height {
		l.itemOffset = l.currentItem + 1 - height
	}

	labelWidth := StringWidth(l.selectedLabel)
	row := y
	for index := l.itemOffset; index < len(l.items) && row < y+height; index++ {
		item := l.items[index]
		selected := index == l.currentItem

		if labelWidth > 0 {
			label := strings.Repeat(" ", labelWidth)
			if selected {
				label = l.selectedLabel
			}
			DrawText(screen, x, row, width, label, l.mainTextStyle)
		}

		style := l.mainTextStyle
		if selected {
			style = l.selectedStyle
		}
		printed := DrawTextStyled(screen, x+labelWidth, row, width-labelWidth, item.mainText, style,
			func(byteIndex int, base tcell.Style) tcell.Style {
				if slices.Contains(item.matchIndices, byteIndex) {
					return base.Foreground(l.matchedColor)
				}
				return base
			})

		if selected && l.highlightFullLine {
			for cx := x + labelWidth + printed; cx < x+width; cx++ {
				screen.SetContent(cx, row, ' ', nil, style)
			}
		}
		row++
	}
}
