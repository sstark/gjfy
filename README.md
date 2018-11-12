![gjfy](https://raw.githubusercontent.com/sstark/gjfy/master/logo.png)

one time link server
====================

What does it do?
----------------

gjfy is a single binary, standalone web server with only one purpose: Create
links that automatically disappear once clicked. On first click it will show a
"secret", for instance a password that somebody wants to send to someone.

The idea is that if the original receiver finds the link invalid, they know
that the secret was intercepted by a third party and the sender can reset the
password. This does not protect against eavesdropping attacks, for this you
need a TLS connection.

There is no persistency: If the server process ends, all secrets are gone.

Please be careful: Using a tool like gjfy is only advised when all other
options are even less secure (mail, e-mail, phone). In any case, if you send a
password, the receiver should be told to change it as soon as possible.

What makes it different?
------------------------

There are other tools available that do similar things. However, usually those
involve installing lots of dependencies or web frameworks and often require
setting up a database. Some of them are even offering a hosted service, so
you would be handing your secrets to a third party.

Gjfy does not need any of this: it is a completely self-contained and
on-premise system.

Probably the most notable difference is that secrets are only kept in memory.
They are never written into a database or a file. So it can never happen that,
because of a program bug or sysadmin mistake, the secrets are left on the disk.
However, it is possible that the operating system will write part of the
program memory into swap temporarily, which is not easy to avoid.

The author believes that tools like this should not load assets from external
sources and also that no javascript should be used. Gjfy will never do that and
instead try to be as simple and privacy respecting as possible.


Features
--------

  - Everything in a single binary
  - No web server or application server needed
  - No database needed
  - No persistence
  - No javascript
  - Simple json API (demo client included)
  - Simple html user interface
  - The CSS styling, logo and user message can be customised
  - Simple token based authentication
  - Supports IPv6, HTTP2, TLS
  - Email notification

Building
--------

A precompiled binary is provided with each release. It is also easy to build
gjfy yourself, in case you prefer that:

If you do not have a go environment installed already, install it from your
linux distribution repository (e. g. `apt-get install golang-go`) or download
it from the [go home page](https://golang.org/dl/).

Download the code and run "go build", it will create a single binary file for
easy deployment.

`go get github.com/sstark/gjfy` will also work.

Installation
------------

Create a directory, e. g. `/usr/local/gjfy`. Then copy the following files to it:

  - gjfy (the binary you just built)<sup>1</sup>
  - auth.db
  - logo.png
  - custom.css

For integration into the various system management environments like upstart or
systemd, check the init/ subdirectory for examples.

<sup>1</sup>If you installed using "go get" the binary will be located at
`$GOPATH/bin/gjfy`, while the rest of the files will be under `$GOPATH/src/github.com/sstark/gjfy`

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

To use TLS security add the `-tls` switch:

    gjfy -tls

The scheme will automatically switch to https unless you set urlbase. Before
you can turn on tls you must create a certificate file called `gjfy.crt` and a
key file called `gjfy.key`.

Use `gjfy -help` for help.

Options
-------

Custom CSS styling can by applied by placing a file "custom.css" in either
`/etc/gjfy/custom.css` or `$PWD/custom.css`.

An authentication token database should placed in either `/etc/gjfy/auth.db` or
`$PWD/auth.db`. An example file is distributed with the software. New secrets
can only be created with a valid auth token in the POST request.

If you are using TLS mode you need to put in place either `/etc/gjfy/gjfy.crt`
or `$PWD/gjfy.crt`. Same applies to the key file `gjfy.key`.

The logo.png can be replaced by a custom logo if needed. (It must be png)

You may create a file `userMessageView.txt` that will contain the message the
user sees when clicking on the link. It will replace the default message. HTML
can not be used.

`$PWD/<file>` will take precedence over `/etc/gjfy/<file>` for above options.

To trigger reloading of auth.db, logo.png, custom.css or userMessageView.txt
you can send SIGHUP to the gjfy process. The TLS certificate or key won't be
reloaded this way.

Authentication
--------------

gjfy has a very simply authentication model. Requests that add tokens are
required to carry an *auth_token* in their json data. This *auth_token* is
looked up in the file `auth.db` and the corresponding email address used for
further processing and notification. If gjfy does not find the provided
auth_token, it will reject the request.

Authentication is only for adding new secrets. It does not give access to the
secrets itself.

This authentication model has some downsides and should probably be replaced by
something better. For now just keep in mind that every user in auth.db needs to
have an individual auth_token, because it is used to identify the "user".

To add an account to `auth.db`, simply edit it using your favorite editor and
add a section to the json list that is contained in it, like this:

    {
        "token": "thesecretauthtoken",
        "email": "test@example.org"
    }

Afterwards send gjfy a hangup signal (`killall -HUP gjfy`) to make it reload
the file. In the logfile you will be informed about success or failure.

Usage
-----

Currently the only way to create new secrets is by using the json API. An
example client (gjfy-post) is included. Basically a request looks like this:

    {"auth_token":"g4uhg3iu4h5i3u4","secret":"someSecret"}

By sending this to `/api/v1/new` you create a new URL which is a hash over that
json structure. The reply from the server will tell you this link in both, a
user friendly version and in an API version. Invocation of that link will
immediately lead to deletion of the secret in the server. However, there is an
exception: you can post a `"max_clicks":n` variable along with the json and it
will allow up to `n` clicks.

The authentication token sent with the request will not be stored in the
server. Instead, the associated email address will be stored with the secret,
so it can be used for email notifications (see below).

A timeout can be set by including `"valid_for:n"` in the request. The secret
will become invalid after n days, even if not clicked. The default timeout is 7
days.

Email notifications
-------------------

To get notified if somebody uses the one time link, add the `-notify` flag to
the gjfy command line. gjfy will use the email address associated with the
authentication token that was used when the secret was generated.

This requires that the email sub system of the server where gjfy is running is
configured properly. In principle, if it is possible to send an email using the
`mail` command as the gjfy user, email notifications from gjfy should also
work.

By default, email notification is not enabled.

gjfy-post
---------

gjfy-post is a demonstration client using bash, curl and jq.

    usage: ./gjfy-post <authtoken> <secret> [maxclicks]

Required arguments are authtoken and the secret itself. Please note that
providing the secret this way makes it readable in the system process listing!

The client can be downloaded from the running server by using the URL

    /gjfy-post

Which is also linked from the root page ("/").

You can change the default URL for gjfy-post by setting the environment
variable `GJFY_POSTURL`. If you downloaded gify-post via the URL, it will
have the correct URL already configured in the script.

FAQ
---

Q: How do you pronounce gjfy?

A: It is pronounced like "jiffy".
