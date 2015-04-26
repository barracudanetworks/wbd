# wbc (Wallboard Control)
A small wallboard program. Acts as a server for wallboards to connect to, and controls what web pages the wallboard computers load.

**NOTE**: This is a work in progress. I am learning Go while writing this app, so please don't expect perfection.

How does it work?
-----------------
Calling `wbc run` will launch a web server on the address and port you specify (0.0.0.0:80 by default).  The web server runs a page containing a full screen iframe and some nifty Javascript so as to allow control over what page the client is viewing.  The Javascript connects back to the Wallboard Control websocket server and listens for commands. The server will tell the client when it's time to load a new webpage.  You can setupa a list of URLs to rotate, shuffle, or stagger (have machines show different pages) through or point a client, or multiple clients, at a specific webpage.

At Barracuda Networks, we use Raspberry Pis hooked up to televisions to drive the wallboards. The wbc server just needs to be run somewhere that the clients can access.
