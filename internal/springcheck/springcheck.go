package springcheck

import (
	"bufio"
	"os"
	"strings"
)

type Result struct {
	SlowQueries []string
	Errors      []string
}

func Analyze(path string) (*Result, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	res := &Result{}
	sc := bufio.NewScanner(file)

	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, "ERROR") {
			res.Errors = append(res.Errors, line)
		}
		if strings.Contains(line, "took") && strings.Contains(line, "ms") {
			res.SlowQueries = append(res.SlowQueries, line)
		}
	}
	return res, nil
}
