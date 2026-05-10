# Car Lock System Backend

This is the backend component of the Car Lock System application, built using Go. The backend is responsible for managing the car locking features, handling user authentication, and interacting with the PostgreSQL database.

## Project Structure

- **cmd/carlock/main.go**: Entry point of the application. Initializes the server and listens for incoming requests.
- **internal/api**: Contains the HTTP handlers and routes for the API.
  - **handlers.go**: Defines functions to handle requests related to car locking features.
  - **routes.go**: Sets up the API routes and associates them with handler functions.
- **internal/auth**: Manages user authentication and authorization logic.
  - **auth.go**: Contains functions for user authentication.
- **internal/lock**: Implements the business logic and data access for the car lock feature.
  - **service.go**: Contains business logic for interacting with the lock repository.
  - **repository.go**: Defines data access functions for lock-related data.
- **internal/db**: Manages the PostgreSQL database connection.
  - **postgres.go**: Contains functions to establish and manage the database connection.
- **internal/config**: Handles application configuration settings.
  - **config.go**: Loads and provides access to configuration settings.
- **pkg/models**: Defines data models related to the car lock feature.
  - **lock.go**: Contains structs representing lock data.
- **migrations**: Contains SQL commands for initializing the database schema.
  - **0001_init.sql**: SQL script for setting up the initial database schema.
- **Dockerfile**: Instructions for building a Docker image for the application.
- **go.mod**: Go module definition specifying dependencies.
- **go.sum**: Contains checksums for module dependencies.
- **.env.example**: Example of environment variables needed for the application.

## Getting Started

1. Clone the repository.
2. Navigate to the `backend` directory.
3. Install dependencies using `go mod tidy`.
4. Set up the database using the provided migration scripts.
5. Run the application using `go run cmd/carlock/main.go`.

## Environment Variables

Make sure to create a `.env` file in the backend directory based on the `.env.example` file, and configure the necessary environment variables for your setup.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or features you would like to add.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.