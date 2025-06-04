# Finplan-structure

**prod by:**
- Khaldee
- Advisor
- Commander
- Guardian Angel

## Environment Variables

Copy the provided `.env.example` file in each service directory to `.env` and
adjust the values for your environment. The following variables are used:

### API Gateway
- `USER_SERVICE_ADDR` - address of the user service
- `AUTH_SERVICE_ADDR` - address of the auth service
- `TRANSACTION_SERVICE_ADDR` - address of the transaction service
- `JWT_SECRET` - JWT signing secret
- `GATEWAY_PORT` - HTTP port for the gateway (default `8080`)

### Auth Service
- `SERVICE_ADDR` - gRPC address to listen on (default `:50051`)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret

### User Service
- `SERVICE_ADDR` - gRPC address to listen on (default `:50053`)
- `DATABASE_URL` - PostgreSQL connection string

### Transaction Service
- `SERVICE_ADDR` - gRPC address to listen on (default `:50052`)
- `DATABASE_URL` - PostgreSQL connection string
