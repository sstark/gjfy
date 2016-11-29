gjfy - one time password link server
====================================

(Transport encryption is not implemented yet, please do not use in production!)

What does it do?
----------------

gjfy is a single binary web server with only one purpose: Create links that
automatically disappear once clicked. On first click it will show a "secret",
for instance a password that somebody wants to send to someone.

The idea is that if the original receiver finds the link invalid, they know
that the secret was intercepted by a third party and the sender can reset the
password. This does not protect against eavesdropping attacks, for this you
need a TLS connection.

There is no persistency: If the server process ends, all secrets are gone.

Please be careful: Using a tool like gjfy is only advised when all other
options are even less secure (mail, e-mail, phone). In any case, if you send a
password, the receiver should be told to change it as soon as possible.

Building
--------

Download the code and run "go build", it will create a single statically linked
binary for easy deployment.

`go get github.com/sstark/gjfy` will also work.

Running
-------

Choose the IP address and port gjfy listens on with the `-listen` parameter.

Examples:

    gjfy -listen '0.0.0.0:1234'    # listen on all IPv4 addresses
    gjfy -listen '[::1]:4123'      # listen on localhost, IPv6 only
    gjfy -listen ':6234'           # listen on all addresses, IPv4 and IPv6

To tell gjfy its name as seen by users of the service, use the `-urlbase` parameter like so:

    gjfy -urlbase 'https://gjfy.example.org'
    gjfy -urlbase 'https://gjfy.example.org:4123'

Use `gjfy -help` for help.

Options
-------

Custom CSS styling can by applied by placing a file "custom.css" in either
`/etc/gjfy/custom.css` or `$PWD/custom.css`.

An authentication token database should placed in either `/etc/gjfy/auth.db` or
`$PWD/auth.db`. An example file is distributed with the software. New secrets
can only be created with a valid auth token in the POST request.

`$PWD/<file>` will take precedence over `/etc/gjfy/<file>` for above options.

To trigger reloading of auth.db or custom.css you can send SIGHUP to the gjfy
process.

Usage
-----

Currently the only way to create new secrets is by using a json API. An example
client (gjfy-post) is included. Basically a request looks like this:

    {"auth_token":"g4uhg3iu4h5i3u4","secret":"someSecret"}

By sending this to `/api/v1/new` you create a new URL which is a hash over that
json structure. The reply from the server will tell you this link in both, a
user friendly version and in an API version. Invocation of that link will
immediately lead to deletion of the secret in the server. However, there is an
exception: you can post a `"max_clicks":n` variable along with the json and it
will allow up to `n` clicks.

The authentication token sent with the request will not be stored in the
server. Instead, the associated email address will be stored with the secret.
Future versions of the server would allow for email notification in case the
link is clicked.

A timeout can be set by including `"valid_for:n"` in the request. The secret
will become invalid after n days, even if not clicked.

gjfy-post
---------

gjfy-post is a demonstration client using zsh, curl and jq.

    usage: ./gjfy-post <authtoken> <secret> [maxclicks]

Required arguments are authtoken and the secret itself. Please note that
providing the secret this way makes it readable in the system process listing!

Upcoming features
-----------------

  - unit tests
  - TLS support. You can use a reverse proxy for now
  - Mail notification to sender
  - Memory protection in server binary
  - customizing of HTML
