package main

import (
	"math/rand"
	"strings"
)

type change struct {
	row int
	col int
}

const letters = "abcdefghijklmnopqrstuvwxyz"

var changes = []change{
	change{row: 0, col: -1},
	change{row: 1, col: 0},
	change{row: 0, col: 1},
	change{row: -1, col: 0},
	change{row: -1, col: -1},
	change{row: 1, col: -1},
	change{row: -1, col: 1},
	change{row: -1, col: -1},
}

func genPos(brd [][]rune, word string) (dir change, x int, y int) {
	valid := false
	for !valid {
		valid = true
		dir = changes[rand.Intn(len(changes))]
		x = rand.Intn(len(brd))
		y = rand.Intn(len(brd))
		cx := x + dir.row*len(word)
		cy := y + dir.col*len(word)
		if (cx > len(brd)) || (cy > len(brd[0])) || (cx < 1) || (cy < 1) { // Section of word out of bounds
			valid = false
		}
	}
	return
}

// GenWordSearch generates a word search
func GenWordSearch(words []string, w, h int) string {
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	brd := make([][]rune, h)
	for i := 0; i < len(brd); i++ {
		brd[i] = make([]rune, w)
	}

	// Put in words
	for _, word := range words {
		var dir change
		var x, y int

		success := false
		for !success {
			success = true
			dir, x, y = genPos(brd, word)
			x--
			y--
			xtemp := x
			ytemp := y
			for _, char := range word {
				if (xtemp < 0) || (ytemp < 0) {
					success = false
					break
				}

				if (brd[xtemp][ytemp] != 0) && (brd[xtemp][ytemp] != char) {
					success = false
					break
				}
				xtemp += dir.row
				ytemp += dir.col
			}
		}

		for _, char := range word {
			brd[x][y] = char
			x += dir.row
			y += dir.col
		}
	}

	// Fill in the rest
	for x := 0; x < len(brd); x++ {
		for y := 0; y < len(brd[x]); y++ {
			if brd[x][y] == 0 {
				brd[x][y] = rune(letters[rand.Intn(len(letters))])
			}
		}
	}

	out := ""
	for _, row := range brd {
		out += string(row) + "\n"
	}

	return out
}
