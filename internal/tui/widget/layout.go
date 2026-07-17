package widget

import "github.com/gdamore/tcell/v2"

// VerticalLayout stacks primitives vertically. Each item has a fixed height,
// except that at most one item may have height 0 and takes the remaining
// space. An overlay primitive, when set, is drawn last, on top.
type VerticalLayout struct {
	*Box
	items   []layoutItem
	overlay Primitive
}

type layoutItem struct {
	primitive Primitive
	height    int // 0 = take the remaining space
}

func NewVerticalLayout() *VerticalLayout {
	return &VerticalLayout{Box: NewBox()}
}

func (v *VerticalLayout) AddItem(p Primitive, height int) *VerticalLayout {
	v.items = append(v.items, layoutItem{primitive: p, height: height})
	return v
}

// ShowOverlay draws p on top of the layout until RemoveOverlay is called.
// The overlay is responsible for its own position and size.
func (v *VerticalLayout) ShowOverlay(p Primitive) {
	v.overlay = p
}

func (v *VerticalLayout) RemoveOverlay() {
	v.overlay = nil
}

func (v *VerticalLayout) SetRect(x, y, width, height int) {
	v.Box.SetRect(x, y, width, height)
	fixed := 0
	for _, item := range v.items {
		fixed += item.height
	}
	remaining := max(height-fixed, 0)
	row := y
	for _, item := range v.items {
		itemHeight := item.height
		if itemHeight == 0 {
			itemHeight = remaining
		}
		item.primitive.SetRect(x, row, width, itemHeight)
		row += itemHeight
	}
}

func (v *VerticalLayout) Draw(screen tcell.Screen) {
	v.Box.Draw(screen)
	for _, item := range v.items {
		item.primitive.Draw(screen)
	}
	if v.overlay != nil {
		v.overlay.Draw(screen)
	}
}
