# Chatroom

This is a simple chatroom server and cli client. I am using this project to practice Go!

Chats are sent between client and server using the Gorilla Websocket package.

# Getting Started

## Running

Use something along the lines of
```
cd cmd/server
go run .
```
and
```
cd cmd/cli
go run .
```

to run.

## Client configuration

The environment variables 

* `GO_CHAT_HOST`
* `GO_CHAT_USERNAME`
* `GO_CHAT_ROOM` 

can be used to skip prompts when starting a client connection.

Use `source util/profile1.sh` to load an example profile.