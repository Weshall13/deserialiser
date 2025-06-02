# Protobuf Deserializer Web App

A web application that dynamically deserializes Protocol Buffer messages using a provided schema and generates protobuf messages from JSON data.

## Features

- Dynamic schema parsing
- Hex-encoded message deserialization
- JSON to protobuf message generation
- JSON output with field mappings
- Simple web interface

## Getting Started

### Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)

### Installation

1. Clone the repository
2. Install dependencies:
```sh
go mod tidy
```

### Running the Application

Start the server:
```sh
go run main.go
```

Access the web interface at: http://localhost:8080

## Usage

### Deserializing Messages

1. Enter your Protocol Buffer schema in the schema field
2. Specify the message type to deserialize
3. Paste the hex-encoded protobuf message
4. Click "Deserialize"

Example:
```proto
syntax = "proto3";

message Person {
  string name = 1;
  int32 age = 2;
}
```

Message Type: `Person`

Hex Message:
```
0a044a6f686e201e
```

Expected Output:
```json
{
  "name": "John",
  "age": 30
}
```

### Generating Messages

1. Enter your Protocol Buffer schema in the schema field
2. Specify the message type to generate
3. Enter the JSON data that matches the schema
4. Click "Generate"

Example:
```proto
syntax = "proto3";

message Person {
  string name = 1;
  int32 age = 2;
}
```

Message Type: `Person`

JSON Input:
```json
{
  "name": "John",
  "age": 30
}
```

Expected Output:
```
0a044a6f686e201e
```

## Supported Features

- Nested messages
- Enums
- Repeated fields
- Maps
- All standard protobuf types

## Project Structure

- `main.go` - Main application code
- `templates/index.html` - Web interface template
- `go.mod` - Go module file

## Deployment with Docker

To build and run the application using Docker:

1.  Navigate to the project root directory in your terminal.
2.  Build the Docker image:

    ```bash
    docker build -t protobuf-serdeser .
    ```

3.  Run the Docker container, mapping the application port 8080 to port 8080 on your host machine:

    ```bash
    docker run -p 8080:8080 protobuf-serdeser
    ```

The application should now be accessible in your web browser at `http://localhost:8080`.