package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var prefix = "freshReadmeSnippet: "

type base int

const (
	NORMAL base = iota
	COMMENT
	FENCE
	INSIDE_FENCE
	END
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func panicIf(condition bool, fileName string, line int, message string, a ...interface{}) {
	if condition {
		err := fmt.Sprintf(message, a...)
		panic(fmt.Sprintf("%s at %s:%d", err, fileName, line))
	}
}

func includeFile(out *os.File, fileName string) {
	in, err := ioutil.ReadFile(fileName)
	check(err)

	_, err = out.Write(in)
	check(err)

	_, err = out.WriteString("\n")
	check(err)
}

func fromFile(out *os.File, header string, fileName string) {
	in, err := os.Open(fileName)
	check(err)
	defer in.Close()

	var state = NORMAL

	var i = 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		i++
		var line = scanner.Text()

		if state == COMMENT {
			state = FENCE
		}

		if strings.Contains(line, header) {
			if state == NORMAL {
				state = COMMENT
			} else if state == FENCE {
				state = END
			} else {
				panicIf(true, fileName, i, "Snippet `%s` appears second time", header)
			}
		}

		if state == FENCE {
			_, err := out.WriteString(line + "\n")
			check(err)
		}
	}

	panicIf(state != END, fileName, i, "Unable to find snippet `%s`", header)
}

func refresh(fileName string) {
	dir := filepath.Dir(fileName)

	refresh, _ := regexp.Compile("^<!--.*\\[freshReadmeSource\\]\\(([^#]+)#*(.*?)\\)")

	in, err := os.Open(fileName)
	check(err)
	defer in.Close()

	out, err := os.Create(fileName + ".tmp")
	check(err)
	defer out.Close()

	state := NORMAL
	header := ""
	source := ""

	i := 0

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		i++

		if state == FENCE {
			state = INSIDE_FENCE
			if header == "" {
				includeFile(out, filepath.Join(dir, source))
			} else {
				fromFile(out, prefix+header, filepath.Join(dir, source))
			}
			header = ""
			source = ""
		}

		if strings.HasPrefix(line, "```") {
			if state == COMMENT {
				state = FENCE
			} else if state == INSIDE_FENCE {
				state = NORMAL
			} else if state != NORMAL {
				panicIf(true, fileName, i, "Unexpected ```")
			}
		}

		matches := refresh.FindStringSubmatch(line)
		if len(matches) > 0 && state != INSIDE_FENCE {
			panicIf(state != NORMAL, fileName, i, "Unable to process include in include")
			state = COMMENT
			header = matches[2]
			source = matches[1]
		}

		if state != INSIDE_FENCE {
			_, err := out.WriteString(line + "\n")
			check(err)
		}
	}
	check(scanner.Err())

	panicIf(state != NORMAL, fileName, i, "Unexpected end of file")

	in.Close()
	out.Close()
	err = os.Rename(fileName+".tmp", fileName)
	check(err)
}

func main() {
	if len(os.Args) == 1 {
		refresh("README.md")
	} else {
		refresh(os.Args[1])
	}
}
