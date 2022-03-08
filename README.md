# Owntracks to MQTT gateway

A simple program I'm writing as part of learning Go.
It takes HTTP posts from the Owntracks Android app and submits the data to an MQTT broker.

This is an alternative to using a Lua script with ot-recorder.

I have targeted a Podman/Docker container, though it will work on its own.

## Building

Taskfile is used for building the app and the docker container.

## Configuration

Configuration is supplied through environment variables.