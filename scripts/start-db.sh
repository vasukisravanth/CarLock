#!/bin/bash

# Start PostgreSQL database service
docker run --name carlock-db -e POSTGRES_USER=youruser -e POSTGRES_PASSWORD=yourpassword -e POSTGRES_DB=carlockdb -p 5432:5432 -d postgres:latest

# Wait for the database to start
echo "Waiting for PostgreSQL to start..."
sleep 10

# Run database migrations
./migrate.sh

echo "Database service started successfully."