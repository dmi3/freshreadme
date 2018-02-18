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

func testRefreshSimple(source string, readme string) string {
	return testRefresh("README.md", readme, "source", source)
}

func testRefresh(readmePath string, readme string, sources ...string) string {
	startDir, e := filepath.Abs("./")
	check(e)
	defer os.Chdir(startDir)

	dir, e := ioutil.TempDir("", "example")
	check(e)
	os.Chdir(dir)
	defer os.RemoveAll(dir)

	dump(dir, readmePath, readme)
	for n := 0; n < len(sources); n += 2 {
		dump(dir, sources[n], sources[n+1])
	}

	refresh(readmePath)

	result, err := ioutil.ReadFile(readmePath)
	check(err)
	return string(result)
}

func TestIncludeFile(t *testing.T) {
	result := testRefreshSimple(
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
	result := testRefreshSimple(
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
		testRefreshSimple("",
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
		testRefreshSimple(
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
		testRefreshSimple(
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
	readme, err := ioutil.ReadFile("examples/README.md.example")
	check(err)

	source1, err := ioutil.ReadFile("examples/Examples.java")
	check(err)

	source2, err := ioutil.ReadFile("examples/examples.py")
	check(err)

	expected, err := ioutil.ReadFile("examples/README.md.result")
	check(err)

	result := testRefresh(
		"README.md", strings.Replace(string(readme), "    ", "", -1),
		"examples/Examples.java", string(source1),
		"examples/examples.py", string(source2))

	if string(expected) != result {
		t.Errorf("Expected \n---\n%s\n---\nGot\n---\n%s\n---\n", expected, result)
	}
}

func TestParentAndSubDir(t *testing.T) {
	result := testRefresh(
		"subdir/README.md",
		`before test
<!-- [freshReadmeSource](../source1) -->
`+"```"+`
replaceMe
`+"```"+`
between includes
<!-- [freshReadmeSource](subsub/source2) -->
`+"```"+`
replaceMe
`+"```"+`
after test
`,
		"source1",
		`test 1
test 2
test 3`,
		"subdir/subsub/source2",
		`test 4
test 5
test 6`)

	var expected = `before test
<!-- [freshReadmeSource](../source1) -->
` + "```" + `
test 1
test 2
test 3
` + "```" + `
between includes
<!-- [freshReadmeSource](subsub/source2) -->
` + "```" + `
test 4
test 5
test 6
` + "```" + `
after test
`

	if expected != string(result) {
		t.Errorf("Expected \n---\n%s\n---\nGot\n---\n%s\n---\n", expected, string(result))
	}
}
