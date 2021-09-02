FROM alpine

COPY --from=mattipaksula/command2http /* /usr/bin

ENV COMMAND="ping -c 3 google.com"
ENTRYPOINT [ "command2http" ]