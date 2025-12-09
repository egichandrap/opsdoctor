package loganalyzer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SlowEntry represents one detected slow request or entry.
type SlowEntry struct {
	Timestamp string `json:"timestamp,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
	Duration  int    `json:"duration_ms"`
	Line      string `json:"line"`
	LineNo    int    `json:"line_no"`
}

// AnalysisResult is the full result of an analyze run.
type AnalysisResult struct {
	Source       string      `json:"source"`
	ThresholdMs  int         `json:"threshold_ms"`
	TotalLines   int         `json:"total_lines"`
	SlowCount    int         `json:"slow_count"`
	SlowExamples []SlowEntry `json:"slow_examples"`
	RunAt        string      `json:"run_at"`
}

// msRegex finds e.g. "123ms"
var msRegex = regexp.MustCompile(`\b([0-9]{1,6})\s*ms\b`)

// traceRegex attempts to find common trace-id patterns (traceId=..., trace-id:..., trace_id=...)
var traceRegex = regexp.MustCompile(`(?i)(?:trace[_\-]?id|tid)[:=]\s*([A-Za-z0-9\-\_:.]+)`)

// Analyze scans a log file at path, returns text report and optionally writes JSON to exportPath.
// thresholdMs: minimal ms to count as slow (e.g. 1000)
// pattern: optional substring pattern to filter lines (if empty, all lines scanned)
// exportPath: if non-empty, JSON result is written to that file
func Analyze(path string, thresholdMs int, pattern string, exportPath string) (AnalysisResult, string, error) {
	res := AnalysisResult{
		Source:      path,
		ThresholdMs: thresholdMs,
		RunAt:       time.Now().Format(time.RFC3339),
	}

	f, err := os.Open(path)
	if err != nil {
		return res, "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNo := 0
	var slowEntries []SlowEntry

	var out strings.Builder
	out.WriteString(fmt.Sprintf("Analyzing %s (threshold=%d ms) ...\n", path, thresholdMs))

	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		// if pattern supplied, skip lines that don't contain it
		if pattern != "" && !strings.Contains(line, pattern) {
			continue
		}

		// search for ms values in the line (may be multiple)
		msMatches := msRegex.FindAllStringSubmatch(line, -1)
		if len(msMatches) > 0 {
			for _, m := range msMatches {
				if len(m) < 2 {
					continue
				}
				msVal, err := strconv.Atoi(m[1])
				if err != nil {
					continue
				}
				if msVal >= thresholdMs {
					// try to extract trace id (optional)
					traceID := ""
					tm := traceRegex.FindStringSubmatch(line)
					if len(tm) >= 2 {
						traceID = tm[1]
					}

					entry := SlowEntry{
						Timestamp: time.Now().Format(time.RFC3339),
						TraceID:   traceID,
						Duration:  msVal,
						Line:      line,
						LineNo:    lineNo,
					}
					slowEntries = append(slowEntries, entry)
					out.WriteString(fmt.Sprintf("Line %d: %dms trace=%s -> %s\n", lineNo, msVal, traceID, summarizeLine(line)))
					break // only count once per line (even if multiple matches)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return res, "", err
	}

	res.TotalLines = lineNo
	res.SlowCount = len(slowEntries)
	// attach up to first 100 slow examples to result
	limit := 100
	if len(slowEntries) < limit {
		limit = len(slowEntries)
	}
	res.SlowExamples = slowEntries[:limit]

	out.WriteString(fmt.Sprintf("Total lines scanned: %d\n", res.TotalLines))
	out.WriteString(fmt.Sprintf("Slow lines found: %d\n", res.SlowCount))

	// If exportPath given -> write JSON
	if exportPath != "" {
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return res, out.String(), fmt.Errorf("failed marshal json: %w", err)
		}
		if err := os.WriteFile(exportPath, b, 0644); err != nil {
			return res, out.String(), fmt.Errorf("failed write export file: %w", err)
		}
		out.WriteString(fmt.Sprintf("\nExported JSON -> %s\n", exportPath))
	}

	return res, out.String(), nil
}

func summarizeLine(s string) string {
	const max = 200
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
