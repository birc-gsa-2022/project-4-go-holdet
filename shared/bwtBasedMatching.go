package shared

import (
	"fmt"
	"sort"
)

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
func BWTPreprocessing(sa []int, genome string) ([]byte, map[byte]int, []map[byte]int) {
	bwtx := make([]byte, len(sa))
	counts := make(map[byte]int)
	buckets := make(map[byte]int)
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

		if v == 0 {
			bwtx[i] = genome[len(sa)-1]
		} else {
			bwtx[i] = genome[v-1]
		}
		counts[bwtx[i]] += 1
	}

	o[len(sa)] = counts
	keys := getSortedKeysOfCountSlice(counts)
	for i, v := range keys {
		if i != 0 {
			buckets[v] = counts[keys[i-1]] + buckets[keys[i-1]]
		}
	}
	return bwtx, buckets, o
}

// this function just prints my stuff the parameters given. Should do the pattern matching.
func FMIndexMatching(bwtx []byte, bucketArray map[byte]int, Oslice []map[byte]int) {
	fmt.Println("bwtx")
	fmt.Println(bwtx)
	fmt.Println("")
	fmt.Println("bucketArray")
	fmt.Println(bucketArray)
	for _, v := range bucketArray {
		print(v)
	}
	fmt.Println("")
	fmt.Println("Oslice")
	//fmt.Println(Oslice)
	for _, v := range Oslice {
		fmt.Println(v)
	}
}
