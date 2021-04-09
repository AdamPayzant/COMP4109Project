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
Additionally port 9090 must be free for the server to use.
Once all of this setup is complete, the server is ready.

### Client Host

---

### Client

The client was built using electron. It should run on linux, mac, and windows. The current implimentation has the user 
run the client using electron. Future plans include distributing an client exported as a full electron application.
(Run from an executable and all libraries are included in the file).

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


To build the client, the system must have node v14.0.0 or higher
Navigate the `./src/clientHost/elec-client/` directory and run:

```
npm install
```

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

##### Startup

    Once the node modules are installed, you can start the program run:
    ```
    nmp start
    ```
    The client can be configured using by modifiying the file `userData.json` 


##### Operation

    Once the client has started enter the address for the host in the 
    center field of the host section and press connect. This should connect 
    you to the host server, if not an error will appear.

    Once connected, to start chatting press on the  *One time Connection* button.
    Then fill in the username of who you want to talk to and hit *Connect*


---


## Notes
