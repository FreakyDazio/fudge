package main

import (
  "os"
  "io"
  "os/exec"
  "bufio"
  "regexp"
  "encoding/json"
  "github.com/codegangsta/cli"
)

var(
  app *cli.App
  defaultFormat string
)

func init() {
  app = cli.NewApp()
  app.Name = "fudge"
  app.Usage = "parse and filter log files like a pro"
  app.Action = filterLogFileCommand
  app.Flags = []cli.Flag{
    cli.StringFlag{"format, f", "", "log format (common, s3, combined)"},
    cli.StringFlag{"script, s", "", "filter script"},
  }
}

func regexpMatches(matcher *regexp.Regexp, line string) map[string]string {
  result := make(map[string]string)
  matches := matcher.FindStringSubmatch(line)
  for i, name := range matcher.SubexpNames() {
    if i == 0 {
      continue
    }
    result[name] = matches[i]
  }
  return result
}

func errorMessage(msg string) {
  os.Stderr.WriteString(msg + "\n")
  os.Exit(1)
}

func handleError(e error) {
  if e != nil {
    errorMessage(e.Error())
  }
}

func filterLogFileCommand(c *cli.Context) {
  var matches map[string]string
  var script *exec.Cmd
  var outputStream io.WriteCloser
  var input io.WriteCloser

  // The first argument is the log file
  if len(c.Args()) < 1 {
    errorMessage("requires file path as first argument")
  }

  // Compile the provided regexp
  regexp, err := regexp.Compile(c.String("format"))
  handleError(err)

  // Open the file for reading
  file, err := os.Open(c.Args()[0])
  handleError(err)

  // Locat the optional filtering script
  scriptPath := c.String("script")
  if scriptPath != "" {
    script = exec.Command(scriptPath)
    // Instead of writting to stdout we want ot write to the scripts stdin
    outputStream, err = script.StdinPipe()
    handleError(err)
    script.Stdout = os.Stdout
    err = script.Start()
    handleError(err)
  } else {
    outputStream = os.Stdout
  }

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    matches = regexpMatches(regexp, scanner.Text())
    rawJson, err := json.Marshal(matches)
    handleError(err)

    _, err = outputStream.Write(rawJson)
    handleError(err)
    _, err = outputStream.Write([]byte("\n"))
    handleError(err)
  }
  if scriptPath != "" {
    _, err = outputStream.Write([]byte("\n"))
    handleError(err)
  }
  err = outputStream.Close()
  handleError(err)
  err = file.Close()
  handleError(err)
}

func main() {
  app.Run(os.Args)
}
