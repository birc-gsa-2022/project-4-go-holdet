package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"birc.au.dk/gsa/shared"
)

func TestVaryingAlphabets(t *testing.T) {
	matches := 0
	Alphabets := []shared.Alphabet{
		shared.English, shared.DNA, shared.AB}

	for _, v := range Alphabets {
		genome, reads := shared.BuildSomeFastaAndFastq(50, 5, 20, v, 11)
		parsedGenomes := shared.GeneralParserStub(genome, shared.Fasta, len(genome)+1)

		//do preprocessing - write read file
		f, err := os.Create("./test.bongo")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		var sa []int
		for _, gen := range parsedGenomes {
			var sb strings.Builder
			//add sentinel if missing
			if gen.Rec[len(gen.Rec)-1] != '$' {
				sb.WriteString(gen.Rec)
				sb.WriteRune('$')
				gen.Rec = sb.String()
			}
			sa = shared.LsdRadixSort(gen.Rec)
			bwt, c := shared.FM_build(sa, gen.Rec)
			//write to file
			f.WriteString(">" + gen.Name + "\n")
			f.WriteString("@")
			f.Write(bwt)
			f.WriteString("\n")
			for k, v := range c {
				f.WriteString("*" + string(k) + fmt.Sprint(v))
				f.WriteString("\n")
			}

		}
		f, err = os.Open("./test.bongo")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		FMParsedGenomes := shared.FMParser(f)
		parsedReads := shared.GeneralParserStub(reads, shared.Fastq, len(reads)+1)

		//iterate all genomes
		for i, gen := range FMParsedGenomes {
			//iterate all reads
			for _, read := range parsedReads {
				start, end := shared.FM_search(gen.Bwt, gen.C, gen.O, read.Rec)
				if start != end {
					matches += end - start
					//find matches
					if len(gen.BS) == 0 {
						//this is only computed if needed
						gen.BS = shared.ReverseBWT(gen.Bwt, gen.C, gen.O)
					}

					//verify that the bounds are correct
					for _, v := range gen.BS {
						idx := gen.BS[v]
						//check if all suffixes in the interval matches and that all suffixes outside do not match.
						if v >= start && v < end {
							if parsedGenomes[i].Rec[idx:len(read.Rec)+idx] != read.Rec {
								t.Error("They ARE NOT identical. But should be at idx:", v)
							}
						} else {
							if gen.BS[v]+len(read.Rec) < len(parsedGenomes[i].Rec) {
								if parsedGenomes[i].Rec[idx:len(read.Rec)+idx] == read.Rec {
									t.Error("They ARE identical. But should be at idx:", v)
								}
							}
						}
					}
				}
			}

		}
	}
	fmt.Println("a total of", matches, " matches was found in the test.")
}

/*func TestBWT(t *testing.T) {
	genome, pattern := "mississippi$", "iss"

	sa := shared.LsdRadixSort(genome)
	bwt, c := shared.FM_build(sa, genome)
	fmt.Println(bwt, c)
	//shared.FMIndexMatching(bwtx, buckets, o)
	fmt.Println("bongo")
	shared.FM_search(bwt, c, _, pattern)

}*/

func TestOTable(t *testing.T) {
	genome := `dagagagaaa$`
	sa := shared.LsdRadixSort(genome)
	bwt, _ := shared.FM_build(sa, genome)

	for i := 0; i <= len(genome); i++ {
		counts := make(map[byte]int)
		for _, v := range []byte(bwt[:i]) {
			counts[v]++

		}
		o := shared.BuildOtable(bwt)
		for k, v := range counts {
			if o[i][k] != v {
				t.Errorf("O table error")
			}
		}

	}
}

/*
func TestFMParser(t *testing.T) {

	f, er := os.Open("data/pre.fm")
	if er != nil {
		panic(er)
	}
	pr_genomes := shared.FMParser(f)

	fmt.Println(pr_genomes[0].Bwt, pr_genomes[0].C, pr_genomes[0].Name)

}
*/

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
