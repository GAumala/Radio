# Distributed Systems Project - Radio

## System dependencies
- [Go Language v1,6+](https://golang.org/dl/)
- [mplayer](https://www.archlinux.org/packages/extra/x86_64/mplayer/)
- [mp3info](https://www.archlinux.org/packages/community/x86_64/mp3info/)
- [mp3splt](https://www.archlinux.org/packages/extra/x86_64/mp3splt/)

## Go dependencies

- [nats](https://github.com/nats-io/nats)
- [gnatsd](https://github.com/nats-io/gnatsd)
- [doublylinkedlist](https://github.com/emirpasic/gods)

## Installation
Use your distro's package manager to install all system dependencies. To install
all Go dependencies you can execute the `install.sh` script.

## Running the Server
First run gnatsd:
```
$GOPATH/bin/gnatsd -DV
```
Then run the radio server (you may want to open a new terminal for this)
```
go run Server/server.go
```

## Running the Client
Once server is running you can start the clients with:

```
go run Client/client.go
```

If the server is running in a different machine, you can pass the server's IP as an argument

```
go run Client/client.go 192.168.11.46
```

### Staff
* Andres Caceres
* Gabriel Aumala
* Jorge Uzca
* Edgar Villaceca
