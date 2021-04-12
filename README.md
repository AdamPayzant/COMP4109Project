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
If using the default settings, a `Mariadb` server must be running locally on port 3306. 
This DB must have a user named `smvs` with the password `"password"` and have a database named `"smvsserver"`.
For a local database server, running the following query will initialize the database.
```
CREATE USER 'smvs'@localhost IDENTIFIED BY 'password';
CREATE DATABASE smvsserver;
GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, UPDATE ON smvsserver.* TO 'smvs'@localhost;
FLUSH PRIVILEGES;
```
Additionally port 9090 must be free for the server to use.
Once all of this setup is complete, the server is ready.

### Client Host
The client host has only been validated in windows, though it is likely compatible with other OS's.
The client host can be either built from source or run from a precompiled binary.
The client host needs a `JSON` settings file as input. The fields in the settings file are:
```
{
	"ClientPublicKeyPath": "<client public key path>.pem",
	"ServerCert": "<Path for the clienthost cert>",
	"ServerKey": "<Path for the clienthost server-key>",
	"DB": "<the name of the DB>",
	"ServerIP": "<The IP of the host mechine i.e LocalHost:<ports>>",
	"Username": "<Username>",
	"CentralServerIP": "<The ip of the centrial server>",
	"CentralServerCACert": "<The CA cert for the central server>",
	"token": "<The token to be sigend to verify the client - clenthost interactions>"
}
```

The client host uses a database that needs to be setup.
If using the default settings, a `Mariadb` server must be running locally on port 3306. 
This DB must have the User, Password, and database specified in the `"DB"` string in the settings file.
For instance the current Default settings are:
```
{
	"ClientPublicKeyPath": "./test/keys/client_public.pem",
	"ServerCert": "./certs/server-cert.pem",
	"ServerKey": "./certs/server-key.pem",
	"DB": "smvs:password@tcp(localhost:3306)/smvsclienthost",
	"ServerIP": "localhost:8081",
	"Username": "Default",
	"CentralServerIP": "localHost:9090",
	"CentralServerCACert": "../server/certs/ca-cert.pem",
	"token": "test"
}
```
Then the database must have the user `smvs` with the password `"password"` and there must be a `"smvsclienthost"` database.
For the Dufault `"DB"` setting and a local database server, running the following query will initialize the database.
```
CREATE USER 'smvs'@localhost IDENTIFIED BY 'password';
CREATE DATABASE smvsclienthost;
GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, UPDATE ON smvsclienthost.* TO 'smvs'@localhost;
FLUSH PRIVILEGES;
```
---

### Client

The client was built using electron. It should run on linux, mac, and windows. The current implimentation has the user 
run the client using electron. ~~Future plans include distributing an client exported as a full electron application.~~
(See section *Final State of the client*)

---

## Building

### Server

---

To build the server, the system must have `golang v1.16`.
Navigate the `./src/protos/` directory and run:

```
make install
```

Once the protocols are install, navigate to the `./src/server/certs` and run:

```
./gen.sh
```
navigate to `./scr/server` then run:
```
go build
```

Or instead after the protocols are install, navigate to `./src/server` and run:
```
make server
```

To clean, navigate to `./src/server` and run:
```
make clean
```

Any changes you want to do to the system (like changing the password for the DB) currently must be by modifying the src files. 
All variables can be found at the top of the `server.go` and `dbhandler.go` files under `const`.

### Client Host
To build the client host, the system must have `golang v1.16`.
Navigate the `./src/protos/` directory and run:

```
make install
```

Once the protocols are install, navigate to the `./src/clientHost/certs` and run:

```
./gen.sh
```
navigate to `./scr/clientHost` then run:
```
go build
```

Or instead after the protocols are install, navigate to `./src/clientHost` and run:
```
make clientHost
```

To clean, navigate to `./src/clientHost` and run:
```
make clean
```


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
Once the setup has been complete, the server is started by the following single command:
```
./clientHost <path to the settings JSON file>
```

---

### Client

##### ~~Startup~~

Once the node modules are installed, you can start the program run:
```
nmp start < dir >
```
The client can be configured by modifiying the file `userData.json` found in the directiony specified as < dir > (The directory which it is run in is default)

The ip of the clienthost and username can be changed in the program, but the address for the centralServer needs to be edited in the `userData.json` file.

The fields in the json file being *clientHostIP* *name* and *centralserverIP* respectively. 


##### ~~Operation~~

~~Once the client has started enter the address for the host in the center field of the host section and press connect. This should connect you to the host server, if not an error will appear.~~

~~Once connected, to start chatting press on the  *One time Connection* button.Then fill in the username of who you want to talk to and hit *Connect*~~

#### Final State of the client
Due to outstanding sercumstances the client will be unable to operate. In theory the client would be able call the client host and the central server.
Due to outstanding issues with the implimentation of the grpc calls in the client, it was abandoned in order to hand something in.  Most of the client functions that would have reached out to the other parts have been gutted since they were broken. The one that remains is the Main connection function that creates a connection to both services, and will result in the error that has been plauging me, *the connection has been dropped*.

The fault for this failure should be placed on myself (Anders Sonstenes).


---

## Notes
