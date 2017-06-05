package main

import (
	"io"
	"text/template"
)

type ClientVars struct {
	DefaultPostURL string
}

const (
	shellClient = `
#!/bin/bash

POSTURL="${GJFY_POSTURL:-{{.DefaultPostURL}}}"

which jq >/dev/null 2>&1 || {
    echo "jq utility not found" >&2
    exit 1
}

if [[ $# -lt 2 ]]
then
    echo "usage: $0 <authtoken> <secret>" >&2
    exit 2
fi

if [[ -z "$3" ]]
then
    postdata="{\"auth_token\":\"$1\",\"secret\":\"$2\"}"
else
    postdata="{\"auth_token\":\"$1\",\"secret\":\"$2\",\"max_clicks\":$3}"
fi
curl -s -X POST -d "$postdata" "$POSTURL" | jq -r '.url,.api_url,.error | select (.!=null)'
`
)

func ClientShellScript(out io.Writer, url string) error {
	cv := ClientVars{url}
	tmpl, err := template.New("shellClient").Parse(shellClient)
	if err != nil {
		return err
	}
	err = tmpl.Execute(out, cv)
	if err != nil {
		return err
	}
	return err
}
