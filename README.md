# **Conway's Game of Life**

<!--toc:start-->

- [**Conway's Game of Life**](#conways-game-of-life)
  - [**About**](#about)
  - [**Features**](#features)
  - [**Usage**](#usage)
  - [**Requirements**](#requirements)
  - [**Contributing**](#contributing)
  - [**License**](#license)
  <!--toc:end-->

A console-based implementation of Conway's Game of Life, written in Go.

## **About**

The Game of Life is a famous example of a cellular automaton, created by John
Horton Conway in 1970. The game is a grid of cells that follow simple rules to
evolve over time. This project implements the Game of Life in a console-based
environment using the Go programming language.

## **Features**

- Supports mouse and keyboard input for cell editing
- Allows for pausing and resuming the simulation
- Displays the current generation number
- Uses a color scheme based on the x and y coordinates to visually represent
  the cells

## **Usage**

- Run the program by executing `main.go` using Go's built-in compiler.
- The game will start in paused mode, allowing you to explore the initial
  state of the grid.
- Use the `p` key to pause or resume the simulation.
- Use the `w`, `a`, `s`, and `d` keys to move the selection cursor.
- Use the `t` key or `space` to toggle the state of the cell at the current
  cursor position.
- Use the `?` key to toggle the help overlay.

> [!NOTE]
> Also supports the arrow keys and vim-like movement keys.

- Use the `t` key to toggle the cell at the current cursor position.
- Use the mouse to select a cell and toggle its state.

## **Requirements**

- Go installed on your system

## **Contributing**

If you'd like to contribute to this project, please feel free to open an issue
or submit a pull request with your proposed changes.
All contributions are appreciated!

## **License**

This project is licensed under the MIT License. See [`LICENSE`](LICENSE) for
more details.
