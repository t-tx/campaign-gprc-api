
# Campaign Service

This project uses Docker Compose for container management and Protocol Buffers (`protoc`) for generating Go and gRPC code. Below are the steps to build, run, and clean up the project using the Makefile.

## Prerequisites

- Docker and Docker Compose
- Go
- Protocol Buffers (protoc) compiler
- gRPC Gateway
- OpenAPI v2

## Makefile Targets

The Makefile defines the following targets:

### `gen`

Generates Go code and gRPC files from the `.proto` definition.

```bash
make gen
```

This target runs the `protoc` compiler with the necessary options for:
- Generating Go code for gRPC (`--go_out` and `--go-grpc_out`).
- Generating gRPC Gateway code (`--grpc-gateway_out`).
- Generating OpenAPI v2 definitions (`--openapiv2_out`).

### `up`

Builds and starts the Docker containers in detached mode.

```bash
make up
```

This target uses `docker-compose` to:
- Build the images if necessary.
- Start the services in the background (`-d` flag).

### `down`

Stops and removes the Docker containers.

```bash
make down
```

This target stops and removes the containers, networks, and other resources associated with the `docker-compose.yml` configuration.

### `clean`

Stops and removes the Docker containers along with any associated volumes.

```bash
make clean
```

This target performs the same as `down`, but also removes any persistent data stored in volumes (`-v` flag).

## Usage

To get started, follow these steps:

1. **Generate the code**:

   Generate the necessary Go and gRPC code using the `gen` target:

   ```bash
   make gen
   ```

2. **Start the services**:

   Build and start the Docker services:

   ```bash
   make up
   ```

3. **Stop the services**:

   Stop the services:

   ```bash
   make down
   ```

4. **Clean up**:

   Stop the services and remove volumes (use this if you want to clean up data as well):

   ```bash
   make clean
   ```

## Notes

- Ensure that Docker is running and that you have the necessary permissions to build and start containers.
- Modify paths in the Makefile as needed based on your project setup (e.g., if `proto/service.proto` is located elsewhere).
- If you encounter issues with submodules or missing dependencies, make sure to fetch all required files before running these commands.

## License

Include your license information here.
