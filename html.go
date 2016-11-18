package main

const (
	htmlMaster = `
	{{define "master"}}
	<html>
	<head><title>GJFY - {{template "title" .}}</title></head>
	<body>
	{{template "content" .}}
	</body>
	</html>
	{{end}}
	`
	htmlView = `
	{{define "title"}}VIEW{{end}}
	{{define "content"}}
	<h1>{{.Id}}</h1>
	<div>The secret is:</div>
	<div>{{.Secret}}</div>
	{{end}}
	`
	htmlViewInfo = `
	{{define "title"}}VIEWINFO{{end}}
	{{define "content"}}
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
	{{end}}
	`
	htmlViewErr = `
	{{define "title"}}ERROR{{end}}
	{{define "content"}}
	<h1>Not available</h1>
	<p>This ID is not valid anymore. Please request another one from the person who sent you this link.</p>
	{{end}}
	`
)
