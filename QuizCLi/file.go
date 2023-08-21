package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	//gets the specific csv file from the user frpm the cli using the flag syntax
	//by default the file that this program reads is problems.csv
	csvFilename := flag.String("csv", "problems.csv", "a csv  file with the format of 'question,answer'")
	flag.Parse()
	timeLimit := flag.Int("limit", 30, "a time limit in between each problem")

	//opens the file specified by the flag/ default
	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open the CSV file: %s\n", *csvFilename)
		os.Exit(1)
	}
	//creating a reader that is able to scan the opened csv file
	r := csv.NewReader(file)
	//defining the variable lines as the placeholder of all the lines read
	lines, err := r.ReadAll()
	if err != nil {
		exit("failed to parse required file")
	}
	//converts lines vector into struct objects
	problems := parseLines(lines)
	//count keeps the score of the user tru the programs lifetime
	count := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// loop traverses thru the range of all problems where at every iteration we utilize the q & a parameter
	for i, p := range problems {

		fmt.Printf("Problem #%d: %s = \n ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("Problems finished , you scored %d out of %d", count, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				count++
			}
		}
		//we use &answer to initialise the answer variable as the placeholder of scanf for each iteration

	}

	fmt.Printf("Problems finished , you scored %d out of %d", count, len(problems))

}

func exit(s string) {
	fmt.Print(s)
	os.Exit(1)
}

// parseLines() takes in a 2D vector (in our case it's a csv) and returns a column space of structs of 1D
// this struct is basically an object in which we can specify the variables as question and answer
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}
