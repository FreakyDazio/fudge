# Fudge

Process NGINX combined logs quickly into JSON control what fields are provided... Simple really.

### Usage

Process input from STDIN

    $ fudge access.log
    $ fudge < access.log

Select what fields you wish to view

    $ fudge -p status access.log
    $ fudge -p httpReferer,httpUserAgent access.log

Be all unixy and stuff by combining it with other commands

    $ fudge -p httpUserAgent access.log | grep google | wc -l # google user agent count
    $ fudge access.log | json -d -c "/pingdom/i.test(this.httpUserAgent)" # view pingdom requests

### Disclaimer

I developed this tool to scratch an itch and to experiment with streams in node.
I have no idea if it will work for everyone and I make no promises it will.
