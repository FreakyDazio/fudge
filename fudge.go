package main

import (
	"bufio"
	"encoding/json"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"regexp"
	"errors"
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
	app     *cli.App = cli.NewApp()
	scanner *bufio.Scanner
	input   io.ReadCloser
	parser  LineParser
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

func perform(c *cli.Context) {
	input = os.Stdin
	if len(c.Args()) == 1 {
		var err error
		input, err = os.Open(c.Args()[0])
		handleError(err)
	}
	scanner = bufio.NewScanner(input)
	selectParser(c.String("format"))
	for scanner.Scan() {
		go processLine(scanner.Text())
	}
	input.Close()
}

func init() {
	app.Name = "fudge"
	app.Version = "1.0.0"
	app.Usage = "parse log files like a pro"
	app.Flags = []cli.Flag{
		cli.StringFlag{"format, f", "combined", "preset format of logs"},
	}
	app.Action = perform
}

func main() {
	app.Run(os.Args)
}
