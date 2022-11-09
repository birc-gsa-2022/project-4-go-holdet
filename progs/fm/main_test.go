package main

import (
	"fmt"
	"testing"

	"birc.au.dk/gsa/shared"
)

func Test_Output(t *testing.T) {

	genome, read := "akasakan$", "aka"

	sa := shared.LsdRadixSort(genome)
	fmt.Println(sa)
	for i, idx := range sa {
		fmt.Println(i, genome[idx:])
	}
	lower, upper := shared.BinarySearch(genome, read, sa)
	fmt.Println("exact matches are in the interval:", lower, "to", upper)
	for _, v := range sa {
		idx := sa[v]

		//check if all suffixes in the interval matches and that all suffixes outside do not match.
		if v >= lower && v < upper {
			if genome[idx:len(read)+idx] != read {
				t.Error("They ARE NOT identical. But should be at idx:", v)
			}
		} else {
			if sa[v]+len(read) < len(genome) {
				if genome[idx:len(read)+idx] == read {
					t.Error("They ARE identical. But should be at idx:", v)
				}
			}
		}
	}

}

func TestVaryingAlphabets(t *testing.T) {
	matches := 0
	Alphabets := []shared.Alphabet{
		shared.English, shared.DNA, shared.AB}

	for _, v := range Alphabets {
		genome, reads := shared.BuildSomeFastaAndFastq(300, 8, 20, v, 11)

		parsedGenomes := shared.GeneralParserStub(genome, shared.Fasta, len(genome)+1)
		parsedReads := shared.GeneralParserStub(reads, shared.Fastq, len(reads)+1)

		//iterate all genomes
		for _, gen := range parsedGenomes {
			sa := shared.LsdRadixSort(gen.Rec)
			//iterate all reads
			for _, read := range parsedReads {
				lower, upper := shared.BinarySearch(gen.Rec, read.Rec, sa)
				matches += upper - lower

				//verify that the bounds are correct
				for _, v := range sa {
					idx := sa[v]

					//check if all suffixes in the interval matches and that all suffixes outside do not match.
					if v >= lower && v < upper {
						if gen.Rec[idx:len(read.Rec)+idx] != read.Rec {
							t.Error("They ARE NOT identical. But should be at idx:", v)
						}
					} else {
						if sa[v]+len(read.Rec) < len(gen.Rec) {
							if gen.Rec[idx:len(read.Rec)+idx] == read.Rec {
								t.Error("They ARE identical. But should be at idx:", v)
							}
						}
					}
				}
			}
		}

	}
	fmt.Println("a total of", matches, " matches was found in the test.")
}

func TestBWT(t *testing.T) {
	genome, pattern := "mississippi$", "iss"

	sa := shared.LsdRadixSort(genome)
	bwt, c, o := shared.FM_build(sa, genome)
	fmt.Println(bwt, c, o)
	//shared.FMIndexMatching(bwtx, buckets, o)
	fmt.Println("bongo")
	shared.FM_search(bwt, c, o, pattern)

}

func TestOTable(t *testing.T) {
	genome := `missijnsooofsjkfndsjkfndsjkhmkdslfjsdlfksdalkfjdsfloooooo 
			   sooooookndsjfndsjkfopopopopoooooooppppooooooooooossippids
			   oooooookndsjfndsjkfopopggopoooasdappppooooooooooossippis$`
	sa := shared.LsdRadixSort(genome)
	bwt, _, o := shared.FM_build(sa, genome)

	for i := 0; i <= len(genome); i++ {
		counts := make(map[byte]int)
		for _, v := range []byte(bwt[:i]) {
			counts[v]++

		}
		for k, v := range counts {
			if o[i][k] != v {
				t.Errorf("O table error")
			}
		}

	}
}

/*
func Test_cmp_with_old_handin(t *testing.T) {
	shared.SortFile("./testdata/output.txt")
	shared.SortFile("./testdata/handin3_reference.txt")
	if !shared.CmpFiles("./testdata/test_result.txt", "./testdata/h1_naive_results.txt") {
		t.Errorf("files are not identical")
	}
}
*/

