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
	htmlViewInfo = `
	<html>
	<head><title>{{.Id}}</title></head>
	<body>
	<h1>{{.Id}}</h1>
	<table>
	<tr>
		<th>Url</th>
		<th>PathQuery</th>
		<th>MaxClicks</th>
		<th>Clicks</th>
		<th>DateAdded</th>
	</tr>
	<tr>
		<td><a href="{{.Url}}">{{.Url}}</a></td>
		<td>{{.PathQuery}}</td>
		<td>{{.MaxClicks}}</td>
		<td>{{.Clicks}}</td>
		<td>{{.DateAdded}}</td>
	</tr>
	</table>
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
