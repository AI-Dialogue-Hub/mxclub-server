# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

mxclub-server is a delivery/errand service platform (similar to Instacart) built with Go 1.23. It uses the jet-web-fasthttp framework (FastHTTP-based) and follows Domain-Driven Design (DDD) principles.

The system consists of two applications:
- **mxclub-admin**: Management backend for operators
- **mxclub-mini**: Client-facing mini program

## Architecture

### Domain-Driven Design Structure

```
domain/
├── common/      # Shared domain models
├── message/     # Messaging domain
├── order/       # Order management (core business logic)
├── payment/     # Payment processing
├── product/     # Product catalog
└── user/        # User management

apps/
├── mxclub-admin/    # Admin panel entry point
└── mxclub-mini/     # Mini program entry point

pkg/common/          # Infrastructure wrappers and utilities
```

Each domain follows this pattern:
- `biz/`: Business logic and use cases
- `entity/`: Domain entities and DTOs
- `po/`: Persistent Objects (GORM models)
- `repo/`: Repository interfaces and implementations
- `util/`: Domain-specific utilities

### Application Structure

Each app (`mxclub-admin` and `mxclub-mini`) contains:
- `main.go`: Application entry point with middleware setup
- `config/`: Configuration loading and validation (YAML-based)
- `controller/`: HTTP request handlers
- `middleware/`: Custom middleware (auth, CORS, etc.)
- `service/`: Application-specific services
- `entity/`: API request/response DTOs

## Common Development Commands

From the `script/` directory:

```bash
# Build both applications
make all

# Build individual apps
make build-admin
make build-mini

# Run applications
make run-admin
make run-mini

# Run tests
make test

# Static analysis
make vet

# Clean build artifacts
make clean
```

Direct Go commands (from project root):

```bash
# Build specific app
go build -o ./bin/mxclub-admin ./apps/mxclub-admin/
go build -o ./bin/mxclub-mini ./apps/mxclub-mini/

# Run tests for all packages
go test ./...

# Static analysis
go vet ./...
```

## Configuration

Both applications require YAML configuration files:

- **Admin**: `./configs/dpp_server.yaml` (config path: `-f ./configs/dpp_server.yaml`)
- **Mini**: Config path in code defaults to developer's local path, use `-f` flag to override

Configuration is validated using `go-playground/validator/v10` on startup.

### Key Configuration Sections
- `server`: Port, JWT key, open API endpoints
- `mysql`: Database connection settings
- `redis`: Connection settings with clustering support
- `wx_pay_config`: WeChat Pay APIv3 credentials
- `upload_config`: File upload (OSS or local storage)
- `logger_config`: Structured logging settings

## Dependency Injection

The framework uses constructor injection via `jet.Provide()`. Services are registered in the `init()` functions of config packages and injected into controllers/services.

Example pattern:
```go
// In config init()
jet.Provide(func() *gorm.DB { return db })
jet.Provide(func() *Service { return service })

// In controller/service constructor
func NewController(db *gorm.DB, service *Service) *Controller {
    return &Controller{db: db, service: service}
}
```

## Database Layer

- **ORM**: GORM v1.25
- **Connection**: Configured via `pkg/common/xmysql`
- **Request-scoped logging**: Request ID is propagated to GORM logger via `xmysql.SetLoggerPrefix(ctx.Logger().ReqId)`
- **Migrations**: Should be handled manually or through GORM AutoMigrate

## Caching

Redis is used via `pkg/common/xredis`:
- Distributed caching with pipeline support
- Custom key naming conventions
- Cluster support

## Key Integrations

- **WeChat Pay APIv3**: `pkg/common/wxpay` - Payment processing
- **WeChat Work (企业微信)**: `pkg/common/wxwork` - Internal communications (admin only)
- **Alibaba Cloud OSS**: `pkg/common/xupload` - File storage
- **SMS**: `pkg/common/txsms` (Tencent) and `pkg/common/sms` (Alibaba)
- **Excel Export**: `pkg/common/xexcel` - Billing and reporting

## Middleware Stack

**Admin application** (apps/mxclub-admin/main.go:15-21):
1. CORS middleware
2. Operator authentication
3. JWT authentication
4. Request tracing
5. Panic recovery

**Mini application**:
1. CORS middleware
2. JWT authentication
3. Request tracing
4. Panic recovery
5. Cron jobs for scheduled tasks

## Business Domain Notes

### Order System
- Supports delayed order visibility for executors ("打手延迟看到订单")
- Tiered executor system (金牌/普通打手)
- Penalty calculations using strategy pattern
- Batch billing export functionality

### Key Concepts
- **打手 (Da Shou)**: Order executors/delivery workers
- **金牌打手**: Premium tier executors with priority access
- **云边电竞**: Multi-tenant support

## Testing

Tests are run using standard Go testing. The Makefile runs tests on all packages except config:
```bash
go list ../... | grep -v config | xargs go test -run .
```

## Deployment

- **Docker**: Separate Dockerfiles for admin and mini apps
- **Base image**: Multi-stage builds with Go 1.23 and Alpine Linux
- **CI/CD**: GitHub Actions workflows deploy to Aliyun Container Registry
- **Strategy**: Multi-instance deployment

## Code Patterns

1. **Controllers**: Handle HTTP requests, delegate to services
2. **Services**: Orchestrate business logic across domains
3. **Repositories**: Data access layer (in domain/*/repo/)
4. **DTOs**: Separate entity packages for API contracts vs domain models
5. **Error Handling**: Use structured errors with context

## Important Notes

- The infrastructure layer (`infra/`) is minimal; most logic resides in domain layer
- Configuration file paths are hardcoded in some places - use `-f` flag to override
- Timezone: Asia/Shanghai
- All logging uses structured logging with request correlation IDs
