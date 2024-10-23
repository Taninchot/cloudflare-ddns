# cloudflare-ddns

## Overview

`cloudflare-ddns` is a Go application that updates a Cloudflare DNS record with your current public IP address. It
periodically checks your public IP and updates the DNS record if the IP has changed.

## Features

- Fetches the current public IP address.
- Compares the public IP with the existing Cloudflare DNS record.
- Updates the DNS record if the public IP has changed.

## Prerequisites

- Go 1.21 or later
- Docker (optional, for containerized deployment)
- Cloudflare account with API token

## Environment Variables

The application requires the following environment variables:

- `CLOUDFLARE_API_TOKEN`: Your Cloudflare API token.
- `CLOUDFLARE_ZONE_ID`: The ID of the Cloudflare zone.
- `CLOUDFLARE_RECORD_NAME`: The DNS record name to update.
- `CHECK_PUBLIC_IP_INTERVAL`: Interval in milliseconds to check the public IP.

## Usage

### Running Locally

1. Clone the repository:
    ```sh
    git clone https://github.com/taninchot/cloudflare-ddns.git
    cd cloudflare-ddns
    ```

2. Set the required environment variables:
    ```sh
    export CLOUDFLARE_API_TOKEN=your_api_token
    export CLOUDFLARE_ZONE_ID=your_zone_id
    export CLOUDFLARE_RECORD_NAME=your_record_name
    export CHECK_PUBLIC_IP_INTERVAL=your_interval
    ```

3. Build and run the application:
    ```sh
    go build -o cloudflare-ddns
    ./cloudflare-ddns
    ```

### Running with Docker

1. Build the Docker image:
    ```sh
    docker build -t cloudflare-ddns .
    ```

2. Run the Docker container:
    ```sh
    docker run -e CLOUDFLARE_API_TOKEN=your_api_token \
               -e CLOUDFLARE_ZONE_ID=your_zone_id \
               -e CLOUDFLARE_RECORD_NAME=your_record_name \
               -e CHECK_PUBLIC_IP_INTERVAL=60000 \
               cloudflare-ddns
    ```

## Acknowledgements

- [icanhazip.com](https://icanhazip.com) for providing a simple API to get the public IP address.
- [Cloudflare](https://www.cloudflare.com) for their DNS services.