package main

import (
	"fmt"
)

func mymax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// lengthOfLongestSubstring возвращает длину самой длинной подстроки
// без повторяющихся символов.
func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}

	cnt := make(map[rune]int)
	rns := []rune(s)
	l, res := 0, 1
	for r := 0; r < len(s); r++ {
		if cnt[rns[r]] != 0 {
			for l <= r && cnt[rns[r]] != 0 {
				cnt[rns[l]]--
				l++
			}
		}
		cnt[rns[r]]++
		res = mymax(res, r-l+1)
	}
	return res
}

func main() {
	tests := []struct {
		input    string
		expected int
	}{
		{"abcabcbb", 3},
		{"cccccccc", 1},
		{"pwwkew", 3},
		{" ", 1},
		{"", 0},
		{"abcdef", 6},
		{"bacabcd", 4},
	}

	for _, tt := range tests {
		result := lengthOfLongestSubstring(tt.input)
		status := "PASS"
		if result != tt.expected {
			status = "FAIL"
		}
		fmt.Printf("[%s] input=%q expected=%d got=%d\n", status, tt.input, tt.expected, result)
	}
}
