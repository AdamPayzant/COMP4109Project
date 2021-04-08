# Secure Messaging and Video System (SMVS)

SMVS is a complete system for secure, peer-to-peer messaging and video.

The system is divided into 3 distinct components: server, client host, and client.
The server stores the usernames, IPs and public keys.
The client host manages all communications going to and from a user.
The client is how individual interfaces with the client.

A more complete description of this system can be found in `Proposal/proposal.pdf`. 

## Installation

### Server

---

The server has only been validated in Linux, though it is likely compatible with other OS's.s
The server can be either built from source or run from a precompiled binary.
If using the default settings, a Mariadb server must be running locally on port 3306. 
This DB must have a user named smvs with the password "password" and have a database named "smvsServer".
Additionally port 50051 must be free for the server to use.
Once all of this setup is complete, the server is ready.

### Client Host

---

### Client

---

## Building

### Server

---

To build the server, the system must have golang v1.16.
Navigate the `./src/protos/` directory and run:

```
make install
```

Once the protocols are install, navigate to the `./src/server/` and run:

```
go build
```

Any changes you want to do to the system (like changing the password for the DB) currently must be by modifying the src files. 
All variables can be found at the top of the `server.go` and `dbhandler.go` files under `const`.

### Client Host

---

### Client

---

## Running

### Server

---
Once the setup has been complete, the server is started by the following single command:
```
./server
```

### Client Host

---

### Client

---


## Notes