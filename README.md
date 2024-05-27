# gdnsd-acme-dns-api

A plugin for gdnsd to automate ACME DNS-01 challenge validation via a lightweight HTTP API.

## Overview

This project provides an HTTP API for automating the DNS-01 challenge validation required for ACME certificate issuance. It integrates with gdnsd, a high-performance, extensible DNS server, allowing for dynamic DNS record updates necessary for the ACME challenge process.

## Features

- **Automated ACME DNS-01 Challenge Handling**: Automatically creates and manages DNS TXT records for ACME challenges.
- **Seamless Integration with gdnsd**: Works as a plugin within the gdnsd environment, ensuring minimal configuration and setup.
- **Support for Multiple ACME Clients**: Compatible with various ACME clients to facilitate SSL certificate automation.
- **High Performance and Reliability**: Leverages gdnsd’s robust architecture for efficient DNS record management.
- **Security and Access Control**: Ensures secure updates to DNS records, with support for access control mechanisms.

## Installation

### Prerequisites

- gdnsd installed and configured on your server.
- An ACME client (e.g., certbot, acme.sh) for certificate requests.
- Docker to build and run the container.

### Building the Docker Image

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/gdnsd-acme-dns-api.git
   cd gdnsd-acme-dns-api
   ```

2. Build the Docker image:

   ```bash
   docker build --tag gdnsd --build-arg GDNS_VER=3.8.2 .
   ```

### Running the Docker Container

```bash
docker-compose up
```

## Usage

1. Configure your ACME client to use the gdnsd-acme-dns-api for DNS-01 challenges. For example, with \`acme.sh\`, you might set the DNS API like this:

   ```bash
   export GDNSD_API="http://yourserver:8080/acme-dns-01"
   export GDNSD_TOKEN="0cb9af673c9af284ba85281053e68820"
   acme.sh --issue --dns dns_gdnsdapi -d yourdomain.com
   ```

2. The API expects a JSON payload with the DNS challenge data. Here's an example of a payload:

   ```json
   {
     "yourdomain.com": "challenge-token"
   }
   ```

3. The API will handle the creation and removal of the necessary DNS records for the challenge.

## Example API Request

You can test the API using \`curl\`:

```bash
curl -X POST http://172.20.0.3:8080/acme-dns-01 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer 0cb9af673c9af284ba85281053e68820" \
     -d '{
              "example.local": "0123456789012345678901234567890123456789012"
         }'

curl -X POST http://172.20.0.3:8080/acme-dns-01-flush \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer 0cb9af673c9af284ba85281053e68820" \
     -d '{
              "submit": "true"
         }'
```

## Development

### Local Development Setup

1. Ensure Go is installed and set up on your machine.
2. Clone the repository and navigate to the project directory.
3. Build and run the Go application:

   ```bash
   go mod download
   go build -o app .
   ./app
   ```

4. The server will start and listen on port \`8080\` by default.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your improvements. Ensure your code adheres to the existing style and includes tests where applicable.

## License

This project is licensed under the MIT License. For more information, please see the [LICENSE](LICENSE) file.

## Built on top of

This project is built on top of [Docker-gdnsd](https://github.com/bedis/Docker-gdnsd).

## Contact

For more information or support, please open an issue on the GitHub repository.