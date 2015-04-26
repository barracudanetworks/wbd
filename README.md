# wbc (Wallboard Control)
A small wallboard program. Acts as a server for wallboards to connect to, and controls what web pages the wallboard computers load.

**NOTE**: This is a work in progress. I am learning Go while writing this app, so please don't expect perfection.

Requirements
------------
1. Go must be installed
2. `GOPATH` enviornment variable must be set
3. `$GOPATH/bin` must be in `PATH` environment variable

Usage
-----
1. `go get github.com/johnmaguire/wbc`
2. `wbc install`
3. `wbc run`

How does it work?
-----------------
Calling `wbc run` will launch a web server on the address and port you specify (`0.0.0.0:80` by default).  The web server runs a simple index page, containing a full screened iframe and some nifty Javascript so as to allow control over what page the client is viewing.

The Javascript on the page connects back to the Wallboard Control websocket server and listens for commands. The server will tell the client when it's time to load a new webpage.  You can setup a list of URLs to rotate, shuffle, or stagger (have machines show different pages) for all connected clients, or specify specific URLs for each client.

At Barracuda Networks, we use Raspberry Pis hooked up to televisions to drive the wallboards. The wbc server just needs to be run somewhere that the clients can access.

Command documentation
---------------------
```
NAME:
   wbc - take back control from your televisions

USAGE:
   wbc [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S):
   John Maguire <contact@johnmaguire.me>

COMMANDS:
   run, r       run the webserver
   url, u       add, remove, or list urls in rotation
   install, i   install the database
   clean        delete the database (WARNING: very destructive)
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h       show help
   --version, -v    print the version
```
