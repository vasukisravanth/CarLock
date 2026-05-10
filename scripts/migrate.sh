#!/bin/bash

# This script is used to run database migrations for the car lock system.

set -e

# Define the database connection parameters
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="your_username"
DB_PASSWORD="your_password"
DB_NAME="car_lock_system"

# Run the migration
echo "Running migrations..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f ../backend/migrations/0001_init.sql

echo "Migrations completed successfully."