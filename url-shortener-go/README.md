# Shortly: URL Shortener

![Go](https://img.shields.io/badge/Go-1.21+-blue) ![Redis](https://img.shields.io/badge/Redis-7.0+-red) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue)

A high-performance, production-ready URL shortener service built with Go, Redis (for caching), and PostgreSQL (for persistence). Designed for scalability, security, and deployment flexibility.

## Features

- **RESTful API** with JSON responses for easy integration.
- **Base62 encoding** for creating short, memorable, and case-sensitive URLs.
- **Redis Caching** for lightning-fast redirects of frequently accessed URLs.
- **PostgreSQL Persistence** for durable and reliable storage of URL data.
- **Custom URL Expiration** to set a lifetime for your shortened links.
- **Docker-Ready** with a `docker-compose.yml` for easy setup and deployment.
- **12-Factor App Compliant** for modern, scalable, and maintainable web applications.
- **Production-Ready Architecture** with a clear separation of concerns and security best practices.

## Architecture

The project follows a clean and modular architecture, separating concerns into different packages:

- `cmd/`: Main application entry point.
- `internal/config`: Configuration loading from environment variables.
- `internal/handlers`: HTTP request handlers for API endpoints.
- `internal/services`: Business logic and coordination between repositories and caches.
- `internal/repositories`: Data access layer for interacting with the PostgreSQL database.
- `internal/models`: Data structures and models for the application.
- `pkg/utils`: Utility functions for tasks like Base62 encoding and URL normalization.

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 15+
- Redis 7.0+

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/s19835/url-shortener-go.git
    cd url-shortener-go
    ```

2. **Set up environment variables:**

    Create a `.env` file in the root of the project and add the following variables:

    ```bash
    DB_URL=postgres://user:securepassword@localhost:5432/table_name?sslmode=disable
    REDIS_URL=redis://localhost:6379/0
    PORT=8080
    ENVIRONMENT=development
    ```

3. **Run the application:**

    ```bash
    go run main.go
    ```

### Using Docker

For a more streamlined setup, you can use the provided `docker-compose.yml` file:

```bash
docker-compose up -d
```

This will start the Go application, a PostgreSQL database, and a Redis instance.

## API Documentation

### Shorten a URL

- **Endpoint:** `POST /shorten`
- **Description:** Creates a new shortened URL.
- **Request Body:**

  ```json
  {
    "url": "https://example.com/a-very-long-url-to-be-shortened",
    "expires_in": 86400
  }
  ```

- **Response:**

  ```json
  {
    "short_url": "http://localhost:8080/jGfW3b"
  }
  ```

### Redirect to Original URL

- **Endpoint:** `GET /:shortCode`
- **Description:** Redirects to the original URL associated with the short code.
- **Example:**

  ```bash
  GET http://localhost:8080/jGfW3b
  ```

  This will result in a `301 Moved Permanently` redirect to the original URL.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs or feature requests.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add some amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
