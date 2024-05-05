package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const MAX = 10

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type cell struct {
	value string
	fg    lipgloss.Color
	bg    lipgloss.Color
}

func (c cell) String() string {
	return c.value
}

type model struct {
	matrix    [][]cell
	textInput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your name..."
	ti.Focus()
	ti.CharLimit = MAX
	ti.Width = 20

	matrix := make([][]cell, MAX)
	for i := range matrix {
		matrix[i] = make([]cell, MAX)
	}

	return model{
		matrix:    matrix,
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	default:
		m.populate()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n",
		m.textInput.View(),
		m.render(),
		"(esc to quit)",
	) + "\n"
}

// populate fills the matrix in diagonal slices from the top right
// https://stackoverflow.com/questions/1779199/traverse-matrix-in-diagonal-strips
func (m model) populate() {
	// for each diagonal slice
	for slice := 0; slice < 2*MAX-1; slice++ {

		// convert the full input to runes so we can handle unicode
		input := []rune(m.textInput.Value())

		// reset current index of the input string for this slice
		// reset the current colors
		var i, fgColor, bgColor int

		// something to do with the slice, honestly not really sure...
		var z int = 0
		if slice >= MAX {
			z = slice - MAX + 1
		}

		// for each cell in this slice
		for j := z; j <= slice-z; j++ {

			// get the next char or " " if we're out of input
			c := " "
			if i < len(input) {
				fgColor = runeToColor(fgColor, input[i])

				// bgColor is always offset from fgColor
				bgColor = runeToColor(fgColor, input[i])
				c = string(input[i])
				i++
			}
			m.matrix[j][slice-j] = cell{
				value: c,
				fg:    lipgloss.Color(fmt.Sprintf("#%x", fgColor)),
				bg:    lipgloss.Color(fmt.Sprintf("#%x", bgColor)),
			}
		}
	}
}

func (m model) render() string {
	// convert the matrix to [][]string
	// man it would be nice if Go had a mapper func
	rows := make([][]string, MAX)
	for i, row := range m.matrix {
		rows[i] = make([]string, MAX)
		for j, c := range row {
			rows[i][j] = c.value
		}
	}

	re := lipgloss.NewRenderer(os.Stdout)
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(rows...).
		StyleFunc(
			func(row, col int) lipgloss.Style {
				return re.NewStyle().
					Align(lipgloss.Center).
					Padding(0, 1).
					Foreground(m.matrix[row-1][col].fg).
					Background(m.matrix[row-1][col].bg)
			})
	return t.Render()
}

// Add the int value of the next rune to the previous color value
// if its outside the trucolor range then start over
func runeToColor(prev int, next rune) int {
	maxColor := 16777216 // the max number of possible truecolors
	return (prev + int(next)) % maxColor
}
