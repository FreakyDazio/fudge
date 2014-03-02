# Fudge

Process logs quickly into JSON making it more convenient to pipe into your own scripts to filter

### Usage

fudge can accept inputs in a variety of ways.

1. stdin - `cat /var/logs/nginx/* | fudge`
2. arguments - `fudge /var/logs/nginx/access.log`
3. globs - `fudge /var/logs/nginx/*.log`
4. multiples - `fudge /var/logs/nginx/{error,access}*.log /src/example.org/logs/*`

fudge has a couple of options available:

- `--format, -f` switch to a different log format parser (default: nginx "common")
- `--gzip, -g` decompress gzip input files on the fly (like gzcat)

## Awesome combinations

My primary use case is to filter logs based on certain information. A good tool to use with fudge is [jq][0] as it allows you to filter and adjust JSON on the command line.

[0]: https://github.com/stedolan/jq

## Contributions

Currently I am looking for common log file formats. I still need to add Apache and many other types of log files into fudge to make it more useful to more people.

Happily accepting pull requests
