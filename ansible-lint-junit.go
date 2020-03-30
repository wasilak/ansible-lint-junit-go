package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

// AppVersion number
const AppVersion = "0.13"

// Failure struct
type Failure struct {
	File    string `xml:"file,attr"`
	Line    string `xml:"line,attr"`
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// TestCase struct
type TestCase struct {
	Name    string    `xml:"name,attr"`
	Failure []Failure `xml:"failure"`
}

// TestSuite struct
type TestSuite struct {
	TestCase []TestCase `xml:"testcase"`
	Errors   int        `xml:"errors,attr"`
	Failures int        `xml:"failures,attr"`
	Tests    int        `xml:"tests,attr"`
	Time     int        `xml:"time,attr"`
}

// TestSuites struct
type TestSuites struct {
	XMLName   xml.Name    `xml:"testsuites"`
	TestSuite []TestSuite `xml:"testsuite"`
}

func appVersion() string {
	return AppVersion
}

func getInput() string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	input := ""

	if input = flag.Arg(0); input != "" {
		content, err := ioutil.ReadFile(input)
		if err != nil {
			log.Fatal(err)
		}

		input = string(content)
	} else {
		if fi.Mode()&os.ModeNamedPipe != 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				input = input + scanner.Text() + "\n"
			}
		} else {
			fmt.Println("Please provide either path to file as [input] or use pipe (stdin).")
			flag.Usage()
			os.Exit(1)
		}
	}

	return input
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func createXML(input string) *TestSuites {
	lines := deleteEmpty(strings.Split(input, "\n"))

	testSuites := &TestSuites{}
	testSuite := TestSuite{
		Errors:   len(lines),
		Failures: 0,
		Tests:    len(lines),
		Time:     0,
	}

	compRegEx := regexp.MustCompile(`^(?P<file>.*?):(?P<line>\d+?):\s\[(?P<code>.*)\]\s(?P<message>.*)$`)

	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		match := compRegEx.FindStringSubmatch(line)
		results := map[string]string{}
		for i, name := range compRegEx.SubexpNames() {
			if i > 0 && i <= len(match) {
				results[name] = match[i]
			}
		}

		text := "[" + results["code"] + "] " + results["message"]

		failure := Failure{
			File:    results["file"],
			Line:    results["line"],
			Message: text,
			Type:    "Ansible Lint",
			Content: text,
		}

		testCase := TestCase{Name: results["file"]}
		testCase.Failure = append(testCase.Failure, failure)
		testSuite.TestCase = append(testSuite.TestCase, testCase)

	}

	testSuites.TestSuite = append(testSuites.TestSuite, testSuite)

	return testSuites
}

func consoleLog(xmlContent *TestSuites, wg *sync.WaitGroup) {
	defer wg.Done()
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "    ")
	if err := enc.Encode(xmlContent); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func writeToFile(xmlContent *TestSuites, filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, _ := os.Create(filename)

	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "    ")
	if err := enc.Encode(xmlContent); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [input]\n", os.Args[0])
		flag.PrintDefaults()
	}

	version := flag.Bool("version", false, "Prints current version")
	verbose := flag.Bool("verbose", false, "Print XML to console as command output")
	output := flag.String("output", "ansible-lint-junit.xml", "Output XML to [output] file")
	flag.Parse()

	if *version {
		fmt.Println(appVersion())
		os.Exit(0)
	}

	var wg sync.WaitGroup

	input := make(chan string)

	go func() {
		input <- getInput()
	}()

	xmlContent := createXML(<-input)

	if *verbose {
		wg.Add(1)
		go consoleLog(xmlContent, &wg)
	}

	wg.Add(1)
	go writeToFile(xmlContent, *output, &wg)

	wg.Wait()
}
