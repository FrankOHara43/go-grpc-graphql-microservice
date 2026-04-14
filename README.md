# go-grpc-graphql-micro

A microservice example using:

- gRPC services: `account`, `catalog`, `order`
- GraphQL gateway: `graphql`
- Postgres for `account` and `order`
- Elasticsearch for `catalog`

## Prerequisites

- Docker
- Docker Compose
- Go 1.25+ (for local non-Docker build/run)

## Run With Docker (Recommended)

From project root:

```bash
docker compose up -d --build
```

Check service status:

```bash
docker compose ps
```

Open GraphQL:

- API endpoint: http://localhost:8080/graphql
- Playground: http://localhost:8080/playground

### Exposed Ports

- `graphql`: `8080`
- `account`: `8081`
- `catalog`: `8082`
- `order`: `8083`
- `account_db` (Postgres): `5433`
- `order_db` (Postgres): `5434`
- `catalog_db` (Elasticsearch): `9200`

## Stop Stack

```bash
docker compose down
```

To also remove volumes:

```bash
docker compose down -v
```

## Local Build Check

From project root:

```bash
go build ./...
```

## GraphQL Smoke Test

You can run these requests in order.

1. Create account

```bash
curl -sS -X POST http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  --data '{"query":"mutation { createAccount(account: {name: \"Alice\"}) { id name } }"}'
```

2. Create product

```bash
curl -sS -X POST http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  --data '{"query":"mutation { createProduct(product: {name: \"Book\", description: \"A test book\", price: 10.5}) { id name price } }"}'
```

3. Create order (replace placeholders)

```bash
curl -sS -X POST http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  --data '{"query":"mutation { createOrder(order: {accountId: \"<ACCOUNT_ID>\", products: [{id: \"<PRODUCT_ID>\", quantity: 2}]}) { id totalPrice products { id quantity } } }"}'
```

4. Query accounts with nested orders

```bash
curl -sS -X POST http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  --data '{"query":"query { accounts { id name orders { id totalPrice products { id quantity } } } }"}'
```

## Notes

- `catalog` uses `olivere/elastic.v5`, so Compose is configured with Elasticsearch `5.6.16`.
- The application images use `golang:1.25-alpine` to match `go.mod`.
