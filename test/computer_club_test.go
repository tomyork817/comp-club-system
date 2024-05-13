package test

import (
	"comp-club-system/internal/parser"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

const (
	inputFile = `input\d+.txt`
	dataPath  = "testdata"
)

func TestComputerClub(t *testing.T) {
	entries, err := os.ReadDir(dataPath)
	if err != nil {
		t.Fatalf("error reading dir: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if matched, _ := regexp.MatchString(inputFile, entry.Name()); matched {
			num := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "input"), ".txt")
			n, err := strconv.Atoi(num)
			if err != nil {
				t.Fatalf("error reading name of input file: %v", err)
			}
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			c, err := parser.ReadComputerClub(fmt.Sprintf("%s/%s", dataPath, entry.Name()))
			if err == nil {
				c.RunIncomingEvents()
				c.Print()
			} else {
				fmt.Println(err.Error())
			}
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = old
			expected, err := os.ReadFile(fmt.Sprintf("%s/output%d.txt", dataPath, n))
			if err != nil {
				t.Fatalf("error reading output: %v", err)
			}
			if strings.TrimSpace(string(out)) != strings.TrimSpace(string(expected)) {
				t.Errorf("Expected:\n%s\nGot:\n%s", expected, out)
			}
		}
	}
}
