package main

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
	cells        [][]cell
	selected     cell
	height       int
	width        int
	paused       bool
	overlayShown bool
	generation   int
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
		case "?":
			m.overlayShown = !m.overlayShown
			return m, nil
		case "p":
			m.paused = !m.paused
			if !m.paused {
				return m, tickCmd()
			}
			return m, nil
		}

		if m.paused {
			switch msg.String() {
			case "w", "up", "k":
				if m.selected.row > 0 {
					m.selected.row -= 1
				}
			case "s", "down", "j":
				if m.selected.row < m.height-1 {
					m.selected.row += 1
				}
			case "a", "left", "h":
				if m.selected.col > 0 {
					m.selected.col -= 1
				}
			case "d", "right", "l":
				if m.selected.col < m.width-1 {
					m.selected.col += 1
				}
			case "t", " ":
				m.cells[m.selected.row][m.selected.col].alive = !m.cells[m.selected.row][m.selected.col].alive
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width / 2
		m.cells = [][]cell{}
		m.overlayShown = true
		m.paused = true
		m.generation = 0
		for row := range m.height {
			cellsRow := []cell{}
			for col := range m.width {
				cellsRow = append(cellsRow, cell{row: row, col: col, alive: false, id: fmt.Sprintf("%4d:%4d", row+1, col+1)})
			}
			m.cells = append(m.cells, cellsRow)
		}
		m.selected = cell{col: m.width / 2, row: m.height / 2}
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
		if m.paused {
			return m, nil
		}

		// GOL
		m.nextState()
		m.generation++
		return m, tickCmd()
	}

	return m, nil
}

func GenerateColor(maxX, maxY, x, y int) string {
	if x < 0 || x > maxX || y < 0 || y > maxY {
		return "#FFFFFF"
	}

	// Normalize x and y to a value between 0 and 255
	red := (x * 255) / maxX
	green := (y * 255) / maxY
	blue := ((maxX - x) * (maxY - y) * 255) / (maxX * maxY)

	// Format the color as a hex string
	color := fmt.Sprintf("#%02X%02X%02X", red, green, blue)

	return color
}

func (m model) GetOverlayChars(row, col int) (string, string) {
	cellColor := GenerateColor(m.width, m.height, col, row)

	if m.paused && m.selected.row == row && m.selected.col == col {
		return "⿻", cellColor
	}

	empty := "  "

	// check if we can send the overlay to the screen
	if !m.overlayShown {
		return "  ", cellColor
	}

	overlay := []string{
		empty,
		"Game Of Life",
		"Generation: " + fmt.Sprintf("%d", m.generation),
		"Paused: " + fmt.Sprintf("%t", m.paused),
		empty,
		"Keybindings:",
		"? to toggle help",
		"p to pause/play",
		"w,a,s,d to move",
		"t to toggle cell",
	}

	if row < len(overlay) {
		curr := empty + overlay[row]

		// get the relevant two characters
		if col < ((len(curr) + 1) / 2) {
			to := col*2 + 2
			suffix := ""
			if to > len(curr) {
				to = col*2 + 1
				suffix = " "
			}
			return curr[col*2:to] + suffix, cellColor
		}
	}

	return "  ", cellColor
}

func (m model) View() string {
	if !m.isInitialized() {
		return ""
	}

	active := lipgloss.NewStyle()

	finalRender := []string{}

	for row := range m.height {
		cellsRow := []string{}
		for col := range m.width {
			cellContent, cellColor := m.GetOverlayChars(row, col)
			if m.cells[row][col].alive {
				cellsRow = append(
					cellsRow,
					active.
						Background(
							lipgloss.Color(cellColor),
						).
						Render(cellContent),
				)
			} else {
				cellsRow = append(cellsRow, cellContent)
			}
		}
		finalRender = append(finalRender, lipgloss.JoinHorizontal(lipgloss.Top, cellsRow...))
	}

	return lipgloss.JoinVertical(lipgloss.Top, finalRender...)
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
