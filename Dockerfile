# Use the official Golang image as the base image for the build stage
FROM golang:1.22.0-alpine AS build

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the warnings-discord-bot binary
RUN go build -o warnings-discord-bot

# Use a smaller base image for the final stage
FROM alpine:3.19

# Copy the warnings-discord-bot binary from the build stage
COPY --from=build /app/warnings-discord-bot /usr/local/bin/

# Expose port 9153 for the warnings-discord-bot to publish metrics
EXPOSE 9153

# Run the warnings-discord-bot binary when the container starts
CMD ["warnings-discord-bot"]
