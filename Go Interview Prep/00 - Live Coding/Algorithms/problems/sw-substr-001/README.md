# Problem sw-substr-001: Longest Substring Without Repeating Characters

**Topic:** sliding window
**Subtopic:** hashmap
**Difficulty:** medium
**Time limit:** 15 мин

## Описание

Дана строка `s`. Нужно найти длину самой длинной подстроки без повторяющихся символов.

Строка `s` состоит из английских букв, цифр, символов и пробелов.

## Requirements

- Сигнатура: `func lengthOfLongestSubstring(s string) int`
- Оптимальное решение за O(n) по времени
- Допустимая память: O(min(n, m)), где m — размер алфавита

## Examples

```
Input:  "abcabcbb"
Output: 3
Explanation: "abc"

Input:  "cccccccc"
Output: 1
Explanation: "c"

Input:  "pwwkew"
Output: 3
Explanation: "wke"

Input:  " "
Output: 1
Explanation: " "

Input:  ""
Output: 0
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Используй два указателя (left, right) — это классический sliding window. Двигай right вправо, расширяя окно, пока не встретишь повтор.
</details>

<details>
<summary>Hint 2</summary>
Храни символы текущего окна в map[byte]int (символ → его последняя позиция). Когда встречаешь повтор — сдвинь left за позицию предыдущего вхождения этого символа.
</details>

<details>
<summary>Hint 3</summary>
Не забывай обновлять максимум на каждом шаге: max = max(max, right - left + 1).
</details>
