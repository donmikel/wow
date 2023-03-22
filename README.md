# World of Wisdom

World of Wisdom is a TCP-server with protection from DDOS based on [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work).

## Description

Design and implement “Word of Wisdom” tcp server. TCP server should be protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
The choice of the PoW algorithm should be explained.
After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
Docker file should be provided both for the server and for the client that solves the PoW challenge.

## Components

World of Wisdom application consists of the following components:

- `server` - provides TCP API endpoints used by clients.
- `guide` - is a special node called by the client when it has to solve a puzzle.
- `client` - calls server to get quote.

### `server` application
The application provides two types of api one to retrieve quotes from the World of Wisdom book, the second is necessary to pass a challenge based on the Guided tour puzzle protocol.

### `guide` application
The application provides api for puzzle solving.

### `client` application
The application calls server to get quote and completes guided tour challenge.

## Protocol

This application uses TCP-based protocol with the following message format.

 `dataSize(4)|id(4)|data(n)`

| segment    | type   | size    | remark                  |
| ---------- | ------ | ------- | ----------------------- |
| `dataSize` | uint32 | 4       | the size of `data` only |
| `id`       | uint32 | 4       |                         |
| `data`     | []byte | dynamic |                         |

All using messages can be found in the `protocol` package [pkg/protocol/protocol.go](pkg/protocol/protocol.go)

## Proof of Work

Idea of Proof of Work for DDOS protection is that client, which wants to get some resource from server, should firstly solve some challenge from server. This challenge should require more computational work on client side and verification of challenge's solution - much less on the server side.

### Algorithm selection

PoW functions can be:

- `CPU-bound` where the computation runs at the speed of the processor, which greatly varies in time, as well as from high-end server to low-end portable devices.
- `Memory-bound` where the computation speed is bound by main memory accesses (either latency or bandwidth), the performance of which is expected to be less sensitive to hardware evolution.
- `Network-bound` if the client must perform few computations, but must collect some tokens from remote servers before querying the final service provider. In this sense, the work is not actually performed by the requester, but it incurs delays anyway because of the latency to get the required tokens.

CPU-bound computational puzzle protocols, such as [Client Puzzle Protocol](https://en.wikipedia.org/wiki/Client_Puzzle_Protocol), [Merkle Tree](https://en.wikipedia.org/wiki/Merkle_tree), [Hashcash](https://en.wikipedia.org/wiki/Hashcash) etc., can mitigate the effect of denial of service attack, because the more an attacker wants to overwhelm the server, the more puzzles it has to compute, and the more it must use its own computational resources. Clients with strong computational power can solve puzzles at much higher rate than destitute clients, and can undesirably take up most of the server resources.

Another crucial shortcoming of computational puzzle protocols is that all clients, including all legitimate clients, are required to perform such CPU-intensive computations that do not contribute to any meaningful service or application.

For implementation i choose [Guided tour puzzle protocol](doc/guidedtour.pdf).

Guided tour puzzle protocol enforces delay on the clients through round trip delays, so that clients' requests arrive at a rate that is sustainable by the server. The advantage of using round-trip delays, as opposed to hard computational problems, is that the round trip delay of a small packet is determined mostly by the processing delays, queuing delays, and propagation delays at the intermediate routers, therefore is beyond the control of end hosts (clients). As such, even an attacker with abundant computational resources cannot prioritize themselves over poorly provisioned legitimate clients.

Furthermore, in guided tour puzzle protocol, the computation required for the client is trivial. Since the length of a guided tour is usually a small number in the order of tens or lower, the bandwidth overhead for completing a guided tour is also trivial. As a result, clients are not burdened with heavy computations (that are usually required by CPU-bound or memory-bound puzzle protocols).

## Repository structure

In addition to the other top-level packages, there are a few special directories that contain specific types of packages:

- **applications** contains packages that compile to applications that are long-running processes (such as API servers).
- **pkg** shared codes, or libraries common across the repo.

## Application (package) source layout

All packages and applications follow the principles of the [DDD, hexagonal architecture](https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together) and some simple conventions for organizing declarations within packages that are designed to help you find the code.

### Layout

- **service/** contains a classic implementation of the business logic service.
- **service.go** contains business logic service interfaces
- **adapters/** contains secondary or driven adapters in terms of [Ports and Adapters](https://herbertograca.com/2017/09/14/ports-adapters-architecture/) architecture
- **handlers/** contains primary or driving adapters in terms of [Ports and Adapters](https://herbertograca.com/2017/09/14/ports-adapters-architecture/) architecture
- **interfaces/ or interfaces.go** contains adapters interface aka ports
- **domain/ or domain.go** contains domain models
- **cmd/** place for `main.go` or other entry points files
- **internal/** contains internal application components, which must follow the same structure as the application itself
- **error.go** this file should contain declarations (both types and vars) for errors that are used by the package.

## Build, Test and Run

### Build

To build all Go binaries run:

    make build

Or build a specific application run: 

    make build_server

To build docker image run:

    make build_docker

### Testing

Tests can be run as easy as:

	make test

### Run locally

    make up

It will start the server, client and two guide services.

We can also look at logs:

    docker-compose -f docker/docker-compose.yml logs -f

## Ways to improve

- implement dynamic activation of PoW protection against an attack based on client behavior heuristics.
- implement a dynamic change in the length of the guided tour depending on the client behavior.
- Implement an online mechanism for exchanging secrets between the server and the guide to simplify configuration and increase system scalability.
- improve unit tests coverage and create integration tests to simulate attacks. 
- use any storage for quotes instead of in-memory array. 