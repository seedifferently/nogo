nogo
====

[![Linux Build Status](https://img.shields.io/travis/seedifferently/nogo.svg?style=flat-square&label=linux+build)](https://travis-ci.org/seedifferently/nogo) [![Windows Build Status](https://img.shields.io/appveyor/ci/seedifferently/nogo.svg?style=flat-square&label=windows+build)](https://ci.appveyor.com/project/seedifferently/nogo)


What?
-----

`nogo` blocks access to various sites (ads, tracking, porn, gambling, etc) by
acting as a DNS proxy server with host blacklist support.


Why?
----

I wanted an open source ad blocker solution that was more universal than a
browser plugin, and:

* Was easy to utilize with unrooted mobile devices (so that [battery life could
  be conserved][1]).
* Had a basic web control panel and API for adding, removing, and "pausing"
  hosts.
* Provided straightforward cross-platform support and acceptable performance (so
  that I could run it from my Raspberry Pi).
* Could be used as a master host "blacklist" service for network-wide ad
  blocking (e.g. by configuring the DNS on my router to point to `nogo`).

[1]: https://lifehacker.com/ad-blockers-on-mobile-can-reduce-battery-drain-by-up-to-1764344384


How?
----

You may simply [download a binary release](https://github.com/seedifferently/nogo/releases)
for your platform, or you can follow the steps below to build from source:

1. Install [Go](https://golang.org/doc/install) (v1.8 or later is required).

2. Clone the repo, then `cd` into it.

3. Install the dependencies by running `make deps`. Or if you don't have `make`:
    * `go get github.com/miekg/dns`
    * `go get github.com/boltdb/bolt`
    * `go get github.com/pressly/chi`

4. Build the app by running `make`. Or if you don't have `make`: `go build`

5. Run the app: `sudo ./nogo`

**Note:**

* Since `nogo` binds to port `:53` by default, it must be given access to
  "privileged" ports (e.g. via `setuid` or `sudo`).
* Run with the `-help` switch for information on additional runtime options
  (such as disabling or password protecting the web control panel).


### Important post-install steps:

#### 1. You must add hosts to the blacklist.

`nogo` doesn't ship with a built-in blacklist, so it won't block any hosts until
you add them.

There are currently two methods for adding hosts to the blacklist:

1. Navigate to the web control panel (default: [http://localhost:8080/][1]) and
   add a host using the form.

2. Download a popular hosts list file (e.g. pick one from the list at
   [https://github.com/StevenBlack/hosts][2]), and execute `nogo` with the
   `-import` switch on its first run.


#### 2. You must reconfigure your DNS.

Your computer/mobile device/etc is probably set up by default to utilize a DNS
server which allows connections to any host. Unless you update your DNS
configuration to point to `nogo` (and *only* to `nogo`), nothing will change.

For those of you who may be unfamiliar with how to update your DNS
configuration, check out Google's guide for their DNS service here:
[https://developers.google.com/speed/public-dns/docs/using][3]

You can follow their instructions, but don't forget to substitute their DNS
service IP addresses with the sole IP address of the machine running `nogo`.

[1]: http://localhost:8080/
[2]: https://github.com/StevenBlack/hosts
[3]: https://developers.google.com/speed/public-dns/docs/using


### Known Issues and Limitations

* The DNS proxy server utilizes a fairly basic configuration, so advanced
  features such as EDNS and DNSSEC are not currently supported.
* Due to the fact that the web control panel utilizes a few modern techniques
  (such as [flexbox][1] and the [Fetch API][2]), you may experience some issues
  with its interface on non-current browsers.

[1]: https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Using_CSS_flexible_boxes
[2]: https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API


Who?
----

My name is Seth and I've been talking to machines in various languages since the
early 90s. If you find this useful and want to say thanks, feel free to
[tweet me][1], [buy me a beer][2], [share some Satoshi][3], or pass
[my resume][4] on to someone you know who is tackling interesting problems with
software.

[1]: https://twitter.com/seedifferently
[2]: https://paypal.me/seedifferently
[3]: https://coinbase.com/seedifferently
[4]: https://resume.sethdavis.name


Copyright (c) 2017 Seth Davis