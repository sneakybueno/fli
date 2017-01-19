package shell

import "fmt"

func FindNextTerm(values []string, term string) (string, error) {
	low := 0
	high := len(values)

	if low < high {
		mid := (high - low) / 2

		// check length before checking if it is a match and find new mid

		comparison := firstNCharactersMatch(len(term), term, values[mid])

		if comparison < 0 {
			return FindNextTerm(values[:mid], term)
		} else if comparison > 0 {
			return FindNextTerm(values[mid+1:], term)
		} else if comparison == 0 {
			value := values[mid]

			for i := mid - 1; i >= 0; i-- {
				if len(values[i]) < len(term) {
					continue
				}

				if firstNCharactersMatch(len(term), term, values[i]) != 0 {
					break
				}

				value = values[i]
			}
			return value, nil
		}
	}

	return "", fmt.Errorf("Error: No matching value for %s ", term)
}

// firstNCharactersMatch compares the first N
// characters of the given strings. Returns:
// 0 if a == b
// -1 if a comes before b in dictionary order
// 1 if a comes after b in dictorary ordefr
func firstNCharactersMatch(n int, a string, b string) int {
	for i := 0; i < n; i++ {
		if i >= len(a) || i >= len(b) {
			return 0
		}

		if a[i] == b[i] {
			continue
		} else if a < b {
			return -1
		} else if a > b {
			return 1
		}
	}

	return 0
}
