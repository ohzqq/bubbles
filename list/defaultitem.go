package list

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/ohzqq/bubbles/key"
)

// ListType describes if the list is ordered or unordered.
type ListType int

const (
	Ul ListType = iota // unordered list
	Ol                 // ordered list
)

// String returns a human-readable string of the list type.
func (f ListType) String() string {
	return [...]string{
		"unordered list",
		"ordered list",
	}[f]
}

// DefaultItemStyles defines styling for a default list item.
// See DefaultItemView for when these come into play.
type DefaultItemStyles struct {
	// The Normal state.
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	// The selected item state.
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style

	Prefix lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 1)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#AFFF87", Dark: "#AFFFAF"}).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#AF87FF", Dark: "#AFAFFF"})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 1)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle(). //Underline(true)
						Foreground(lipgloss.AdaptiveColor{Light: "#5FAFFF", Dark: "#AFFFFF"})

	s.Prefix = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FF87FF", Dark: "#FFAFFF"})

	return s
}

// DefaultItem describes an items designed to work with DefaultDelegate.
type DefaultItem interface {
	Item
	Title() string
	Description() string
}

// DefaultDelegate is a standard delegate designed to work in lists. It's
// styled by DefaultItemStyles, which can be customized as you like.
//
// The description line can be hidden by setting Description to false, which
// renders the list as single-line-items. The spacing between items can be set
// with the SetSpacing method.
//
// Setting UpdateFunc is optional. If it's set it will be called when the
// ItemDelegate called, which is called when the list's Update function is
// invoked.
//
// Settings ShortHelpFunc and FullHelpFunc is optional. They can be set to
// include items in the list's default short and full help menus.
type DefaultDelegate struct {
	ShowDescription bool
	Styles          DefaultItemStyles
	UpdateFunc      func(tea.Msg, *Model) tea.Cmd
	ShortHelpFunc   func() []key.Binding
	FullHelpFunc    func() [][]key.Binding
	listType        ListType
	height          int
	spacing         int
}

// NewDefaultDelegate creates a new delegate with default styles.
func NewDefaultDelegate() DefaultDelegate {
	return DefaultDelegate{
		Styles:   NewDefaultItemStyles(),
		height:   2,
		spacing:  1,
		listType: Ul,
	}
}

// SetHeight sets delegate's preferred height.
func (d *DefaultDelegate) SetHeight(i int) {
	d.height = i
}

// Height returns the delegate's preferred height.
// This has effect only if ShowDescription is true,
// otherwise height is always 1.
func (d DefaultDelegate) Height() int {
	if d.ShowDescription {
		return d.height
	}
	return 1
}

// SetSpacing sets the delegate's spacing.
func (d *DefaultDelegate) SetSpacing(i int) {
	d.spacing = i
}

// Spacing returns the delegate's spacing.
func (d DefaultDelegate) Spacing() int {
	return d.spacing
}

// SetListType sets the list type: ordered or unordered.
func (d *DefaultDelegate) SetListType(lt ListType) {
	d.listType = lt
}

// ListType returns the list type.
func (d *DefaultDelegate) ListType() ListType {
	return d.listType
}

// Update checks whether the delegate's UpdateFunc is set and calls it.
func (d DefaultDelegate) Update(msg tea.Msg, m *Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

// Render prints an item.
func (d DefaultDelegate) Render(w io.Writer, m Model, index int, item Item) {
	var (
		title, desc  string
		matchedRunes []int
		s            = &d.Styles
	)

	if i, ok := item.(DefaultItem); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	if m.width <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := uint(m.width - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight() - s.Prefix.GetPaddingLeft())
	title = truncate.StringWithTail(title, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.height-1 {
				break
			}
			lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
		}
		desc = strings.Join(lines, "\n")
	}

	// Conditions
	var (
		prefix      string
		padding     = len(strconv.Itoa(len(m.Items())))
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == Filtering || m.FilterState() == FilterApplied
		isToggled   = m.ItemIsToggled(item)
	)

	// style prefix
	switch d.listType {
	case Ol:
		p := "%" + strconv.Itoa(padding) + "d."
		prefix = fmt.Sprintf(p, index+1)
	default:
		prefix = " "
	}

	if m.MultiSelectable() {
		if isToggled {
			prefix = s.Prefix.Render(m.toggledPrefix + prefix)
		} else {
			prefix = s.Prefix.Render(m.untoggledPrefix + prefix)
		}
	}

	if isSelected {
		if d.listType == Ul && !m.MultiSelectable() {
			prefix = m.prefix
		}
		prefix = s.Prefix.Render(prefix)
	}

	// style filtered items
	if isFiltered && index < len(m.filteredItems) {
		// Get indices of matched characters
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		desc = s.DimmedDesc.Render(desc)
	} else if isSelected && m.FilterState() != Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := s.FilterMatch
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
		desc = s.SelectedDesc.Render(desc)
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := s.FilterMatch
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s%s\n%s", prefix, title, desc)
		return
	}
	fmt.Fprintf(w, "%s%s", prefix, title)
}

// ShortHelp returns the delegate's short help.
func (d DefaultDelegate) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

// FullHelp returns the delegate's full help.
func (d DefaultDelegate) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
