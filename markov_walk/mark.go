package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

//inner is the map of suffix to int
//whole is the entire map of prefix- suffix- frequency

func main() {

	command := os.Args[1] //build or generate

	if command == "read" {
		numberOfPrefixString := os.Args[2]                           //how many word for prefix
		numberOfPrefixInt, err := strconv.Atoi(numberOfPrefixString) //convert prefix# into int
		if err != nil {
			fmt.Print("number of prefix wrong")
		}

		freqMapFileName := os.Args[3]
		inputTextFiles := os.Args[4:]                                   //all the input text filenames as array
		CreateChain(inputTextFiles, numberOfPrefixInt, freqMapFileName) //make frequency map txt file
	}

	if command == "generate" {
		freqMapFile := os.Args[2]
		numGeneratedWords := os.Args[3]
		numGeneratedWordsInt, err := strconv.Atoi(numGeneratedWords)
		if err != nil {
			fmt.Print("number of words generated wrong")
		}
		GenerateWords(freqMapFile, numGeneratedWordsInt)
	}

}

//generate words from model frequency map text file
func GenerateWords(freqFilename string, numGenerated int) {
	rand.Seed(time.Now().Unix())
	// open file
	// read file and parse into individual strings
	// make a map (ready line by line?)
	// (when output, erase "" into emptyspace )

	//open file
	fd, err := os.Open(freqFilename)
	if err != nil {
		fmt.Println("there is an error opening the frequency table file")
	}
	defer fd.Close()

	//read file and store each line as one big string
	var lineByLine []string
	reader := bufio.NewReader(fd)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error readling file")
			break
		}
		lineByLine = append(lineByLine, line)
	}

	freqMap := make(map[string]map[string]int) //prefix-suffix frequency map

	length := strings.Replace(lineByLine[0], "\n", "", -1)
	prefixLength, err := strconv.Atoi(length)
	if err != nil {
		fmt.Println("cannot read number of prefixes")
	}

	for i := 1; i < len(lineByLine); i++ { // from the second line till the end of file
		splitLine := strings.Split(lineByLine[i], " ") // split line is an array of words in each line

		prefix := strings.Join(splitLine[0:prefixLength], " ")
		suffixFreq := make(map[string]int)

		for j := 0; j < (len(splitLine)-prefixLength)/2; j++ {
			countString := splitLine[prefixLength+2*j+1]
			countInt, err := strconv.Atoi(countString)
			if err != nil {
				fmt.Println("count data cannot be read")
			}
			suffixFreq[splitLine[prefixLength+2*j]] = countInt
		}

		freqMap[prefix] = suffixFreq
		// initialize split words to prefix
		// loop through each line and add suffix - int pairs to suffixFreq map

	}
	//fmt.Print(freqMap)

	// random seed
	// generate random int between 0 to maplength-1
	var words []string
	var firstPrefix string

	// maps are not ordered, first key chosen will be different
	for k, _ := range freqMap {
		firstPrefix = k
		break
	}

	splitFirstPrefix := strings.Split(firstPrefix, " ")
	words = append(words, splitFirstPrefix...) // add split prefix to word string
	// make slice with all suffix

	// for each suffix
	for j := 0; j < (numGenerated - prefixLength); j++ {
		prefix := strings.Join(words[(len(words)-prefixLength):], " ")
		var suffixList []string
		//read through suffix of the given prefix
		for kk, vv := range freqMap[prefix] {

			// add suffix to suffix list (suffix count)number of times
			for i := 0; i < vv; i++ {
				suffixList = append(suffixList, kk)
			}
		}
		suffixIndex := rand.Intn(len(suffixList))
		words = append(words, suffixList[suffixIndex])
	}

	for i := 0; i < len(words); i++ {
		if words[i] == "\"\"" {
			words[i] = " "
		}
	}
	fmt.Println(words)
	//WHEN PRINT ERASE ""
	//use freqMap to generate numGenerated number of words
	//randomly choose key from map
	//how to find value of specific key in map
	//before printing, get rid of double quotation marks

	// make a slice of all keys from map
	// pick a random element of the slice
	// for k v range freqmap --

}

// open file and output an array of words of the input text file
func CreateChain(filenames []string, prefixLength int, outputFilename string) {

	//whole is the entire map of prefix- suffix- frequency
	whole := make(map[string]map[string]int) // combine all the prefixes into one string (because map cannot take slice as keys and arrays need to have a set size)
	var allWords []string                    //all the strings of all files must be added here

	for i := 0; i < (len(filenames)); i++ { //loop through each text file
		b, err := ioutil.ReadFile(filenames[i])
		if err != nil {
			fmt.Print(err)
		}
		entireFileString := string(b)
		words := strings.Fields(entireFileString)
		allWords = append(allWords, words...)
	}
	//for i := 0; i < len(allWords)-1; i++ {
	//	fmt.Printf(allWords[i])
	//}

	// set prefix into prefix length number of double quotation ""
	// make string with ""
	// make an empty slice
	// with a for loop make an array of "" strings with the length of prefix length
	// join array into one string := prefix
	emptyString := "\"\""
	var slice []string
	for i := 0; i < prefixLength; i++ {
		slice = append(slice, emptyString)
	}
	prefix := strings.Join(slice, " ")

	for i := 0; i < len(allWords); i++ {
		var suffix string = allWords[i]

		if _, ok := whole[prefix]; ok { //if the prefix is part of the whole map
			if _, kk := whole[prefix][suffix]; kk { //if the prefix and suffix combination exists
				whole[prefix][suffix] = whole[prefix][suffix] + 1 // add +1 to frequeny to existing prefix-suffix pair
			} else { //prefix exists but not with suffix pair
				whole[prefix][suffix] = 1 //initialize new prefix-suffix
			}

		} else { //not even the prefix is part of the whole map
			whole[prefix] = make(map[string]int) // make suffix map
			whole[prefix][suffix] = 1

		}
		//shift everything to left to make room for new suffix
		//join the two prefix together

		//make array of prefixes + suffix
		splitPrefix := strings.Split(prefix, " ")
		copy(splitPrefix, splitPrefix[1:])
		splitPrefix[len(splitPrefix)-1] = suffix

		//join
		prefix = strings.Join(splitPrefix, " ")
		//fmt.Println(prefix) // check if prefix correct
	}

	// print map in desired format prefix suffix int suffix int ...
	/*for k, v := range whole {
		fmt.Printf("%s ", k)
		for kk, vv := range v {
			fmt.Printf("%s %d ", kk, vv)
		}
		fmt.Printf("\n")
	}*/

	//create frequency map file
	file, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println("error: cannot create frequency map file")
	}
	defer file.Close()

	//save map in desired format prefix suffix int suffix int ...
	fmt.Fprintln(file, prefixLength)
	for k, v := range whole {
		fmt.Fprintf(file, "%s ", k)
		for kk, vv := range v {
			fmt.Fprintf(file, "%s %d ", kk, vv)
		}
		fmt.Fprintf(file, "\n")
	}
}
