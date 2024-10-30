#!/bin/sh

# Wait for postgres
echo "Waiting for postgres..."
until pg_isready -h db -U postgres; do
    sleep 1
done

# Start the application
./main 