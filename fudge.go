package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type Line struct {
	RemoteAddr           string `json:"remote_addr,omitempty"`
	RemoteUser           string `json:"remote_user,omitempty"`
	TimeLocal            string `json:"time_local,omitempty"`
	Request              string `json:"request,omitempty"`
	Status               string `json:"status,omitempty"`
	BodyBytesSent        string `json:"body_bytes_sent,omitempty"`
	HTTPReferer          string `json:"http_referer,omitempty"`
	HTTPUserAgent        string `json:"http_user_agent,omitempty"`
	HTTPXForwardedFor    string `json:"http_x_forwarded_for,omitempty"`
	Connection           string `json:"connection,omitempty"`
	ConnectionRequests   string `json:"connection_requests,omitempty"`
	MSec                 string `json:"msec,omitempty"`
	Pipe                 string `json:"pipe,omitempty"`
	RequestLength        string `json:"request_length,omitempty"`
	RequestTime          string `json:"request_time,omitempty"`
	TimeISO8601          string `json:"time_iso8601,omitempty"`
	UpstreamResponseTime string `json:"upstream_response_time,omitempty"`
}

type LineParser interface {
	Parse(string) (*Line, error)
}

type CombinedParser struct {
	Matcher *regexp.Regexp
}

var (
	app    *cli.App = cli.NewApp()
	parser LineParser
)

const UnparsableLine string = "failed to parse line"

func (parser *CombinedParser) Parse(line string) (*Line, error) {
	lineStruct := &Line{}
	matches := parser.Matcher.FindStringSubmatch(line)
	if len(matches) < 9 {
		return lineStruct, errors.New(UnparsableLine)
	}
	lineStruct.RemoteAddr = matches[1]
	lineStruct.RemoteUser = matches[2]
	lineStruct.TimeLocal = matches[3]
	lineStruct.Request = matches[4]
	lineStruct.Status = matches[5]
	lineStruct.BodyBytesSent = matches[6]
	lineStruct.HTTPReferer = matches[7]
	lineStruct.HTTPUserAgent = matches[8]
	return lineStruct, nil
}

func displayError(message string) {
	os.Stderr.Write([]byte(message + "\n"))
}

func exitWithMessage(message string) {
	displayError(message)
	os.Exit(1)
}

func handleError(err error) {
	if err != nil {
		exitWithMessage(err.Error())
	}
}

func selectParser(format string) {
	var regex *regexp.Regexp
	switch format {
	case "combined":
		regex = regexp.MustCompile(
			`(.+) - (.+) \[(.+)\] "(.+)" (.+) (.+) "(.+)" "(.+)"`,
		)
		parser = &CombinedParser{regex}
	default:
		exitWithMessage("unknown format provided")
	}
}

func processLine(line string) {
	lineStruct, err := parser.Parse(line)
	if err != nil {
		displayError(err.Error())
		return
	}
	output, err := json.Marshal(lineStruct)
	handleError(err)
	os.Stdout.Write(output)
	os.Stdout.Write([]byte("\n"))
}

func processInput(input io.ReadCloser, gzipped bool) {
	var scanner *bufio.Scanner
	if gzipped {
		gzipInput, err := gzip.NewReader(input)
		handleError(err)
		scanner = bufio.NewScanner(gzipInput)
	} else {
		scanner = bufio.NewScanner(input)
	}
	for scanner.Scan() {
		processLine(scanner.Text())
	}
	input.Close()
}

func listInputs(pattern string) []io.ReadCloser {
	inputs := make([]io.ReadCloser, 0)
	paths, err := filepath.Glob(pattern)
	if err != nil {
		displayError(err.Error())
		return inputs
	}
	for _, path := range paths {
		input, err := os.Open(path)
		if err == nil {
			inputs = append(inputs, input)
		} else {
			displayError(err.Error())
		}
	}
	return inputs
}

func perform(c *cli.Context) {
	selectParser(c.String("format"))
	inputs := make([]io.ReadCloser, 0)
	if len(c.Args()) == 0 {
		inputs = append(inputs, os.Stdin)
	} else {
		for _, pattern := range c.Args() {
			inputs = append(inputs, listInputs(pattern)...)
		}
	}
	for _, input := range inputs {
		processInput(input, c.Bool("gzip"))
	}
}

func init() {
	app.Name = "fudge"
	app.Version = "1.0.0"
	app.Usage = "parse log files like a pro"
	app.Flags = []cli.Flag{
		cli.StringFlag{"format, f", "combined", "preset format of logs"},
		cli.BoolFlag{"gzip, g", "Decompress logs on the fly"},
	}
	app.Action = perform
}

func main() {
	app.Run(os.Args)
}
