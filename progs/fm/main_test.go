package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestMakeDataPP(t *testing.T) {
	csvFile, err := os.Create("./data/search_time.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"x_size", "time"})

	num_of_n := 1000

	//num_of_m := 0

	for i := 1; i < 3; i++ {
		num_of_n += 1000
		genomes, _ := shared.BuildSomeFastaAndFastq(num_of_n, 0, 1, shared.English, 102)
		parsedGenomes := shared.GeneralParserStub(genomes, shared.Fasta, 50000+1)
		if len(parsedGenomes) != 1 {
			t.Errorf("should only be 1.")
		}

		gen := parsedGenomes[0]
		fmt.Println("creating sa")
		sa := shared.LsdRadixSort(gen.Rec)
		fmt.Println("sa created")

		for i := 0; i < 10; i++ {
			time_start := time.Now()

			f, err := os.Create("./data/banko.dk")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			var sb strings.Builder
			//add sentinel if missing
			if gen.Rec[len(gen.Rec)-1] != '$' {
				sb.WriteString(gen.Rec)
				sb.WriteRune('$')
				gen.Rec = sb.String()
			}
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
			time_end := int(time.Since(time_start))
			fmt.Println("time", int((time_end)))
			_ = csvwriter.Write([]string{strconv.Itoa(num_of_n), strconv.Itoa(time_end)})
			csvwriter.Flush()
		}
	}
}

/*
func TestMakeDataSearch(t *testing.T) {
	csvFile, err := os.Create("./data/search_time_2.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"x_size", "time"})

	num_of_n := 30000

	//num_of_m := 0

	genomes, _ := shared.BuildSomeFastaAndFastq(num_of_n, 0, 1, shared.A, int64(num_of_n)+2)
	parsedGenomes := shared.GeneralParserStub(genomes, shared.Fasta, 50000+1)
	gen := parsedGenomes[0]
	fmt.Println("creating sa")
	sa := shared.LsdRadixSort(gen.Rec)
	fmt.Println("sa created")

	f, err := os.Create("./data/banko.dk")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var sb strings.Builder
	//add sentinel if missing
	if gen.Rec[len(gen.Rec)-1] != '$' {
		sb.WriteString(gen.Rec)
		sb.WriteRune('$')
		gen.Rec = sb.String()
	}
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

	f, err = os.Open("./data/banko.dk")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	p_genome := shared.FMParser(f)

	gen_unp := p_genome[0]
	//find the original genome
	gen_unp.BS = shared.ReverseBWT(gen_unp.Bwt, gen_unp.C, gen_unp.O)

	num_of_m := 0

	for i := 1; i < 4; i++ {
		num_of_m += 1000
		_, reads := shared.BuildSomeFastaAndFastq(num_of_m+2, num_of_m, 1, shared.A, int64(num_of_m)*2+5)

		parsedRead := shared.GeneralParserStub(reads, shared.Fastq, 50000+1)
		read := parsedRead[0]

		if len(parsedGenomes) != 1 {
			t.Errorf("should only be 1.")
		}

		for i := 0; i < 10; i++ {
			time_start := time.Now()
			_, _ = shared.FM_search(gen_unp.Bwt, gen_unp.C, gen_unp.O, read.Rec)
			time_end := int(time.Since(time_start))
			fmt.Println("time", int((time_end)))
			_ = csvwriter.Write([]string{strconv.Itoa(num_of_m), strconv.Itoa(time_end)})
			csvwriter.Flush()
		}
	}
}
*/
