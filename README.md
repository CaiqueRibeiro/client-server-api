# USD to BRL Quotation Service

A simple and efficient distributed system for fetching, storing, and retrieving USD to BRL exchange rates. Built with Go, this project demonstrates the use of contexts, HTTP clients/servers, timeouts, and SQLite for data persistence.

## 🌟 Overview

This project consists of two main components:

- **Server**: Fetches USD-BRL exchange rates from an external API, stores them in SQLite, and serves them to clients
- **Client**: Requests the current exchange rate from the server and saves it to a text file

The system implements strict timeout controls to ensure responsiveness, and properly handles context cancellation for resource management.

## 🏗️ Architecture

```
├── client/                  # Client application
│   ├── src/
│   │   ├── entities/        # Domain models
│   │   ├── usecases/        # Business logic
│   │   ├── tests/           # Unit and integration tests
│   │   └── main.go          # Entry point
│   └── go.mod               # Dependencies
│
├── server/                  # Server application
│   ├── src/
│   │   ├── gateways/        # External API communication
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── repositories/    # Database operations
│   │   ├── tests/           # Unit and integration tests
│   │   └── main.go          # Entry point
│   ├── go.mod               # Dependencies
│   └── quotations.db        # SQLite database
│
└── Makefile                 # Build and test commands
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later

### Running the Application

1. **Start the server**:
   ```
   make run-server
   ```

2. **Run the client** (in a separate terminal):
   ```
   make run-client
   ```

The client will create a `cotacao.txt` file with the current USD-BRL exchange rate.

### Command-line Options

#### Server
```
go run server/src/main.go -port <port> -db <database_path>
```

#### Client
```
go run client/src/main.go -server <server_url> -output <output_file_path>
```

## ⏱️ Timeout Management

One of the key features of this project is timeout management:

- Server API calls are limited to **200ms**
- Database operations are limited to **10ms**
- Client requests have a timeout of **300ms**

This ensures the system maintains responsiveness even when external services are slow.

## 🧪 Testing

The project includes comprehensive test coverage:

```
make test-all      # Run all tests
make test-unit     # Run unit tests only
make test-integration  # Run integration tests only
```

Tests are structured using the testify package and follow Go best practices:
- Table-driven tests for multiple scenarios
- Mocks for external dependencies
- Integration tests for end-to-end verification
- Context deadline tests

## 🔍 Implementation Details

### Server
The server uses a clean architecture approach:
- **Handlers**: Process HTTP requests and coordinate responses
- **Gateways**: Communicate with external APIs
- **Repositories**: Manage data persistence

### Client
The client follows a similar clean architecture:
- **Entities**: Define the domain models
- **Usecases**: Implement the business logic

Both applications are designed with dependency injection to facilitate testing and maintainability.

## 💭 Design Decisions

1. **SQLite**: Used for simplicity and zero-configuration. For a production system, PostgreSQL or another RDBMS might be more appropriate.

2. **Context Timeouts**: Strict timeouts ensure the system remains responsive even when external services slow down.

3. **Clean Architecture**: Separation of concerns makes the codebase maintainable and testable.

4. **Dependency Injection**: All components accept their dependencies, making them easily testable with mocks.

5. **Unified Error Handling**: Consistent approach to error propagation and logging.