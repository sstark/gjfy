Example Dockerfile for gjfy
===========================

Use this Dockerfile as a basis for your own setups.

Build container:

    docker build . -t gjfy

Run container:

    docker run -d -p 9154:9154 gjfy

Add parameters:

    docker run -d -p 9154:9154 gjfy -urlbase https://example.org/
