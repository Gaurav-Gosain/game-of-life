package main

/*
* The RULES
*
* - Any live cell with fewer than two live neighbors dies, as if by underpopulation.
* - Any live cell with two or three live neighbors lives on to the next generation.
* - Any live cell with more than three live neighbors dies, as if by overpopulation.
* - Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
*
* */

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cell struct {
	id    string
	row   int
	col   int
	alive bool
}

type tickMsg time.Time

type model struct {
	cells  [][]cell
	height int
	width  int
}

// Function to get the number of alive neighbors for a given cell
func aliveNeighbors(m *model, row int, col int) int {
	directions := []struct{ rowOffset, colOffset int }{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	count := 0
	for _, direction := range directions {
		newRow := row + direction.rowOffset
		newCol := col + direction.colOffset

		if newRow >= 0 && newRow < m.height && newCol >= 0 && newCol < m.width && m.cells[newRow][newCol].alive {
			count++
		}
	}
	return count
}

// Function to calculate the next state of the grid
func (m *model) nextState() {
	// Create a new grid to store the next state
	newCells := make([][]cell, m.height)
	for i := range newCells {
		newCells[i] = make([]cell, m.width)
	}

	for row := 0; row < m.height; row++ {
		for col := 0; col < m.width; col++ {
			newCells[row][col] = m.cells[row][col] // Copy current cell to new state

			aliveNeighborsCount := aliveNeighbors(m, row, col)
			if m.cells[row][col].alive {
				// Rule 1: Any live cell with fewer than two live neighbors dies (underpopulation).
				// Rule 2: Any live cell with two or three live neighbors lives on to the next generation.
				// Rule 3: Any live cell with more than three live neighbors dies (overpopulation).
				if aliveNeighborsCount < 2 || aliveNeighborsCount > 3 {
					newCells[row][col].alive = false
				}
			} else {
				// Rule 4: Any dead cell with exactly three live neighbors becomes a live cell (reproduction).
				if aliveNeighborsCount == 3 {
					newCells[row][col].alive = true
				}
			}
		}
	}

	m.cells = newCells
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isInitialized() {
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width / 2
		m.cells = [][]cell{}
		for row := range m.height {
			cellsRow := []cell{}
			for col := range m.width {
				cellsRow = append(cellsRow, cell{row: row, col: col, alive: false, id: fmt.Sprintf("%4d:%4d", row+1, col+1)})
			}
			m.cells = append(m.cells, cellsRow)
		}
	case tea.MouseMsg:
		mouseButton := tea.MouseEvent(msg).Button
		isLeftClick := mouseButton == tea.MouseButtonLeft
		isRightClick := mouseButton == tea.MouseButtonRight

		if !isLeftClick && !isRightClick {
			return m, nil
		}

		x, y := msg.X, msg.Y

		col := x / 2
		row := y

		if col < 0 || col > m.width-1 || row < 0 || row > m.height-1 {
			return m, nil
		}

		m.cells[row][col].alive = isLeftClick

		return m, nil
	case tickMsg:

		// GOL
		m.nextState()
		return m, tickCmd()
	}

	return m, nil
}

func (m model) View() string {
	if !m.isInitialized() {
		return ""
	}

	active := lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF"))

	str := ""

	for row := range m.height {
		str += "\n"
		for col := range m.width {
			if m.cells[row][col].alive {
				str += active.Render("  ")
			} else {
				str += "  "
			}
		}
	}

	return str
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	m := &model{}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
