package gh

import (
	"log"
	"os"
)

func writeFile(path string, data []byte) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(string(data))
}

func SplitPages(data string) []string {
	start := 0
	parts := []string{}
	stackDepth := 0
	for index, runeValue := range data {
		switch runeValue {
		case '{':
			if stackDepth == 0 {
				start = index
			}
			stackDepth++
		case '}':
			if stackDepth > 0 {
				stackDepth--
			}
			if stackDepth == 0 {
				parts = append(parts, data[start:index+1])
			}
		}
	}
	return parts
}
