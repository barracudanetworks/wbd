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
	function SiteRotator (elementId, urls, duration) {
		if (typeof elementId === 'undefined') return;
		if (typeof urls === 'undefined' || urls.length < 1) return;

		this.elementId = elementId;
		this.urls = urls;
		this.currentIndex = 0;

		this.init = function(duration) {
			// Load first URL when initialized
			console.log("Initial load")

			this.load(urls[this.currentIndex]);

			// If a duration is passed in, setup rotation
			if (typeof duration !== 'undefined')
			{
				this.rotateEvery(duration * 1000);
			}
		};

		this.load = function(url) {
			console.log("Loading URL:", url)

			document.getElementById(this.elementId).src = url;
		};

		this.next = function() {
			console.log("Moving to next URL");

			if (urls.length < 1) {
				return;
			}

			this.currentIndex++;
			if (this.currentIndex >= urls.length) {
				this.currentIndex = 0;
			}

			this.load(urls[this.currentIndex]);
		};

		this.previous = function() {
			if (urls.length < 1) {
				return;
			}

			this.currentIndex--;
			if (this.currentIndex < 0) {
				this.currentIndex = this.urls.length - 1;
			}

			this.load(urls[this.currentIndex]);
		};

		this.rotateEvery = function(duration) {
			console.log("Setting up rotation at every", duration / 1000, "seconds")
			setInterval(this.next.bind(this), duration);
		};

		this.init(duration);
	}

	function wbcConnect(endpoint) {
		if (typeof endpoint === 'undefined') return false;
		if (!window["WebSocket"]) return false;

		conn = new WebSocket(endpoint);
		conn.onopen = function(evt) {
			console.log("Connected to websocket server");
			conn.send("Hello");
		}
		conn.onclose = function(evt) {
			console.log("Disconnected from websocket server");
		}
		conn.onmessage = function(evt) {
			console.log("Message from websocket server: ", evt.data);
		}
	}

	var urls = [
		{{ range .URLs }}'{{ . }}',
		{{ else }}'/welcome?client={{ .Client }}',
		{{ end }}
	];

	document.addEventListener("DOMContentLoaded", function(event) {
		var rotate = new SiteRotator('frame', urls, 60);
		{{ if ne .Client "" }}
		// Connect to WebSocket server (provides control)
		wbcConnect("ws://{{.Address}}/ws?client={{ .Client }}");
		{{ else }}
		wbcConnect("ws://{{.Address}}/ws");
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
