# Hyperactive README

Hyperactive is an API and directory of services on
[Hyperboria](https://wiki.projectmeshnet.org/Hyperboria).


## Usage

To add your Hyperboria server to the directory, run _something like_
the following command (or the programmatic equivalent)

    curl -6 -X POST http://activity.hype/services/new -d \
    '{"name": "Hyperactive", "url": "http://activity.hype/services", "description": "Directory of Hyperboria services exposed as a RESTful-ish API"}'

If it worked, the response will be just what you POSTed, but with a
few key/value pairs added.  Otherwise, you'll see an error (in plain
text, not JSON).

To see the current list of services, visit
<http://activity.hype/services> in your browser or run

    curl -6 -g http://activity.hype/services

at the command line.


## Why did you create this project?

To make it easy for Hyperboria users to see the services and activity
on the network, and much more... soon.  Hopefully.


## TODO

* WebHooks. This would allow servers to subscribe to be notified when
  this directory is updated.

* Programmatic verification that the supposedly-existent
  server/service is actually at the given URL
