# Car Lock System

This project is a car lock management system that consists of a backend application built with Go and a frontend application built with React. The system allows users to control and manage car locks through a web interface.

## Project Structure

```
car-lock-system
├── backend              # Backend application
│   ├── cmd              # Command line applications
│   ├── internal         # Internal packages
│   ├── pkg              # Public packages
│   ├── migrations       # Database migrations
│   ├── Dockerfile       # Docker configuration
│   ├── go.mod           # Go module definition
│   ├── go.sum           # Go module checksums
│   ├── .env.example     # Example environment variables
│   └── README.md        # Backend documentation
├── frontend             # Frontend application
│   ├── web              # Web application files
│   └── README.md        # Frontend documentation
├── scripts              # Utility scripts
├── .gitignore           # Git ignore file
└── README.md            # Overall project documentation
```

## Backend

The backend is responsible for handling API requests related to car locking features. It includes:

- **cmd/carlock/main.go**: Entry point of the application.
- **internal/api**: Contains HTTP handlers and routes for the API.
- **internal/auth**: Manages user authentication and authorization.
- **internal/lock**: Contains business logic and data access for car locks.
- **internal/db**: Manages the PostgreSQL database connection.
- **internal/config**: Handles application configuration.
- **pkg/models**: Defines data models related to car locks.
- **migrations**: Contains SQL commands for initializing the database schema.

## Frontend

The frontend is a React application that provides a user interface for interacting with the car lock system. It includes:

- **web/src/App.tsx**: Main component of the application.
- **web/src/components/LockControl.tsx**: Component for controlling the car lock feature.
- **web/src/api/client.ts**: API client for making requests to the backend.

## Getting Started

To get started with the project, follow these steps:

1. Clone the repository.
2. Navigate to the `backend` directory and run `go mod tidy` to install dependencies.
3. Set up the database using the SQL commands in the `migrations` folder.
4. Configure environment variables using the `.env.example` file.
5. Start the backend server.
6. Navigate to the `frontend/web` directory and install dependencies using `npm install`.
7. Start the frontend application.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.