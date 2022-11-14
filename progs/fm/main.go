package main

import (
	"fmt"
	"os"
	"strings"

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
		//create file to store preprocessing (data/pre.fa)
		f, err := os.Open(os.Args[2])
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// preprocess
		genome := os.Args[2]
		p_genomes := shared.GeneralParser(genome, shared.Fasta)

		f, err = os.Create(os.Args[2])
		if err != nil {
			panic(err)
		}
		defer f.Close()

		for _, gen := range p_genomes {
			var sb strings.Builder
			//add sentinel if missing
			if gen.Rec[len(gen.Rec)-1] != '$' {
				sb.WriteString(gen.Rec)
				sb.WriteRune('$')
				gen.Rec = sb.String()
			}
			sa := shared.LsdRadixSort(gen.Rec)
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
		//fmt.Println(shared.TodoPreprocess(os.Args[2]))
	} else {
		//fmt.Println(shared.TodoMap(os.Args[1], os.Args[2]))

		//perform exact pattern matching on already precomputed data
		reads := os.Args[2]

		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		p_genomes := shared.FMParser(file)
		p_reads := shared.GeneralParser(reads, shared.Fastq)

		fo, err := os.Create("./data/output.txt")
		if err != nil {
			panic(err)
		}
		for _, gen := range p_genomes {
			for _, read := range p_reads {
				start, end := shared.FM_search(gen.Bwt, gen.C, gen.O, read.Rec)
				if start != end {
					if len(gen.BS) == 0 {
						//this is only computed if needed
						gen.BS = shared.ReverseBWT(gen.Bwt, gen.C, gen.O)
					}
					for i := start; i < end; i++ {

						shared.Sam(read.Name, gen.Name, gen.BS[i], read.Rec)

						res := shared.SamStub(read.Name, gen.Name, gen.BS[i], read.Rec)
						fo.Write([]byte(res))

					}
				}
			}
		}
	}
}
