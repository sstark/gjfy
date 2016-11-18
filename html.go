package main

const (
	htmlView = `
	<html>
	<head><title>{{.Id}}</title></head>
	<body>
	<h1>{{.Id}}</h1>
	<div>The secret is:</div>
	<div>{{.Secret}}</div>
	</body>
	</html>
	`
	htmlViewErr = `
	<html>
	<head><title>Error</title></head>
	<body>
	<h1>Not available</h1>
	<p>This ID is not valid anymore. Please request another one from the person who sent you this link.</p>
	</body>
	</html>
	`
)
