package list

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type item string

type defaultItem struct {
	title string
	desc  string
}

type testList struct {
	*Model
}

var appStyle = lipgloss.NewStyle().Padding(1, 2)

func (m testList) Init() tea.Cmd { return nil }

func (m testList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//switch msg := msg.(type) {
	//case tea.WindowSizeMsg:
	//m.SetSize(msg.Width, msg.Height)
	//m.SetHeight(msg.Height)
	//case tea.KeyMsg:
	//if msg.String() == "enter" {
	//return m, tea.Quit
	//}
	//}
	l, cmd := m.Model.Update(msg)
	m.Model = &l
	return m, cmd
}

func (i defaultItem) Title() string {
	return i.title
}

func (i defaultItem) Description() string {
	return ""
}

func (i defaultItem) FilterValue() string { return i.title }

func newDefaultItem(s string) defaultItem {
	return defaultItem{
		title: s,
		desc:  "",
	}
}

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                          { return 1 }
func (d itemDelegate) Spacing() int                         { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m Model, index int, listItem Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)
	fmt.Fprint(w, m.Styles.TitleBar.Render(str))
}

func TestList(t *testing.T) {
	var tItems []Item
	for _, c := range itemSlice[:10] {
		tItems = append(tItems, newDefaultItem(c))
	}

	w := 100
	h := 20
	//w, h := TermSize()
	//println(w)
	//println(h)

	del := NewDefaultDelegate()
	del.SetListType(Ul)
	//del.ShowDescription = false
	l := New(tItems, del, w, h)
	l.SetLimit(3)
	l.SetShowStatusBar(false)
	//l.SetNoLimit()
	l.SetShowTitle(true)
	m := testList{
		Model: &l,
	}
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		t.Error(err)
	}

	togItems := m.ToggledItems()
	if len(togItems) != len(m.toggledItems) {
		t.Errorf("toggled items %v != toggled %v\n", m.toggledItems, togItems)
	}
	if m.Selectable() {
		if m.Limit() > 0 && len(togItems) > m.Limit() {
			t.Errorf("toggled items %v > limit %v\n", len(togItems), m.Limit())
		}
	}
	for _, item := range togItems {
		println(item)
		fmt.Printf("%v\n", tItems[item])
	}
}

func TermSize() (int, int) {
	w, h, _ := term.GetSize(int(os.Stdin.Fd()))
	return w, h
}

func TestStatusBarItemName(t *testing.T) {
	list := New([]Item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
	expected := "2 items"
	if !strings.Contains(list.statusView(), expected) {
		t.Fatalf("Error: expected view to contain %s", expected)
	}

	list.SetItems([]Item{item("foo")})
	expected = "1 item"
	if !strings.Contains(list.statusView(), expected) {
		t.Fatalf("Error: expected view to contain %s", expected)
	}
}

func TestStatusBarWithoutItems(t *testing.T) {
	list := New([]Item{}, itemDelegate{}, 10, 10)

	expected := "No items"
	if !strings.Contains(list.statusView(), expected) {
		t.Fatalf("Error: expected view to contain %s", expected)
	}
}

func TestCustomStatusBarItemName(t *testing.T) {
	list := New([]Item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
	list.SetStatusBarItemName("connection", "connections")

	expected := "2 connections"
	if !strings.Contains(list.statusView(), expected) {
		t.Fatalf("Error: expected view to contain %s", expected)
	}

	list.SetItems([]Item{item("foo")})
	expected = "1 connection"
	if !strings.Contains(list.statusView(), expected) {
		t.Fatalf("Error: expected view to contain %s", expected)
	}

	list.SetItems([]Item{})
	expected = "No connections"
	if !strings.Contains(list.statusView(), expected) {
		t.Fatalf("Error: expected view to contain %s", expected)
	}
}

var itemSlice = []string{
	"Artichoke",
	"Baking Flour",
	"Bananas",
	"Barley",
	"Bean Sprouts",
	"Bitter Melon",
	"Blood Orange",
	"Brown Sugar",
	"Cashew Apple",
	"Cashews",
	"Cat Food",
	"Coconut Milk",
	"Cucumber",
	"Curry Paste",
	"Currywurst",
	"Dill",
	"Dragonfruit",
	"Dried Shrimp",
	"Eggs",
	"Fish Cake",
	"Furikake",
	"Garlic",
	"Gherkin",
	"Ginger",
	"Granulated Sugar",
	"Grapefruit",
	"Green Onion",
	"Hazelnuts",
	"Heavy whipping cream",
	"Honey Dew",
	"Horseradish",
	"Jicama",
	"Kohlrabi",
	"Leeks",
	"Lentils",
	"Licorice Root",
	"Meyer Lemons",
	"Milk",
	"Molasses",
	"Muesli",
	"Nectarine",
	"Niagamo Root",
	"Nopal",
	"Nutella",
	"Oat Milk",
	"Oatmeal",
	"Olives",
	"Papaya",
	"Party Gherkin",
	"Peppers",
	"Persian Lemons",
	"Pickle",
	"Pineapple",
	"Plantains",
	"Pocky",
	"Powdered Sugar",
	"Quince",
	"Radish",
	"Ramps",
	"Star Anise",
	"Sweet Potato",
	"Tamarind",
	"Unsalted Butter",
	"Watermelon",
	"Weißwurst",
	"Yams",
	"Yeast",
	"Yuzu",
}
