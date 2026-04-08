package common

func SuggestPath(unknownPath string, expectedPaths map[string]struct{}) (string, bool) {
	if unknownPath == "" || len(expectedPaths) == 0 {
		return "", false
	}

	closestPath := ""
	closestDistance := int(^uint(0) >> 1)
	for expectedPath := range expectedPaths {
		distance := levenshteinDistance(unknownPath, expectedPath)
		if distance < closestDistance || (distance == closestDistance && (closestPath == "" || expectedPath < closestPath)) {
			closestPath = expectedPath
			closestDistance = distance
		}
	}

	if closestPath == "" {
		return "", false
	}

	maxAllowedDistance := len(unknownPath) / 3
	maxAllowedDistance = max(1, maxAllowedDistance)
	maxAllowedDistance = min(3, maxAllowedDistance)
	if closestDistance > maxAllowedDistance {
		return "", false
	}

	return closestPath, true
}

func levenshteinDistance(left, right string) int {
	leftRunes := []rune(left)
	rightRunes := []rune(right)

	if len(leftRunes) == 0 {
		return len(rightRunes)
	}
	if len(rightRunes) == 0 {
		return len(leftRunes)
	}

	previousRow := make([]int, len(rightRunes)+1)
	currentRow := make([]int, len(rightRunes)+1)
	for rightIndex := 0; rightIndex <= len(rightRunes); rightIndex++ {
		previousRow[rightIndex] = rightIndex
	}

	for leftIndex := 1; leftIndex <= len(leftRunes); leftIndex++ {
		currentRow[0] = leftIndex
		for rightIndex := 1; rightIndex <= len(rightRunes); rightIndex++ {
			substitutionCost := 0
			if leftRunes[leftIndex-1] != rightRunes[rightIndex-1] {
				substitutionCost = 1
			}

			currentRow[rightIndex] = min(
				previousRow[rightIndex]+1,
				min(currentRow[rightIndex-1]+1, previousRow[rightIndex-1]+substitutionCost),
			)
		}
		previousRow, currentRow = currentRow, previousRow
	}

	return previousRow[len(rightRunes)]
}
