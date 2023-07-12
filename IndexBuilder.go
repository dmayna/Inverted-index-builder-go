package main

import (
  "path/filepath"
  "fmt"
  "os"
  "bufio"
  "io/ioutil"
  "log"
  "github.com/kljensen/snowball"
)

type InvertedIndex struct {
	Map map[string]map[string][]uint16
	WordCount map[string]uint16
  }
  
  func InitInvertedIndex() * InvertedIndex {
	index := InvertedIndex{}
	index.Map = map[string]map[string][]uint16{}
	index.WordCount = map[string]uint16{}
	return &index
  }
  
  func (index *InvertedIndex) getWordCount(location string) uint16 {
	return index.WordCount[location]
  }
  
  func (index *InvertedIndex) Add(word string, location string, position uint16) {
	if(index.Map[word] == nil) {
	  index.Map[word] = make(map[string][]uint16)
	}
	index.Map[word][location] = append(index.Map[word][location], position)
	index.WordCount[location] = position
	fmt.Println("Added word to index")
  }
  
  func (index *InvertedIndex) getLocationsOfWord(word string) []string {
	keys := make([]string, 0, len(index.Map[word]))
	for k := range index.Map[word] {
		keys = append(keys, k)
	}
	return keys
  }
  
  func (index *InvertedIndex) getWords() []string {
	keys := make([]string, 0, len(index.Map))
	for k := range index.Map {
		keys = append(keys, k)
	}
	return keys
  }
  
  func (index *InvertedIndex) getLocations(word string, location string) []uint16 {
	keys := make([]uint16, 0, len(index.Map[word][location]))
	for _, k := range index.Map[word][location] {
		keys = append(keys, k)
	}
	return keys
  }
  
  func (index *InvertedIndex) containsWord(word string) bool {
	if _, ok := index.Map[word]; ok {
	  return true
	}
	return false
  }
  
  func (index *InvertedIndex) containsWordLocation(word string, location string) bool {
	if _, ok := index.Map[word][location]; ok {
	  return true
	}
	return false
  }
  
  func (index *InvertedIndex) contains(word string, location string, position uint16) bool {
	for _, item := range index.Map[word][location] {
	  if item == position {
		return true
	  }
  
	}
	return false
  }
  
  func (index *InvertedIndex) wordSize() int {
	return len(index.Map)
  }
  
  func (index *InvertedIndex) locationSize(word string) int {
	return len(index.Map[word])
  }
  
  func (index *InvertedIndex) positionSize(word string, location string) int {
	return len(index.Map[word][location])
  }


/**
* Builds Inverted Index
*
* @param path Path of file or directory
* @param index Inverted index to use 
*/
func buildIndex(path string, index *InvertedIndex) {
	info, err := os.Stat(path)
	if err != nil {
    	fmt.Println(err)
    	return
	}

	if info.IsDir() {
		files, _ := ioutil.ReadDir(path)
		for _, file := range files {
			parseFile(filepath.Join(path, file.Name()), index)
		}
	} else {
		parseFile(path, index)
	}
}

/**
* parses a file and adds every stem word to inverted index 
*
* @param file the file to parse
* @param index the inverted index to add to
*
*/
func parseFile(file string, index *InvertedIndex) {
	oFile, err := os.Open(file)
    if err != nil {
        log.Fatal(err)
    }
	Scanner := bufio.NewScanner(oFile)
    Scanner.Split(bufio.ScanWords)
	position := uint16(1)
	location := file

	for Scanner.Scan() {
		stem, _ := snowball.Stem(Scanner.Text(), "english", true)
		index.Add(location, stem, position)
		fmt.Println("Location: ", location, " | word: ", stem, " | position: ", position)
		position++
		//fmt.Println(snowball.Stem(Scanner.Text(), "english", false))
    }
    if err := Scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func main() {
	index := InitInvertedIndex()
	buildIndex("/path/to/directory/or/file", index)
	fmt.Println(index)
}