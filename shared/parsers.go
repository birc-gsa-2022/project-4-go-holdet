// You can create modules at this level and they will be
// interpreted as under module birc.au.dk, so to import
// package `shared` you need `import "birc.au.dk/gsa/shared"`

package shared

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Sam(readName string, chrom string, pos int, readString string) {
	output := readName + "	" + chrom + "	" + fmt.Sprint(pos+1) + "	" + strconv.Itoa(len(readString)) + "M" + "	" + readString + "\n"

	fmt.Print(output)
}

func SamStub(readName string, chrom string, pos int, readString string) string {
	output := readName + "	" + chrom + "	" + fmt.Sprint(pos+1) + "	" + strconv.Itoa(len(readString)) + "M" + "	" + readString + "\n"

	return output
}

type Recs struct {
	Name string
	Rec  string
}

type Format string

const (
	Fasta Format = ">"
	Fastq Format = "@"
)

func GeneralParser(file string, format Format) []Recs {
	f, err := os.Open(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	//###########################################

	output := ""
	fileScanner := bufio.NewScanner(f)
	activeRec := new(Recs)

	recs := make([]Recs, 0)
	//scan file line by line
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(line) == 0 {
			continue
		}

		//handle 'name of sequence' cases
		if line[0:1] == string(format) {
			if len(output) != 0 {
				activeRec.Rec = output
				recs = append(recs, *activeRec)
			}
			output = ""
			activeRec = new(Recs)
			activeRec.Name = strings.TrimSpace(line[1:])
			//handle 'sequence' cases
		} else {
			output = output + line
		}
	}
	activeRec.Rec = output
	recs = append(recs, *activeRec)

	return recs
}

func GeneralParserStub(file string, format Format, maxCapacity int) []Recs {
	output := ""
	fileScanner := bufio.NewScanner(strings.NewReader(file))
	buf := make([]byte, maxCapacity)
	fileScanner.Buffer(buf, maxCapacity)
	activeRec := new(Recs)

	recs := make([]Recs, 0)
	//scan file line by line
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(line) == 0 {
			continue
		}

		//handle 'name of sequence' cases
		if line[0:1] == string(format) {
			if len(output) != 0 {
				activeRec.Rec = output
				recs = append(recs, *activeRec)
			}
			output = ""
			activeRec = new(Recs)
			activeRec.Name = strings.TrimSpace(line[1:])
			//handle 'sequence' cases
		} else {
			output = output + line
		}
	}
	activeRec.Rec = output
	recs = append(recs, *activeRec)

	return recs
}