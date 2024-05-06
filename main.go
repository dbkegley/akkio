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

// height/width of the identicon
const MAX = 15

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

// Each cell value in the matrix holds a single letter
// and a background/foreground color.
// This way each cell can be styled independently.
type cell struct {
	value string
	fg    lipgloss.Color
	bg    lipgloss.Color
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

	// init the 2d matrix to hold the identicon
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
		// populate the matrix everytime the input changes
		m.populate()

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

// Add the int value of the next rune to the previous color value.
// If it's outside the truecolor range then start over from the beginning.
func runeToColor(prev int, next rune) int {
	// the max number of possible truecolors
	maxColor := 16777216
	return (prev + int(next)) % maxColor
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
		var i int
		var fg, bg int = 32, 64

		// z is something to do with the slice
		// honestly not really sure... ¯\_(ツ)_/¯
		var z int = 0
		if slice >= MAX {
			z = slice - MAX + 1
		}

		// for each cell in this slice
		for j := z; j <= slice-z; j++ {

			// get the next char of the string
			// or " " if we're out of input
			c := " "
			if i < len(input) {
				fg = runeToColor(fg, input[i])

				// bg is intentionally offset again
				// from bg so that they're not he same values
				bg = runeToColor(fg, input[i])

				// convert our rune back to a string for display.
				c = string(input[i])
				i++
			}
			// assign the cell the current letter and some colors
			m.matrix[j][slice-j] = cell{
				value: c,
				fg:    lipgloss.Color(fmt.Sprintf("#%x", fg)),
				bg:    lipgloss.Color(fmt.Sprintf("#%x", bg)),
			}
		}
	}
}

func (m model) render() string {
	// convert the matrix to back to a [][]string
	rows := make([][]string, MAX)
	for i, row := range m.matrix {
		rows[i] = make([]string, MAX)
		for j, c := range row {
			rows[i][j] = c.value
		}
	}

	// render our matrix as an identicon w/ styling
	re := lipgloss.NewRenderer(os.Stdout)
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderRow(true).
		//BorderColumn(true).
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
