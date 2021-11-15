Documentation
===

This Documentation serves as a way to illustrate how to go through
this worked out solution

The Application can be run as a CLI application by using:-
- Maker file command
- bash script file
- docker container


## Requirements

1. **Docker**
2. **Golang Installed**

## Execution

###  Docker
The application needs to be build to a docker image for it to run

    ~ docker build -t zerohash/vwap:latest .

Then you can use the following script to run a new container

    ~ docker run --rm zerohash/vwap:latest

### Golang
To run the application with go, you need to make sure you have golang installed and running on your machine, and then run the following command:

    ~ go run . \


## Using Makefile
The Makerfile has the script to build image as well as running the application
- `make build-image` - is used to package the application with the help into executable image. The image can be run on any docker
  environment machine. basically `build-image` script consist of script to build an image
- `make go-run` - executes the as a golang application


## Structure
These are modules in this application
* [cmd](./cmd)
    * [server.go](./cmd/server.go)

* [internal](./internal)
    * [core](./internal/core)
        * [services](./internal/core/services)
            * [coinbase](./internal/core/services/coinbase)
            * [numbers](./internal/core/services/numbers)
            * [vwap](./internal/core/services/vwap)
    * [domain](./internal/domain)
    * [handlers](./internal/handlers)


[comment]: <> (#### cmd)

[comment]: <> (entry point for our server)

[comment]: <> (#### internal)

[comment]: <> (Contains all our application)

[comment]: <> (#### core)

[comment]: <> (container port and services to access our app services)

[comment]: <> (#### services)

[comment]: <> (#### coinbase)

[comment]: <> (#### numbers)

[comment]: <> (#### vwap)

[comment]: <> (#### domain)

[comment]: <> (#### handlers)

[comment]: <> (## Tasks)

### Functional Requirements
- [x] read trade from coinbase websocket on a 'matches' channel
- [x] Pull data by currencies
     - BTC-USD 
     - ETH-USD
     - ETH-BTC
- [x] Calculate VWap with a sliding window of 200 data point
- [x] stream out the trade update through a go channel
- [x] print out the result on every update


VWap Calculator
====
![alt vwap](https://zamezii-warehouse-test.s3.us-east-2.amazonaws.com/images/temp/site/vwap.png)


Instructions
-----

1. Clone this repository.
2. Create a new branch called `dev`.
3. Create a pull request from your `dev` branch to the master branch.
