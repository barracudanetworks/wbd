package web

const (
	indexTemplate string = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title> Wallboard Control </title>

	<style type='text/css'>
	/* Remove padding around iframe */
	html, body {
		margin: 0;
		height: 100%;
		overflow: hidden;
	}

	iframe {
		border: 0;
		width: 100%;
		height: 100%;
	}
	</style>

	<script type='text/javascript'>
	Array.prototype.equals = function (array) {
		// if the other array is a falsy value, return
		if (!array)
			return false;

		// compare lengths - can save a lot of time
		if (this.length != array.length)
			return false;

		for (var i = 0, l=this.length; i < l; i++) {
			// Check if we have nested arrays
			if (this[i] instanceof Array && array[i] instanceof Array) {
				// recurse into the nested arrays
				if (!this[i].equals(array[i]))
					return false;
			}
			else if (this[i] != array[i]) {
				// Warning - two different object instances will never be equal: {x:20} != {x:20}
				return false;
			}
		}
		return true;
	}

	function SiteRotator (elementId, defaultUrls, duration) {
		if (typeof elementId === 'undefined') return;
		if (typeof defaultUrls === 'undefined' || defaultUrls.length < 1) return;

		this.elementId = elementId;
		this.defaultUrls = defaultUrls;
		this.duration = duration;

		this.urls = defaultUrls;
		this.currentIndex = 0;

		this.init = function() {
			// Load first URL when initialized
			console.log("Initializing rotator");

			// Try to use the default URLs if we don't have any
			if (this.urls.length < 1) {
				if (this.defaultUrls.length < 1) {
					if (typeof this.interval !== 'undefined') {
						clearInterval(this.interval);
					}

					console.error("Can't run rotator -- no URLs to rotate");
					return;
				}

				this.urls = this.defaultUrls;
			}

			this.currentIndex = 0;
			this.load(this.urls[this.currentIndex]);

			// If a duration is passed in, setup rotation
			if (typeof this.duration !== 'undefined')
			{
				this.rotateEvery(this.duration);
			}
		};

		this.setUrls = function(urls) {
			if (urls === null || urls === undefined) {
				urls = [];
			}

			// Only update if URLs changed -- this reinits the rotator
			if (this.urls.equals(urls) === false) {
				console.log("Current URLs:", this.urls);
				console.log("Updated URLs:", urls);

				this.urls = urls;
				this.init();
			}
		};

		this.load = function(url) {
			console.log("Loading URL:", url)

			document.getElementById(this.elementId).src = url;
		};

		this.next = function() {
			console.log("Moving to next URL");

			if (this.urls.length < 1) {
				return;
			}

			this.currentIndex++;
			if (this.currentIndex >= this.urls.length) {
				this.currentIndex = 0;
			}

			this.load(this.urls[this.currentIndex]);
		};

		this.previous = function() {
			if (this.urls.length < 1) {
				return;
			}

			this.currentIndex--;
			if (this.currentIndex < 0) {
				this.currentIndex = this.urls.length - 1;
			}

			this.load(this.urls[this.currentIndex]);
		};

		this.pause = function(duration) {
			// If we're already paused, bail
			if (typeof this._load !== 'undefined') { return; }

			// Save this.load and noop it
			this._load = this.load;
			this.load = function(url) { }

			setInterval(function() {
				this.resume();
			}, duration * 1000);
		}

		this.resume = function() {
			if (typeof this._load === 'undefined') { return; }

			// Return the function
			this.load = this._load;

			// Load the page we should be on
			this.load(this.urls[this.currentIndex]);
		}

		this.rotateEvery = function(duration) {
			if (typeof this.interval !== 'undefined') {
				clearInterval(this.interval);
			}

			var t = this;
			this.interval = setInterval(function() {
				t.next();
			}, duration * 1000);

			console.log("Rotate scheduled for every ", duration, "seconds");
		};

		this.init();
	}

	function wbcConnect(endpoint, rotator) {
		if (typeof endpoint === 'undefined') return false;
		if (!window["WebSocket"]) return false;

		conn = new WebSocket(endpoint);
		conn.onopen = function(evt) {
			console.log("Connected to websocket server");
			conn.send(JSON.stringify({
				"action": "sendUrls"
			}));
		}
		conn.onclose = function(evt) {
			console.log("Disconnected from websocket server");
		}
		conn.onmessage = function(evt) {
			message = JSON.parse(evt.data);

			if (typeof message.action === 'undefined')
			{
				console.error("No action in message from server:", message)
				return;
			}

			switch (message.action) {
			case 'updateUrls':
				rotator.setUrls(message.data.urls);

				break;
			case 'flashUrl':
				rotator.pause(message.data.url);
				rotator.load(message.data.duration);

				break;
			default:
				console.error("Unknown action in message from server:", message)
				break;
			}
		}
	}

	document.addEventListener("DOMContentLoaded", function(event) {
		var defaultUrls = [
			'/welcome?client={{ .Client }}'
		];

		var rotator = new SiteRotator('frame', defaultUrls, 60);

		{{ if ne .Client "" }}
		// Connect to WebSocket server (provides control)
		wbcConnect("ws://{{.Address}}/ws?client={{ .Client }}", rotator);
		{{ else }}
		wbcConnect("ws://{{.Address}}/ws", rotator);
		{{ end }}
	});
	</script>
</head>
<body>
	<iframe id='frame'>Oops, something went wrong with the Wallboard page!</iframe>
</body>
</html>
`

	welcomeTemplate string = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title> Welcome </title>
	<style type='text/css'>
	html, body {
		height: 100%;
		width: 100%;
		margin: 0;
		padding: 0;
	}
	div.wrapper {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		-webkit-transform: translate(-50%, -50%);
		-moz-transform: translate(-50%, -50%);
		-ms-transform: translate(-50%, -50%);
	}
	h1 {
		font-size: 6em;
	}
	</style>
</head>
<body>
	<div class='wrapper'>
		<h1>wbc</h1>
		{{ if ne .Client "" }}<h2>Client: {{ .Client }}</h2>{{end}}
		<h2>IP Addr: {{ .RemoteAddr }}</h2>
		<p>Add a URL or two and this page will disappear. :)</p>
	</div>
</body>
</html>
`
)
