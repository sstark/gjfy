FROM debian:bookworm
ARG VERSION="1.2"
RUN adduser --system --home=/home/gjfy gjfy
WORKDIR /tmp
RUN apt-get update
RUN apt-get -y install curl unzip
RUN curl -L -O https://github.com/sstark/gjfy/releases/download/v"$VERSION"/gjfy"$VERSION"-linux-x86_64.zip
RUN unzip gjfy"$VERSION"-linux-x86_64.zip
RUN mkdir /etc/gjfy
RUN mv gjfy/auth.db gjfy/logo.png gjfy/custom.css /etc/gjfy
RUN mv gjfy/gjfy /home/gjfy
USER gjfy:nogroup
ENTRYPOINT ["/home/gjfy/gjfy"]
CMD [""]
