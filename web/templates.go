package web

var indexTemplate string = `
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

	var urls = [
		{{ range .URLs }}'{{ . }}',
		{{ else }}'/welcome?client={{ .Client }}',
		{{ end }}
	];

	document.addEventListener("DOMContentLoaded", function(event) {
		var rotate = new SiteRotator('frame', urls, 60);
	});
	</script>
</head>
<body>
	<iframe id='frame'>Oops, something went wrong with the Wallboard page!</iframe>
</body>
</html>
`

var welcomeTemplate string = `
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
		position: relative;
		top: 50%;
		text-align: center;
		transform: translateY(-50%);
	}
	h1 {
		font-size: 6em;
	}
	</style>
</head>
<body>
	<div class='wrapper'>
		{{ if ne .Client "" }}<h1>Client: {{ .Client }}</h1>{{end}}
		<h1>IP Addr: {{ .RemoteAddr }}</h1>
	</div>
</body>
</html>
`
