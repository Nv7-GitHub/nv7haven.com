package main

import (
	"strings"
)

var board [][]rune
var words []string

// offsets are all the neighbors of a given area
var offsets = [][]int{
	[]int{-1, -1},
	[]int{0, -1},
	[]int{1, -1},
	[]int{1, 0},
	[]int{1, 1},
	[]int{0, 1},
	[]int{-1, 1},
	[]int{-1, 0},
}

// parseText parses multiline text into a board, a 2-dimensional slice of runes
func parseText(text string) [][]rune {
	lines := strings.Split(text, "\n")   // Get each line of text, gets the lines and the height of the board
	output := make([][]rune, len(lines)) // Initialize a board with the correct amount of rows
	// Populate each row with a row the right size
	for i := range output {
		output[i] = make([]rune, len(lines[0]))
	}
	// Put the data into the board
	for x, val := range lines {
		for y, char := range val {
			output[x][y] = char
		}
	}
	return output
}

// findWord is the actual word-searching algorithm, returns a slice of (row, col) coordinates
func findWord(board [][]rune, word string) [][]int {
	// Create the slice of (row, col) coordinates
	output := make([][]int, len(word))
	for i := range output {
		output[i] = make([]int, 2)
	}

	// Initialize variables
	var xoff int
	var yoff int

	for x, val := range board { // Loop through rows
		for y, char := range val { // Loop through values in row
			if char == rune(word[0]) { // If letter matches,
				for _, offset := range offsets { // Loop through neighbors
					for i, lett := range word[1:] { // Loop through characters in the word after the first letter
						// Calculate the offset off of the word it would be based on the direction you are checking
						xoff = offset[0] * (i + 1)
						yoff = offset[1] * (i + 1)
						// Check if out of the bounds of the board
						if (x+xoff >= 0) && (x+xoff < len(board)) && (y+yoff >= 0) && (y+yoff < len(val)) {
							// Check if it is correct, if not continue to the next offset
							if board[x+xoff][y+yoff] == lett {
								// If this character is correct, update the output
								output[i+1] = []int{x + xoff, y + yoff}
								// If this is the last character of the word, you found the word and can finish the calculation
								if i == len(word)-2 {
									output[0] = []int{x, y} // Save first character because that is not included in loop
									return output
								}
							} else {
								break
							}
						} // Check bounds
					} // Loop through word
				} // Offsets
			} // Checking first char
		} // Y Loop
	} // X Loop

	// Nothing found, return nil to show that nothing found
	return nil
}

// main is the main function
/*func main() {
	// Multuline text, user should be inputting this
	text := `xymbitackkj
kxfyofasokd
sxtmttcbmwh
rranatbspkh
htotcdwlzor
oeptvvzopif
rkemcsgjndn
ssecqacizlw
eahcfxrbkyk
cbslmqxtrla`
	// The list of words to find
	words = []string{"cow", "tomato", "tractor", "horse", "sheep", "basket", "cat"}

	// Parse board and print original board
	board = parseText(text)
	printBoard(board)

	// Loop through words, find word, print board with the letters that were found missing
	for _, word := range words {
		fmt.Println("")
		fmt.Println(word)
		fmt.Println("")
		wordposes := findWord(board, word)
		if wordposes != nil {
			dupl := duplBoard(board)
			for _, pos := range wordposes {
				dupl[pos[0]][pos[1]] = rune(" "[0])
			}
			printBoard(dupl)
		} else {
			fmt.Println("Couldn't solve")
		}
	}
}*/
