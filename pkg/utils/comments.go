package utils

import "math"

func SplitComment(comment string, maxSize int, sepEnd string, sepStart string, maxCommentsPerCommand int, truncationHeader string) []string {
	if len(comment) <= maxSize {
		return []string{comment}
	}

	// No comment contains both sepEnd and truncationHeader, so we only have to count their max.
	maxWithSep := maxSize - max(len(sepEnd), len(truncationHeader)) - len(sepStart)
	var comments []string
	numPotentialComments := int(math.Ceil(float64(len(comment)) / float64(maxWithSep)))
	var numComments int
	if maxCommentsPerCommand == 0 {
		numComments = numPotentialComments
	} else {
		numComments = min(numPotentialComments, maxCommentsPerCommand)
	}
	isTruncated := numComments < numPotentialComments
	upTo := len(comment)
	for len(comments) < numComments {
		downFrom := max(0, upTo-maxWithSep)
		portion := comment[downFrom:upTo]
		if len(comments)+1 != numComments {
			portion = sepStart + portion
		} else if len(comments)+1 == numComments && isTruncated {
			portion = truncationHeader + portion
		}
		if len(comments) != 0 {
			portion = portion + sepEnd
		}
		comments = append([]string{portion}, comments...)
		upTo = downFrom
	}
	return comments
}
