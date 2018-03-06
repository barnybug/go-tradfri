# go-tradfri

A library and command line tool to control the Ikea Tradfri smart home gateway
and bulbs.

## Installation

	$ go get github.com/barnybug/go-tradfri/cmd/tradfri

## Usage

	$ tradfri help

The first time you run you'll need to provide the gateway key (on the
underside of the gateway):

	$ tradfri --gateway 192.168.10.123 --key <KEY> info

After the first run, a shared key is generated, so --key is no longer required.

Search for devices:

	$ tradfri --gateway 192.168.10.123 devices

Switch a bulb on 50% brightness:

	$ tradfri --gateway 192.168.10.123 set --device 65536 --level 50
