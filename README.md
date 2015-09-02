# wbd (Wallboard Daemon)
[![Build Status](https://travis-ci.org/barracudanetworks/wbd.svg?branch=master)](https://travis-ci.org/barracudanetworks/wbd)
A wallboard daemon. Runs a full-screened iframe and some Javascript that allows for dynamic updating of wallboards. Useful if you have many clients that you would like to display metrics on.

Requirements
------------
1. Go must be installed
2. `GOPATH` enviornment variable must be set
3. `$GOPATH/bin` must be in `PATH` environment variable

Usage
-----
1. `go get github.com/barracudanetworks/wbd`
2. `wbd install`
3. `wbd run`

If you would like to specify a custom listen address, port, or database location, you may do so with some command-line options (try `wbd help install` or `wbd help run`).

How does it work?
-----------------
Calling `wbd run` will launch a web server on the address and port you specify (`0.0.0.0:80` by default). The web server runs a simple index page, containing a full screened iframe and some nifty Javascript so as to allow control over what page the client is viewing.

The Javascript on the page connects back to the Wallboard Control websocket server and listens for commands. The server will keep the client updated with which URLs it should rotate through. Right now, clients can only rotate through a global pre-defined list of URLs. In the future, you will be able to setup a list of URLs to rotate, shuffle, or stagger (have machines show different pages) for all, or specific, clients.

At Barracuda Networks, we use Raspberry Pis hooked up to televisions to drive the wallboards. The wbd server just needs to be run somewhere that the clients can access.

Command documentation
---------------------
```
NAME:
   wbd - take back control from your televisions

USAGE:
   wbd [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S):
   John Maguire <jmaguire@barracuda.com>

COMMANDS:
   run, r	run the webserver
   url, u	add, remove, or list urls in rotation
   list, l	add, remove, or list url lists
   client, c	alias, remove, or list clients
   assign, a	assign a client or url to a list
   install, i	install the database
   clean	delete the database (WARNING: very destructive)
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version

```
