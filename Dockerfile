FROM golang:1.23

# Install Trivy dependencies and Trivy itself
RUN apt-get update && apt-get install -y wget apt-transport-https gnupg lsb-release && \
    wget https://github.com/aquasecurity/trivy/releases/download/v0.51.1/trivy_0.51.1_Linux-64bit.deb && \
    dpkg -i trivy_0.51.1_Linux-64bit.deb && \
    rm trivy_0.51.1_Linux-64bit.deb

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the app code
COPY . .

# Expose app port
EXPOSE 8080

# Run the app using go run
CMD ["go", "run", "."]
