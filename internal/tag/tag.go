package tag

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func GetTags(filename string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("cannot open file")
		fmt.Printf("%s\n", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#+filetags: ") {
			if strings.Contains(line, ":gist:") {
				tmp := strings.TrimPrefix(line, "#+filetags: ")
				tmp = strings.Replace(tmp, ":gist:", "", 1)
				results <- tmp
			}
			return
		}
	}
}

func PrintUnsortedTags(results chan string) {
	var uniqueTags []string

	for result := range results {
		tags := strings.Split(result, ":")
		for _, tag := range tags {
			found := false
			for _, uniqueTag := range uniqueTags {
				if tag == uniqueTag {
					found = true
					break
				}
			}
			if ! found && tag != "" {
				uniqueTags = append(uniqueTags, tag)
			}
		}
	}

	for _, tag := range uniqueTags {
		fmt.Println(tag)
	}
}
