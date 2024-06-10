# Fiber Admin

Fiber Admin is a rapid development platform for building management systems based on Fiber + MongoDB + Redis, provides
a set of common features such as user management, role management, documentation management, etc.

## Project Layout

```
.
├── main.go             # Entry point
├── docs                # Swagger API documentation
├── cmd                 # Command line interface
├── configs             # Configuration files
├── internal            # Source code
│   ├── app             # Application
│   └── pkg             # Internal packages
│       ├── api         # API Layer
│       ├── config      # Configuration
│       ├── dao         # Data access object
│       ├── domain      # Domain Layer
│       │   ├── entity  # Entity Struct
│       │   └── vo      # Value Object Struct
│       ├── errors      # Error handling
│       ├── middleware  # Middleware
│       ├── router      # Router Layer
│       ├── service     # Business Logic Layer
│       ├── tasks       # Scheduled tasks
│       ├── validator   # Request validation
│       └── wire        # Dependency Injection
├── pkg                 # Common packages
│   ├── cron            # Cron encapsulation
│   ├── errors          # Custom error
│   ├── jwt             # JWT encapsulation
│   ├── mongo           # MongoDB encapsulation
│   ├── prometheus      # Prometheus encapsulation
│   ├── redis           # Redis encapsulation
│   ├── utils           # Common utils
│   └── zap             # Zap encapsulation
└── test                # Test files

```

## Features

* **Fiber:** Provides a high-performance, minimalist web framework for building RESTful APIs.
* **MongoDB:** Used as the primary database for storing and retrieving data.
* **Redis:** Provides caching capabilities.
* **Casbin:** Role-Based Access Control (RBAC) for managing user permissions, provides a flexible access control model.
* **Zap:** A fast, structured logging library for detailed and efficient logging.
* **Viper:** Used for configuration management, allowing easy configuration handling.
* **Wire:** Dependency injection framework.
* **Swagger:** Integrated for API documentation, allowing automatic generation of API docs.
* **Cron:** Scheduling library for running scheduled tasks.
* **Cobra:** Framework for creating powerful modern CLI applications.
* **Prometheus:** Monitoring and alerting toolkit to track application performance and health.
* **RESTful API:** Provides a RESTful API for interacting with the system.

## Installation

### Clone the repository

```bash
git clone https://github.com/6ixGODD/fiber-admin.git
cd fiber-admin
```

### Install dependencies

```bash
go mod tidy
```

## Usage

### Generate Swagger API documentation

To generate, run:

```bash
swag init
```

### Generate Wire Dependencies

To generate, run:

```bash
wire gen ./internal/pkg/wire
```

### Run the application

To run the application, run:

```bash
go run main.go
```

Here are the available command line parameters:

* `--config`: Specify the configuration file
* `--port`, `-p`: Specify the port to listen on (default is 8080)
* `--host`, `-H`: Specify the host to listen on (default is localhost)
* `--log-level`, `-l`: Specify the log level (default is info)
* `--tls`: Enable TLS (default is false)
* `--tls-cert-file`: Specify the TLS certificate file path (default is empty)
* `--tls-key-file`: Specify the TLS key file path (default is empty)

e.g.:

```bash
go run main.go --config config.yaml --port 8080 --host 127.0.0.1 --log-level debug --tls --tls-cert-file cert.pem --tls-key-file key.pem
```

## License

[Apache Lincense 2.0](LICENSE).
