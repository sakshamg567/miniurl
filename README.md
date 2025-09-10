# Miniurl: url shortener with microservice architecture

A production-style URL shortener built in Go, designed with a microservices architecture.
It includes caching, database persistence, and is built for scalability.

<hr>

Architecture: 

## Core Services
1. Shortener service - Generates short code for long URL
2. Redirect Service - Takes a long url and redirects to apt longURL
3. Token Service - For generating unique urls everytime based on a base10 counter.

# Project Structure
```
├── cache/                          # cache-ops
├── db/                               # db-ops
├── cmd               
│   ├── gateway/                # main endpoint server
│   └── server/                    # spins up all the services
├── services/
│   ├── redirect/                  # redirection service
│   ├── shortener/               # shortening service
│   └── token/                     # token generation service
├── shared/
│    └── proto/                     #  protobuf definitions
└── docker-compose.yml    # spins db and cache container 
```

<!-- # Services

## Shortener Service
Takes a long url and returns a short code -->

