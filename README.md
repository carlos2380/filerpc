# Filerpc

## Overview

The Filerpc is designed to read files from the disk, calculate their hashes, and store the content in Redis. The service accepts a payload with `type`, `version`, and `hash` as optional parameters, using default values if they are not provided. The service's primary function is to return the file's content and metadata, ensuring data integrity by comparing provided and calculated hashes.

## Service Description

The Filerpc reads files and returns their content along with metadata. If the provided hash differs from the calculated hash, the service returns the correct hash with an empty content field. If the hash is correct and the file exists, the service saves the file's content and its hash in Redis.


## 1- How to run
### Prerequisites

To run the API you need to have Docker and Docker compose installed on the machine.
- Docker (Min version: 20.10.12): https://docs.docker.com/get-docker/
- Docker compose (Min version: 2.27.1): https://docs.docker.com/compose/install/

### Build and run
In the main folder of the project, where the file compose.yml is located. Execute:
```
docker-compose up --build
```
This command builds the Dockerfile and pull a redis and swagger from DockerHub.
Once docker-compose has finished building and running the images and logs. 

#### Swagger
Swagger allows you to visualize and interact with grpc via http using grpc-gateway.
```
http://localhost:8081
```
Initialize the browser on this URL and access the documentation to know the request and response.

![swagger1](https://github.com/carlos2380/webCarlos2380/blob/master/filerpc/Swagger.png)

You can interact with swagger and make requests and see the responses.

![swagger2](https://github.com/carlos2380/webCarlos2380/blob/master/filerpc/Swagger2.png)

- The documentation that uses swagger to work is here:
## 2- Performance
I tested the performance using the client and server on the same host. The results are different than in a real environment where the client and server do not share resources.

The environment to test the performance was an Apple M1 8CPUs and 16GB RAM where Docker resources are 8CPUs and 8GB RAM.

I create a simple json files to make the build and run the code fast. In real development should be good use files similars are the production.

### Own client
I have created a simple client to test server performance.

- https://github.com/carlos2380/filerpc/blob/main/cmd/client/main.go

To run the client, after running the compose. Create the build:
```
docker build -t client --target client .
```
and then, run the client.
```
docker run --network host -it client /client -c 8 -nc 250
```
Where c is the number of threads, nc the number of transactions per thread and url the url to do the get.

To make sure the server has files, I created a basic script createFiles.go which is run in the Dockerfile by default. 
- https://github.com/carlos2380/filerpc/blob/main/createFiles.go

This script creates 8000 files, with 1000 files per thread. To create more files, modify the loop. 

![loop](https://github.com/carlos2380/webCarlos2380/blob/master/filerpc/loop.png)

It is a simple script, as it is assumed that in a real environment, the files would already be present on the server.

#### Results
Executing concurrency 1 and 20000 transactions per thread we have a TPS (Transactions Per Second) of 2319
![TPS 1C](https://github.com/carlos2380/webCarlos2380/blob/master/filerpc/TPS2319.png)

Executing concurrency 4 and 5000 transactions per thread we have a TPS of 7253.
![TPS 4C](https://github.com/carlos2380/webCarlos2380/blob/master/filerpc/TPS7255.png)

### AB Apache
Using AB testing it is easy to check the transactions per second specifying the number of threads and the number of total transactions 

#### Results
Executing concurrency 1 and 15000 transaction we have a TPS of 1482

![AB 1C](https://github.com/carlos2380/webCarlos2380/blob/master/filerpc/AB1.png)
Executing concurrency 4 and 15000 transactions we have a TPS of 3559

![AB 4C](https://github.com/carlos2380/webCarlos2380/blob/master/filerpc/AB4.png)

### Conclusion

There is a significant improvement when using more than one thread, as Go can easily parallelize threads. Testing in a more realistic environment with real files would provide more accurate data regarding TPS (Transactions Per Second) and how calls are parallelized relative to the number of threads.

## 3- Documentation

## gRPC

In this project, I use gRPC to handle communication between services due to its high performance and efficiency. gRPC uses HTTP/2 for transport and Protocol Buffers for serialization, offering features like authentication, load balancing, and health checking. It ensures faster performance and smaller binary sizes compared to traditional RPC. The strongly-typed schema of Protocol Buffers reduces data inconsistencies, making gRPC an ideal choice for robust and scalable service communication.

### Setup

1. Install necessary tools:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

2. Create the file_service.proto file:

    - https://github.com/carlos2380/filerpc/blob/main/internal/proto/file_service.proto

This file defines the gRPC service and messages.

3. Generate Go files from the .proto file:

```
protoc -I . --go_out=. --go-grpc_out=. internal/proto/file_service.proto
```

This setup ensures that your gRPC services are defined and the necessary Go code is generated automatically.
## Database (Redis)

I use Redis to store file content because it is fast and efficient for accessing data in memory, improving the application's speed. Storing files directly on the server is more complex and can lead to data loss, so Redis provides a safer and more scalable solution.

Redis is used to store the file content and it hash. Data is stored using the `HSET` command, where the key is the file path and the fields are the `content` and `hash`.

#### Storing Data in Redis

When a file is read successfully, the file content and its hash are stored in Redis using the following structure:


| Key          | Field     | Description                        |
|--------------|-----------|------------------------------------|
|`<type>/<version>.json`    | `content` | The content of the file            |
|              | `hash`    | The hash of the file content       |

### Accessing Redis via CLI
To interact with the Redis database via the command-line interface (CLI), use the following commands:

1. **Connect to Redis CLI:**
   ```sh
   redis-cli -h <redis-host> -p <redis-port>
   ```

2. **Retrieve data:**
   ```
   HGETALL <path>
   ```
   Example:
   ```
   HGETALL core/1.0.0.json
   ```

## Flags

Flags are used to configure the application at startup.

```go
network := flag.String("network", "tcp", "Network type to use (e.g., tcp, tcp4, tcp6, unix)")
grpcPort := flag.String("grpc-port", "50051", "Port or address to listen on for gRPC")
dbAddr := flag.String("redis-addr", "redis:6379", "Address of the Redis server")
host := flag.String("host", "0.0.0.0", "Host address for the server")
gatewayPort := flag.String("gateway-port", "8080", "Port to run the gRPC-Gateway on")
flag.Parse()
```

## TESTS


### Table Tests

I use table-driven tests to validate our code with various inputs. This method keeps tests organized and maintainable.

- https://github.com/carlos2380/filerpc/blob/main/internal/handler/rpc_test.go

### Integration Tests

Integration tests interact directly with the database and file system, ensuring that all components work together correctly in a real-world environment.

### Mock Tests

Mocks are used to simulate dependencies, allowing isolated testing of individual components. This results in more focused and faster tests.

 - https://github.com/carlos2380/filerpc/blob/main/internal/handler/rpc_mock_test.go

#### GoMock

To install and use mocks with GoMock:

```
go install github.com/golang/mock/mockgen@v1.6.0
````

Generate Mocks:

```
mockgen -source=path/to/your/interface.go -destination=path/to/your/mock/file.go -package=yourpackage
```

## Swagger Documentation with gRPC-Gateway

To create Swagger documentation for our gRPC services and enable testing of the files endpoint directly from Swagger, we used gRPC-Gateway to translate the gRPC API into a RESTful HTTP API.

#### Steps to Implement

##### 1 Install gRPC-Gateway and OpenAPI Plugin:
```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```
##### 2 Clone Google APIs:
```
git clone https://github.com/googleapis/googleapis.git
```
##### 3 Define HTTP Rules in Protobuf:
Annotate service definitions with HTTP rules.
```
import "google/api/annotations.proto";

service FileService {
    rpc ReadFile (FileRequest) returns (FileResponse) {
        option (google.api.http) = {
            get: "/v1/file"
        };
    }
}
```
##### 4 Generate Code:
```
protoc -I . -I ./googleapis --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out . internal/proto/file_service.proto
```
##### 5 Set Up and Serve Swagger:
Implement the gRPC-Gateway server to handle HTTP requests and serve the generated Swagger JSON file.
- https://github.com/carlos2380/filerpc/blob/main/internal/gateway/gateway.go

#### Cors (Cross-Origin Resource Sharing)

CORS is enabled to allow connections from different origins, especially useful for connecting Swagger UI with the API.

```go
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

## Linter

I have integrated a linter into our project to ensure code quality and consistency. The linter checks for potential issues and enforces coding standards.
- https://github.com/carlos2380/filerpc/blob/main/.golangci.yml
##### Intall
````
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
`````
##### Run
```
golangci-lint run
```
## Continuous Integration (CI)

I set up a Continuous Integration (CI) pipeline to automatically test and validate the codebase with each push to the repository. This ensures that our code remains reliable and that new changes do not introduce any issues.

- https://github.com/carlos2380/filerpc/blob/main/.github/workflows/ci.yml

## Error Handling

I implemented a  centralized error handling mechanism with a dedicated errors.go file. 
- https://github.com/carlos2380/filerpc/blob/main/internal/errors/errors.go

This approach ensures consistent and clear error messages, making the application easier to maintain and debug, while adhering to best programming practices.

## Decoupling with Interfaces

In the project, I have decoupled the code by using interfaces for file storage and Redis and handlers. This design allows modularity, making the codebase more flexible, maintainable, and easier to test and extend. For example, by defining interfaces for file storage and Redis interactions, we can easily mock these components during testing, ensuring that each part of our application can be developed and tested in isolation.

## Server Shutdown

To prevent the server from shutting down while there are pending tasks, a graceful shutdown mechanism is implemented. The server captures shutdown signals and waits for pending tasks to complete before shutting down.
```
    sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
	log.Logger.Info("Shutting down server...")
	cancel()

	time.Sleep(2 * time.Second)
	log.Logger.Info("Server stopped")
```

## Areas for Improvement

### Enhanced Test Coverage

More Test Cases: Expand the current test suite to cover a wider range of scenarios, including edge cases and error conditions. This ensures that the application handles all possible inputs and states correctly.

### Automatic Documentation Generation

Integrate automatic generation of documentation into the CI pipeline. This ensures that the documentation is always up-to-date with the latest code changes.

### Monitoring and Metrics

Prometheus Integration: Integrate Prometheus for monitoring application performance and gathering metrics. This helps in identifying performance bottlenecks and monitoring the health of the application.

Grafana Dashboards: Set up Grafana dashboards to visualize the metrics collected by Prometheus, providing a clear and accessible overview of application performance and health.

Health Check: Implement an endpoint dedicated to verifying the status of key system components, such as database connection, service availability, etc.

### Security Enhancements

Implement a robust authentication and authorization mechanisms to secure the application.
