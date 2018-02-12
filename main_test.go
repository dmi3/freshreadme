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
	_ = os.MkdirAll(filepath.Dir(fn), os.ModePerm)
	var e = ioutil.WriteFile(fn, []byte(content), 0666)
	check(e)
}

func refresh(source string, readme string) string {
	startDir, e := filepath.Abs("./")
	check(e)
	defer os.Chdir(startDir)

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

func TestTheTest(t *testing.T) {
	readme, err := ioutil.ReadFile("examples/README.example.md")
	check(err)

	source1, err := ioutil.ReadFile("examples/Examples.java")
	check(err)

	source2, err := ioutil.ReadFile("examples/examples.py")
	check(err)

	expected, err := ioutil.ReadFile("examples/README.result.md")
	check(err)

	startDir, e := filepath.Abs("./")
	check(e)
	defer os.Chdir(startDir)

	dir, e := ioutil.TempDir("", "example")
	check(e)
	os.Chdir(dir)
	defer os.RemoveAll(dir)

	dump(dir, "README.md", strings.Replace(string(readme), "    ", "", -1))
	dump(dir, "examples/Examples.java", string(source1))
	dump(dir, "examples/examples.py", string(source2))

	main()

	result, err := ioutil.ReadFile("README.md")
	check(err)

	if string(expected) != string(result) {
		t.Errorf("Expected \n---\n%s\n---\nGot\n---\n%s\n---\n", expected, result)
	}
}
