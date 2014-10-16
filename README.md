# Hyperactive README

Hyperactive is an API and directory of services on
[Hyperboria](https://wiki.projectmeshnet.org/Hyperboria).


## Usage

To add your Hyperboria server to the directory, run _something like_
the following command (or the programmatic equivalent)

    curl -6 -g -X POST http://[fcaa:c3ef:7d17:db5a:baca:809b:8376:1e6e]:9999/services/new -d \
    '{"name": "Hyperactive", "url": "http://[fcaa:c3ef:7d17:db5a:baca:809b:8376:1e6e]:9999/services", "description": "Directory of Hyperboria services exposed as a RESTful-ish API"}'

If it worked, the response will be just what you POSTed, but with a
few key/value pairs added.  Otherwise, you'll see an error.

To see the current list of services, visit
<http://[fcaa:c3ef:7d17:db5a:baca:809b:8376:1e6e]:9999/services> in your browser or run

    curl -6 -g http://[fcaa:c3ef:7d17:db5a:baca:809b:8376:1e6e]:9999/services

at the command line.


## Why did you create this project?

To make it easy for Hyperboria users to see the services and activity
on the network, and much more... soon.  Hopefully.


## TODO

* WebHooks

  * This would allow servers to subscribe to be notified when this
    directory is updated.

  * WebHooks to specific sites should also be possible.

* Programmatic verification that the supposedly-existent
  server/service is actually at the given URL