/*
func TestMakeDataCons(t *testing.T) {
	csvFile, err := os.Create("./testdata/construction_time.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"x_size", "quadratic"})
	num_of_n := 0
	time_sq := 0
	for i := 1; i < 20; i++ {
		num_of_n += 500
		num_of_m := 1
		genome, _ := shared.BuildSomeFastaAndFastq(num_of_n, 0, 1, shared.English, 78)
		parsedGenomes := shared.GeneralParserStub(genome, shared.Fasta, num_of_n*num_of_m+1)
		//parsedReads := shared.GeneralParserStub(reads, shared.Fastq, num_of_n*num_of_m+1)
		for i := 0; i < 5; i++ {
			for _, gen := range parsedGenomes {
				time_start := time.Now()
				shared.LsdRadixSort(gen.Rec)
				time_end := int(time.Since(time_start))
				time_sq += time_end
				fmt.Println("time", int((time_sq)))
				_ = csvwriter.Write([]string{strconv.Itoa(num_of_n), strconv.Itoa(time_sq)})
				csvwriter.Flush()
				time_sq = 0
			}
		}
	}
}
func TestMakeDataSearch(t *testing.T) {
	csvFile, err := os.Create("./testdata/search_time.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"x_size", "quadratic"})
	//num_of_n := 4000
	num_of_m := 0
	time_sq := 0
	//always use the same genome in order to make the sa process go faster.
	genomes, _ := shared.BuildSomeFastaAndFastq(50000, 0, 1, shared.A, 102)
	parsedGenomes := shared.GeneralParserStub(genomes, shared.Fasta, 50000+1)
	if len(parsedGenomes) != 1 {
		t.Errorf("should only be 1.")
	}
	gen := parsedGenomes[0]
	fmt.Println("creating sa")
	sa := shared.LsdRadixSort(gen.Rec)
	fmt.Println("sa created")
	for i := 1; i < 100; i++ {
		//num_of_n += 500
		num_of_m += 500
		_, reads := shared.BuildSomeFastaAndFastq(50000, num_of_m, 1, shared.A, 102)
		parsedReads := shared.GeneralParserStub(reads, shared.Fastq, 40000*num_of_m+1)
		for i := 0; i < 5; i++ {
			for _, read := range parsedReads {
				time_start := time.Now()
				shared.BinarySearch(gen.Rec, read.Rec, sa)
				time_end := int(time.Since(time_start))
				time_sq += time_end
			}
			fmt.Println("time", int((time_sq)))
			_ = csvwriter.Write([]string{strconv.Itoa(num_of_m), strconv.Itoa(time_sq)})
			csvwriter.Flush()
			time_sq = 0
		}
	}
}
*/
/*
func TestMakeDataOneMore(t *testing.T) {
	csvFile, err := os.Create("./testdata/search_time.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"x_size", "fixed_log", "fixed_log2"})
	num_of_m := 0
	num_of_n := 50000
	time_sq, time_sq2 := 0, 0
	genome, _ := shared.BuildSomeFastaAndFastq(num_of_n, 500, 1, shared.A, 78)
	parsedGenomes := shared.GeneralParserStub(genome, shared.Fasta, num_of_n+1)
	gen := parsedGenomes[0].Rec
	sa := shared.LsdRadixSort(gen)
	for i := 1; i < 51; i++ {
		num_of_m += 500
		_, reads := shared.BuildSomeFastaAndFastq(num_of_m, num_of_m, 1, shared.A, 78)
		var sb strings.Builder
		sb.WriteString(reads[:len(reads)-1])
		sb.WriteRune('b')
		reads = sb.String()
		sb.Reset()
		_, reads2 := shared.BuildSomeFastaAndFastq(num_of_m/2, num_of_m/2, 1, shared.A, 78)
		sb.WriteString(reads2[:len(reads2)-1])
		sb.WriteRune('b')
		reads2 = sb.String()
		sb.Reset()
		parsedReads := shared.GeneralParserStub(reads, shared.Fastq, num_of_n*num_of_m+1)
		parsedReads2 := shared.GeneralParserStub(reads2, shared.Fastq, num_of_n*num_of_m+1)
		for i := 0; i < 5; i++ {
			for _, read := range parsedReads {
				time_start := time.Now()
				shared.BinarySearch(gen, read.Rec, sa)
				time_end := int(time.Since(time_start))
				time_sq += time_end
			}
			for _, read := range parsedReads2 {
				time_start := time.Now()
				shared.BinarySearch(gen, read.Rec, sa)
				time_end := int(time.Since(time_start))
				time_sq2 += time_end
			}
			fmt.Println("time", int((time_sq)))
			_ = csvwriter.Write([]string{strconv.Itoa(num_of_m), strconv.Itoa(time_sq), strconv.Itoa(time_sq2)})
			csvwriter.Flush()
			time_sq, time_sq2 = 0, 0
		}
	}
}
*/
