package main

import (
	"sort"
)

var lists [][3]byte

func isIn(list [3]byte, character byte) bool {
	return (list[0] == character) || (list[1] == character) || (list[2] == character)
}

func isValid(word string) bool {
	if len(word) < 2 {
		return false
	}
	itemNum := -1
	for i, val := range lists {
		if isIn(val, word[0]) {
			itemNum = i
		}
	}
	if itemNum == -1 {
		return false
	}
	for _, char := range []byte(word[1:]) {
		oldItemNum := itemNum
		itemNum = -1
		for i, val := range lists {
			if isIn(val, char) {
				itemNum = i
			}
		}
		if itemNum == -1 {
			return false
		}
		if itemNum == oldItemNum {
			return false
		}
	}
	return true
}

func (c *LetterBoxed) solveLetterBoxed() {
	lists = [][3]byte{[3]byte{c.List1[0], c.List1[1], c.List1[2]}, [3]byte{c.List2[0], c.List2[1], c.List2[2]}, [3]byte{c.List3[0], c.List3[1], c.List3[2]}, [3]byte{c.List4[0], c.List4[1], c.List4[2]}}

	words := c.Words
	valid := make([]string, len(words))
	count := 0

	for _, val := range words {
		if isValid(val) {
			valid[count] = val
			count++
		}
	}

	valid = valid[:count]

	categories := make(map[byte][]string, 0)
	for _, val := range valid {
		_, exists := categories[val[0]]
		if exists {
			categories[val[0]] = append(categories[val[0]], val)
		} else {
			categories[val[0]] = []string{val}
		}
	}

	sort.Slice(valid, func(i, j int) bool { return len(valid[i]) > len(valid[j]) })
	for key := range categories {
		sort.Slice(categories[key], func(i, j int) bool { return len(categories[key][i]) > len(categories[key][j]) })
	}

	combos := make([][]string, 0)
	for _, val := range valid {
		combo := []string{val}
		end := val[len(val)-1]
		_, exists := categories[end]
		if !exists {
			continue
		}
		for i := 0; i < 5; i++ {
			index := 0
			_, contains := categories[categories[end][index][len(categories[end][index])-1]]
			success := true
			for !contains {
				index++
				if index > len(categories[end]) {
					success = false
					break
				}
				_, contains = categories[categories[end][index][len(categories[end][index])-1]]
			}
			if success {
				combo = append(combo, categories[end][index])
				end = categories[end][index][len(categories[end][index])-1]
			}
		}
		cleaned := make(map[string]bool, 0)
		for _, val := range combo {
			cleaned[val] = true
		}
		finalCombo := make([]string, 0)
		for key := range cleaned {
			finalCombo = append(finalCombo, key)
		}
		combos = append(combos, finalCombo)
	}

	scores := make([]Scored, 0)
	for _, val := range combos {
		set := make(map[byte]bool, 0)
		for _, word := range val {
			for _, char := range []byte(word) {
				set[char] = true
			}
		}
		score := Scored{
			Val:   val,
			Score: 0,
		}
		for range set {
			score.Score++
		}
		scores = append(scores, score)
	}
	sort.Slice(scores, func(i, j int) bool { return scores[i].Score > scores[j].Score })

	for _, val := range scores {
		if val.Score >= 12 {
			c.Output = append(c.Output, val.Val)
		}
	}
}

// Scored has the data on a value and what it scored
type Scored struct {
	Val   []string
	Score int
}
