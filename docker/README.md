Example Dockerfile for gjfy
===========================

Use this Dockerfile as a basis for your own setups.

Build container for version 1.2:

    docker build --build-arg version=1.2 . -t gjfy

Run container:

    docker run -d -p 9154:9154 gjfy server

Add parameters:

    docker run -d -p 9154:9154 gjfy --urlbase https://example.org/
