package shared

import (
	"strings"
)

/*Performs binary search on
First we identify a match. When we the match, we branch out and find upper and lower bound.*/
func BinarySearch(genome string, read string, suffixArray []int) (int, int) {

	//add sentinel if not present
	var sb strings.Builder
	if genome[len(genome)-1] != '$' {
		sb.WriteString(genome)
		sb.WriteByte('$')
		genome = sb.String()
	}

	//edgecase where read is empty or longer than genome. Return empty interval (genome[0:0[)
	if read == "" || len(read) > len(genome) {
		return 0, 0
	}
	upperBound := upperBound(suffixArray, genome, read)
	lowerBound := lowerBound(suffixArray, genome, read)
	if upperBound != lowerBound {
		return lowerBound, upperBound
	}
	return 0, 0

}

func upperBound(suffixArray []int, genome string, read string) int {
	low := 0
	high := len(genome) - 1

	//we break when we are left with two values. The last element that matches, and the first element that does not.
	for low < high-1 {
		mid := low + (high-low)/2
		saIndex := suffixArray[mid]

		if suffixArray[mid]+len(read) < len(genome) {
			if genome[suffixArray[mid]:suffixArray[mid]+len(read)] == read {
				low = mid
				continue
			}
		}
		if genome[saIndex:] < read {
			//mid is too high
			low = mid
		} else {
			//mid is too low
			high = mid
		}

	}

	//We return the idx of first element that is not included, since we define intervals as '[i,j['. Result is either high or high+1.
	if suffixArray[high]+len(read) < len(genome) {
		if genome[suffixArray[high]:suffixArray[high]+len(read)] == read {
			return high + 1
		} else {
			return high
		}
	} else {
		//we look at higher element that is also shorter than our pattern
		return high
	}
}
func lowerBound(suffixArray []int, genome string, read string) int {
	low := 0
	high := len(genome) - 1

	//we break when we are left with two values. The last element before match, and the first match.
	for low < high-1 {
		mid := low + (high-low)/2
		saIndex := suffixArray[mid]

		if suffixArray[mid]+len(read) < len(genome) {
			if genome[suffixArray[mid]:suffixArray[mid]+len(read)] == read {
				high = mid
				continue
			}
		}
		if genome[saIndex:] < read {
			//mid is too low
			low = mid
		} else {
			//mid is too high
			high = mid
		}

	}

	//check if value low is the match, and return if it is. Otherwise match is low+1. (high)
	if suffixArray[low]+len(read) < len(genome) {
		if genome[suffixArray[low]:suffixArray[low]+len(read)] == read {
			return low
		} else {
			return high
		}
	} else {
		return high
	}
}