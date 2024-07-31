package main

import (
	"flag"
	"fmt"

	"github.com/eiannone/keyboard"
)

type move int

const (
	playerX move = iota - 1
	empty
	playerY
)

func main() {
	var size int
	flag.IntVar(&size, "size", 3, "Size of the Tic-Tac-Goe grid (n x n)")
	flag.Parse()

	board := createBoard(size) // two dimensional slice to track board moves
	cursorX, cursorY := 0, 0   // current cursor location
	counter := 0               // count of moves made

	// start keyboard listener
	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return
	}
	defer keyboard.Close()

	for {
		printBoard(board, cursorX, cursorY)

		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key press:", err)
			continue
		}

		if key == keyboard.KeyEsc {
			// exit the game
			break
		}

		switch key {
		case keyboard.KeyArrowUp:
			if cursorX > 0 {
				cursorX--
			}
		case keyboard.KeyArrowDown:
			if cursorX < size-1 {
				cursorX++
			}
		case keyboard.KeyArrowLeft:
			if cursorY > 0 {
				cursorY--
			}
		case keyboard.KeyArrowRight:
			if cursorY < size-1 {
				cursorY++
			}
		case keyboard.KeySpace:
			// set player move if cursor location is empty
			if board[cursorX][cursorY] == empty {
				m := currentMove(counter)
				board[cursorX][cursorY] = m
				counter++
			}
		default:
			// no-op
		}

		winner := findWinner(board)
		if winner != empty {
			fmt.Printf("We have a winner! %s wins!\n", moveDisplayValue(winner))
			break
		}

		if counter == size*size {
			fmt.Println("Game is tied! nobody wins!")
			break
		}
	}
}

func createBoard(n int) [][]move {
	// create rows by n
	board := make([][]move, n)

	for i := range board {
		// create row
		board[i] = make([]move, n)

		// populate empty state
		for j := range board[i] {
			board[i][j] = empty
		}
	}

	return board
}

func printBoard(board [][]move, x, y int) {
	size := len(board[0])

	// Clear the terminal and move to top left corner
	fmt.Print("\033[H\033[2J")

	// Print instructions
	fmt.Println("Arrow: move cursor, Space: make move, ESC: exit")

	// figure out row divider sizing
	div := ""
	for i := 0; i < size; i++ {
		div += "---"

		// add spacer for grid
		if i != size-1 {
			div += "-"
		}
	}

	// iterate through the board and draw lines
	for i, row := range board {
		line := ""

		for j, val := range row {
			// check for cursor location and indicate with brackets
			if x == i && y == j {
				line += fmt.Sprintf("[%s]", moveDisplayValue(val))
			} else {
				line += fmt.Sprintf(" %s ", moveDisplayValue(val))
			}

			// check if vertical divider is needed
			if j != size-1 {
				line += "|"
			}
		}

		fmt.Println(line)

		// check if horizontal divider is needed
		if i != size-1 {
			fmt.Println(div)
		}
	}

}

func currentMove(n int) move {
	if n%2 == 0 {
		return playerX
	}

	return playerY
}

func moveDisplayValue(m move) string {
	switch m {
	case playerX:
		return "x"
	case playerY:
		return "o"
	default:
		return " "
	}
}

func findWinner(board [][]move) move {
	size := len(board)

	// scan rows and columns
	for i, row := range board {
		winner := findRowWinner(row)
		if winner != empty {
			return winner
		}

		winner = findColumnWinner(board, i)
		if winner != empty {
			return winner
		}

	}

	if size%2 != 0 {
		winner := findDiagonalWinner(board)
		if winner != empty {
			return winner
		}
	}

	return empty
}

func findRowWinner(row []move) move {
	score := 0
	for _, val := range row {
		score += int(val)
	}

	return findWinnerFromScore(len(row), score)
}

func findColumnWinner(board [][]move, index int) move {
	score := 0
	for _, row := range board {
		score += int(row[index])
	}

	return findWinnerFromScore(len(board), score)
}

func findDiagonalWinner(board [][]move) move {
	size := len(board)
	leftScore, rightScore := 0, 0

	for i, row := range board {
		leftScore += int(row[i])
		rightScore += int(row[size-1])
	}

	winner := findWinnerFromScore(size, leftScore)
	if winner != empty {
		return winner
	}

	return findWinnerFromScore(size, rightScore)
}

func findWinnerFromScore(size, score int) move {
	if score < 0 && score == size*-1 {
		return playerX
	}

	if score > 0 && score == size {
		return playerY
	}

	return empty
}
