package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assertPanic(t *testing.T, msg string, f func()) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			if !strings.Contains(r.(string), msg) {
				t.Errorf("Expected `%s` error, got `%s`", msg, r)
			}
		}
	}()
	f()
}

func dump(dir string, fileName string, content string) {
	var fn = filepath.Join(dir, fileName)
	var e = ioutil.WriteFile(fn, []byte(content), 0666)
	check(e)
}

func refresh(source string, readme string) string {
	dir, e := ioutil.TempDir("", "example")
	check(e)
	os.Chdir(dir)
	defer os.RemoveAll(dir)

	dump(dir, "README.md", readme)
	dump(dir, "source", source)

	main()

	result, err := ioutil.ReadFile("README.md")
	check(err)
	return string(result)
}

func TestIncludeFile(t *testing.T) {
	result := refresh(
		`test 1
test 2
test 3`,

		`
before test
<!-- [freshReadmeSource](source) -->
`+"```"+`
replaceMe
`+"```"+`
after test
`)

	var expected = `
before test
<!-- [freshReadmeSource](source) -->
` + "```" + `
test 1
test 2
test 3
` + "```" + `
after test
`

	if expected != string(result) {
		t.Errorf("Expected \n---\n%s\n---\nGot\n---\n%s\n---\n", expected, string(result))
	}
}

func TestIncludeSnippet(t *testing.T) {
	result := refresh(
		`
before test
# freshReadmeSnippet: header
test 1
test 2
test 3
# freshReadmeSnippet: header
after test
`,

		`
before test
<!-- [freshReadmeSource](source#header) -->
`+"```"+`
replaceMe
`+"```"+`
after test
`)

	var expected = `
before test
<!-- [freshReadmeSource](source#header) -->
` + "```" + `
test 1
test 2
test 3
` + "```" + `
after test
`

	if expected != string(result) {
		t.Errorf("Expected \n---\n%s\n---\nGot\n---\n%s\n---\n", expected, string(result))
	}
}

func TestEmptyInclude(t *testing.T) {
	assertPanic(t, "Unable to find snippet", func() {
		refresh("",
			`
before test
<!-- [freshReadmeSource](source#header)  -->
`+"```"+`
replaceMe
`+"```"+`
after test
`)

	})
}

func TestIncludeMultipleTimes(t *testing.T) {
	assertPanic(t, "appears second time", func() {
		refresh(
			`
before test
# freshReadmeSnippet: header
test 1
# freshReadmeSnippet: header
repeat
# freshReadmeSnippet: header
test 1
# freshReadmeSnippet: header
after test
`,
			`
before test
<!-- [freshReadmeSource](source#header) -->
`+"```"+`
replaceMe
`+"```"+`
after test
`)

	})
}

func TestEOF(t *testing.T) {
	assertPanic(t, "Unexpected end of file", func() {
		refresh(
			`
before test
# freshReadmeSnippet: header
test 1
# freshReadmeSnippet: header
after test
`,
			`
before test
<!-- [freshReadmeSource](source#header) -->
`+"```"+`
replaceMe
after test
`)

	})
}
