# gin-go
RBAC scaffolding based on Gin + Gorm/MongoDB + Casbin + Dig
# Project Structure
```bash
.
├── cmd
│   └── server            # Main
├── internal
│   └── app               # Your Application
│       ├── api           # api
│       ├── config        # Config's Structure
│       ├── icontext      # 
│       ├── middleware    # Gin's Middleware
│       ├── models        # Model gorm mongodb
│       ├── repositories  # Repository DB
│       ├── schema        # Schemas
│       ├── services      # Business Logic Layer
│       │   └── impl      # BLL Implement
│       └── test          # Test Cases
└── pkg                   # Common Packages
    ├── app               # Extend Gin
    ├── errors            # Define message, code errors
    ├── jwt               # JWT Auth
    ├── logger            # Logers
    └── utils             # Utilities
```
## Start project
```bash
# start make file
    make start 
# OR run with go command
    go run cmd/server/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml
```


