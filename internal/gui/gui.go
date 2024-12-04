package gui

import (
	"fmt"
	"strings"

	"github.com/burgr033/t00lbox/internal/file"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/skratchdot/open-golang/open"
)

var (
	docStyle               = lipgloss.NewStyle().Margin(1, 2)
	borderStyle            = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))
	levelBeginnerStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#baed91", Dark: "#baed91"})
	levelIntermediateStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#f8b88b", Dark: "#f8b88b"})
	levelExpertStyle       = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#fea3aa", Dark: "#fea3aa"})
)

// item custom item struct
type item struct {
	title, desc, markdown string
	isSoftware            bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc + "\n" + i.desc }
func (i item) Markdown() string    { return i.markdown }
func (i item) FilterValue() string { return i.title + " " + i.desc }

// model custom model struct
type model struct {
	list       list.Model
	showing    bool
	markdown   string
	isSoftware bool
}

// NewModel gerneates the list from the yaml file and manipulates the items
func NewModel(resourcesFile *file.ResourcesFile) model {
	items := []list.Item{}

	for _, v := range resourcesFile.Resources {
		categories := strings.Join(v.Categories, " #")
		categories = "#" + categories
		sep := " • "
		desc := v.Description + "\n[" + categories + "]"
		name := v.Name
		if v.Recommended {
			desc += sep + "recommended"
		}
		if v.Documentation != "" && v.IsSoftware {
			desc += sep + "docs available"
		}
		if v.IsSoftware {
			name = "🛠️ " + name
		} else {
			name = "📖 " + name
		}
		switch v.Level {
		case 1:
			desc += sep + levelBeginnerStyle.Render("beginner")
		case 2:
			desc += sep + levelIntermediateStyle.Render("intermediate")
		case 3:
			desc += sep + levelExpertStyle.Render("expert")
		default:
			desc += sep + string(v.Level)
		}
		items = append(items, item{title: name, desc: desc, markdown: v.Documentation, isSoftware: v.IsSoftware})
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(3)

	l := list.New(items, delegate, 0, 0)
	l.Title = "🧰 t00lbox"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"),
				key.WithHelp("enter", "opens documentation"),
			),
		}
	}

	return model{list: l}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Update handles the default model interactions and also adds the handling of the enter key
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if !i.isSoftware {
					err := open.Run(i.markdown)
					if err == nil {
						return m, nil
					} else {
						m.markdown = "error opening URL"
						m.showing = !m.showing
						return m, nil
					}
				}
				m.showing = !m.showing
				if m.showing {
					m.markdown = i.Markdown()
				} else {
					m.markdown = ""
				}
			}
			return m, nil
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View overriding the default bubbletea View method for enabling the markdown window
func (m model) View() string {
	if m.showing {
		out, err := glamour.Render(m.markdown, "dark")
		if err != nil {
			out = "Error rendering markdown"
		}
		return borderStyle.Render(out)
	}
	return docStyle.Render(m.list.View())
}

// CreateProgram loads the resources file, creates the model and launches the actual program
func CreateProgram(resourcesPath string) (*tea.Program, error) {
	resourcesFile, err := file.LoadResources()
	if err != nil {
		return nil, fmt.Errorf("error loading %s: %w", resourcesPath, err)
	}

	m := NewModel(resourcesFile)
	return tea.NewProgram(m, tea.WithAltScreen()), nil
}
