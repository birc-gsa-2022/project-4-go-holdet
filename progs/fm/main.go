package main

import (
	"fmt"
	"os"

	// Directories in the root of the repo can be imported
	// as long as we pretend that they sit relative to the
	// url birc.au.dk/gsa, like this for the example 'shared':
	"birc.au.dk/gsa/shared"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-p genome")
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "genome reads")
		os.Exit(1)
	}
	if os.Args[1] == "-p" {
		// preprocess

		genome := os.Args[2]
		p_genomes := shared.GeneralParser(genome, shared.Fasta)
		for _, gen := range p_genomes {
			sa := shared.LsdRadixSort(gen.Rec)
			bwtx, buckets, o := shared.BWTPreprocessing(sa, gen.Rec)
			shared.FMIndexMatching(bwtx, buckets, o)
		}
		fmt.Println(shared.TodoPreprocess(os.Args[2]))
	} else {
		fmt.Println(shared.TodoMap(os.Args[1], os.Args[2]))
	}
	genome := os.Args[1]
	reads := os.Args[2]

	p_genomes := shared.GeneralParser(genome, shared.Fasta)
	p_reads := shared.GeneralParser(reads, shared.Fastq)

	/*
		fo, err := os.Create("./testdata/output.txt")
		if err != nil {
			panic(err)
		}*/

	for _, gen := range p_genomes {
		sa := shared.LsdRadixSort(gen.Rec)
		fmt.Println(sa)

		for _, read := range p_reads {
			start, end := shared.BinarySearch(gen.Rec, read.Rec, sa)
			for i := start; i < end; i++ {
				shared.Sam(read.Name, gen.Name, sa[i], read.Rec)
				/*
					res := shared.SamStub(read.Name, gen.Name, sa[i], read.Rec)
					fo.Write([]byte(res))
				*/
			}
		}

	}

}
