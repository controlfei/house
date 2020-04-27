# PostHousesImage Service

This is the PostHousesImage service

Generated with

```
micro new house/PostHousesImage --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.PostHousesImage
- Type: srv
- Alias: PostHousesImage

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./PostHousesImage-srv
```

Build a docker image
```
make docker
```