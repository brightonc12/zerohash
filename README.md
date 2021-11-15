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


#### cmd
entry point for our server

#### internal
Contains all our application

#### core
container port and services to access our app services
#### services
#### coinbase
#### numbers
#### vwap
#### domain
#### handlers

## Tasks

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



Instructions
-----

1. Clone this repository.
2. Create a new branch called `dev`.
3. Create a pull request from your `dev` branch to the master branch.
