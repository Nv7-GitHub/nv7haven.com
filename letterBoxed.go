package main

var sides [][]byte

// getSideNum gets the index of the side the letter is on (used in isValid)
func getSideNum(char byte) int {
	for i := 0; i < 4; i++ {
		for _, sideChar := range sides[i] {
			if char == sideChar {
				return i
			}
		}
	}
	return -1
}

// isValid checks if a word can be entered into the puzzle
func isValid(word string) bool {
	sideNum := 0
	for _, c := range []byte(word) {
		newSideNum := getSideNum(c)
		if newSideNum == -1 { // check that its a valid letter
			return false
		}
		if sideNum == newSideNum { // check that its different than the side before
			return false
		}
		sideNum = newSideNum
	}
	return true
}

var words []string
var letterMap map[byte][]int

func iddfs(v []int, maxLen int) []int {
	if len(v) == maxLen {
		return nil
	}

	// Check if it is valid
	needed := [256]bool{}        // This is basically map[byte]bool where needed[char] is true if the letter is on the board
	for _, side := range sides { // Set all spots where a letter is on the board to true
		for _, c := range side {
			needed[int(c)] = true
		}
	}
	for _, ind := range v { // Set all spots where a letter is in the guess to false
		for _, c := range []byte(words[ind]) {
			needed[int(c)] = false
		}
	}
	// Check if any spots are still true
	done := true
	for _, v := range needed {
		if v {
			done = false
			break
		}
	}
	if done { // If no spots are true, then we have a solution
		return v
	}

	// Go through all the words that start with the last letter we have and test if there is a solution with them
	last := words[v[len(v)-1]]
	for _, next := range letterMap[last[len(last)-1]] {
		if next == v[len(v)-1] {
			continue
		}
		v = append(v, next) // Add the word to the guess
		res := iddfs(v, maxLen)
		if res != nil {
			return res // Return the guess, it worked!
		} else {
			v = v[:len(v)-1] // Remove since it didn't work
		}
	}

	// None of the words starting with the last letter worked, so there aren't any solutions for this guess
	return nil
}

func (c *LetterBoxed) solveLetterBoxed() {
	// Sides
	sides = [][]byte{
		[]byte(c.List1),
		[]byte(c.List2),
		[]byte(c.List3),
		[]byte(c.List4),
	}

	// Build letterMap (map of first letter of word to all words with that first letter)
	old := words
	words = make([]string, 0)
	for _, word := range old {
		if isValid(word) {
			words = append(words, word)
		}
	}

	letterMap = make(map[byte][]int)
	for i, word := range words {
		_, exists := letterMap[word[0]]
		if exists {
			letterMap[word[0]] = append(letterMap[word[0]], i)
		} else {
			letterMap[word[0]] = []int{i}
		}
	}

	// Solve using iterative deepening depth-first search
	max := 0
	for {
		max++
		for i := range words {
			res := iddfs([]int{i}, max)
			if res != nil {
				c.Output = make([]string, len(res))
				for i, ind := range res {
					c.Output[i] = words[ind]
				}
				return
			}
		}
	}
}
