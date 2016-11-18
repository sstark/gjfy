package main

const (
	htmlMaster = `
	{{define "master"}}
	<html>
	<head>
		<title>GJFY - {{template "title" .}}</title>
		<link rel="shortcut icon" type="image/x-icon" href="favicon.ico" />
	</head>
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

var (
	favicon = []byte{0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x10, 0x10, 0x10, 0x0, 0x1,
		0x0, 0x4, 0x0, 0x28, 0x1, 0x0, 0x0, 0x16, 0x0, 0x0, 0x0, 0x28, 0x0, 0x0,
		0x0, 0x10, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x1, 0x0, 0x4, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff,
		0xff, 0x0, 0x0, 0xfe, 0x7f, 0x0, 0x0, 0xfe, 0x7f, 0x0, 0x0, 0xff, 0xff,
		0x0, 0x0, 0xfe, 0x7f, 0x0, 0x0, 0xfe, 0x7f, 0x0, 0x0, 0xff, 0x3f, 0x0, 0x0,
		0xff, 0x9f, 0x0, 0x0, 0xfd, 0xdf, 0x0, 0x0, 0xfc, 0x9f, 0x0, 0x0, 0xfe,
		0x3f, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff,
		0x0, 0x0}
)