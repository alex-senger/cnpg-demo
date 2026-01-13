# cnpg-demo

This repository contains example manifests and a CLI tool for demonstrating CloudNativePG, a Kubernetes operator for PostgreSQL databases. These materials were used for a presentation on Kubernetes Operators for PostgreSQL.

## Overview

CloudNativePG is a Kubernetes operator that manages the full lifecycle of PostgreSQL database clusters with high availability, automated failover, and backup capabilities.

## Getting Started

### Prerequisites

- Kubernetes cluster
- CloudNativePG operator installed ([installation guide](https://cloudnative-pg.io/docs/1.28/))
- kubectl configured

### Deployment

1. Create the namespace:

   ```bash
   kubectl apply -f namespace.yaml
   ```

2. Deploy the PostgreSQL cluster:

   ```bash
   kubectl apply -f cluster-example.yaml
   ```

### Using the Inserter Tool

1. Build the inserter:

   ```bash
   cd inserter
   go build -o inserter
   ```

2. Forward the PostgreSQL port (adjust service name as needed):

   ```bash
   kubectl port-forward -n cnpg-demo svc/cluster-example-rw 5432:5432
   ```

3. Run the inserter:

   ```bash
   ./inserter -host localhost -port 5432 -user alex -password admin -database alex
   ```

## Credentials

Default credentials (defined in [cluster-example.yaml](cluster-example.yaml)):

- Username: `alex`
- Password: `admin`

⚠️ **Note:** These are example credentials for demonstration purposes. Never use these in production.
