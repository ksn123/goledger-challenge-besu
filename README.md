# GoLedger Besu Challenge

This project is a backend Go application that interacts with a **Hyperledger Besu** private blockchain and a **PostgreSQL** database.

It provides a REST API to:

- `SET` a value in the smart contract  
- `GET` the current value from the smart contract  
- `SYNC` the blockchain value to the SQL database  
- `CHECK` if the blockchain value matches the latest database record  

---

## Project Structure

```
goledger-challenge-besu/
├── app/                   # Go application (REST API)
│   ├── main.go
│   ├── handlers/          # HTTP endpoints
│   ├── services/          # Blockchain & DB logic
│   ├── test/              # Unit tests
│   ├── .env
│   ├── docker-compose.yml # PostgreSQL setup
│   └── scripts/init.sql   # DB table creation
├── besu/                  # Besu network configuration
│   ├── artifacts/         # Smart contract ABI
│   └── startDev.sh        # Starts Besu network
```

---

## Requirements

- Docker & Docker Compose  
- Go 1.20+  
- Git  
- Node.js (with npm and npx)  
- [HardHat](https://hardhat.org)  
- Besu Network  
- Unix-based OS or WSL (for Windows users)

---

## Setup & Run

### 1. Clone the Project

```bash
git clone https://github.com/ksn123/goledger-challenge-besu.git
cd goledger-challenge-besu
```

### 2. Start the Besu Network

```bash
cd besu
npm install
chmod +x startDev.sh
sudo ./startDev.sh
```

> If `./startDev.sh` fails due to environment paths:

```bash
sudo env "PATH=$PATH" ./startDev.sh
```

This script will:

- Start a local Besu private network  
- Deploy the smart contract  
- Generate ABI JSON at: `besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json`  
- Output contract address to: `besu/ignition/deployments/chain-1337/deployed_addresses.json`  

> ** Copy the contract address and use it in the `.env` file.**

---

### 3. Configure Environment Variables

Edit the `.env` file inside `app/`:

```
BESU_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=0xYourDeployedContractAddress
ABI_JSON=/absolute/path/to/besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json
PRIVATE_KEY=YourPrivateKey
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=besu_data
```

---

### 4. Setup PostgreSQL Database

#### If using Docker:

```bash
cd app
docker compose up -d
```

#### If using WSL (e.g. Ubuntu on Windows):

```bash
sudo -i -u postgres
psql
CREATE DATABASE besu_data;
\q
```

Create the table:

```bash
psql -h localhost -p 5433 -U postgres -d besu_data -f scripts/init.sql
```

---

### 5. Run the Application

```bash
cd app
go mod tidy
go run main.go
```

The REST API will be running at: `http://localhost:8080`

---

##  API Endpoints

| Method | Endpoint  | Description                                 |
|--------|-----------|---------------------------------------------|
| GET    | `/get`    | Read the current value from the contract    |
| POST   | `/set`    | Set a new value in the smart contract       |
| POST   | `/sync`   | Sync contract value to PostgreSQL database  |
| GET    | `/check`  | Compare blockchain and database values      |

####  Example Calls

```bash
curl -X POST http://localhost:8080/set
curl -X GET http://localhost:8080/get
curl -X POST http://localhost:8080/sync
curl -X GET http://localhost:8080/check
```

---

##  Running Tests

```bash
cd app
go test ./... -v
```

> Ensure `.env` is properly configured and both Besu and PostgreSQL are running.

---

##  Architecture Overview

**Smart Contract:**
- `SimpleStorage.sol` — exposes `set(uint256)` and `get()` methods.

**Go REST API:**
- `handlers/` — HTTP layer
- `services/` — business logic for blockchain + DB

**Database:**
- `contract_state` table stores the latest synced value with timestamps.

**Integration:**
- Interacts with Ethereum JSON-RPC (via `go-ethereum`)
- Uses `database/sql` and `github.com/lib/pq` for PostgreSQL

---

##  Authors

- Forked from: [GoLedgerDev/goledger-challenge-besu](https://github.com/GoLedgerDev/goledger-challenge-besu)  
- Developed by: [Kleison (ksn123)](https://github.com/ksn123)