# tsutil
Tailscale utility for deploying Nodes, Subnet Routers, and App Connectors with Tailscale OAuth Clients.

## Installation Options:
1) Download the binary that matches your system from the release page.
2) Build from source (requires Go 1.21.6 or greater):
   ```bash
   git clone https://github.com/sureapp/tsutil
   cd tsutil/
   make build
   cd bin/
   ./tsutil -h
   ```