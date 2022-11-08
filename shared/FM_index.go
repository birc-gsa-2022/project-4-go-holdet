package shared

import (
	"fmt"
	"sort"
)

/*¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤ THIS IS PROBABLY NOT IDEAL....
¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤ SHOULD PERHAPS USE BYTE TO
¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤ REPRESENT BUCKETS TO AVOID
¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤¤ n log n SEARCH*/
func getSortedKeysOfCountSlice(counts map[byte]int) []byte {
	keys := make([]byte, len(counts))
	i := 0
	for key := range counts {
		keys[i] = key
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// Data might need to represented differently
func FM_build(sa []int, genome string) ([]byte, map[byte]int, []map[byte]int) {
	bwt := make([]byte, len(sa))
	counts := make(map[byte]int)
	c := make(map[byte]int)
	o := make([]map[byte]int, len(sa)+1)
	fmt.Println(len(sa))
	fmt.Println(len(o))
	for i, v := range sa {
		fmt.Println(i)
		copyOfCounts := make(map[byte]int)
		// Copy from the original map to the target map
		for key, value := range counts {
			copyOfCounts[key] = value
		}
		o[i] = copyOfCounts

		//add current letter to o table
		if v == 0 {
			bwt[i] = genome[len(sa)-1]
		} else {
			bwt[i] = genome[v-1]
		}
		counts[bwt[i]] += 1
	}
	//last idx with all values
	o[len(sa)] = counts

	//create buckets
	keys := getSortedKeysOfCountSlice(counts)
	for i, v := range keys {
		if i != 0 {
			c[v] = counts[keys[i-1]] + c[keys[i-1]]
		}
	}

	return bwt, c, o
}

//locate interval for pattern p
func FM_search(bwt []byte, c map[byte]int, o []map[byte]int, p string) (int, int) {
	L := 0
	R := len(bwt) - 1

	for i := len(p) - 1; i >= 0; i-- {
		if L == R {
			return L, R
		}

		a := p[i]

		L = c[a] + o[L][a]
		R = c[a] + o[R][a]
		fmt.Println(L, R)
	}

	for i := len(p) - 1; i > 0; i-- {

	}
	return L, R
}

/*
// this function just prints my stuff the parameters given. Should do the pattern matching.
func FMIndexMatching(bwt []byte, c map[byte]int, Oslice []map[byte]int) {
	fmt.Println("bwt")
	fmt.Println(bwt)
	fmt.Println("")
	fmt.Println("bucketArray")
	fmt.Println(c)
	for _, v := range c {
		print(v)
	}
	fmt.Println("")
	fmt.Println("Oslice")
	//fmt.Println(Oslice)
	for _, v := range Oslice {
		fmt.Println(v)
	}
}
*/
