# Hexagonal Architecture Example in Go

A minimal example of hexagonal architecture (ports and adapters) in Go, featuring a REST API, external API client, database CRUD operations, and comprehensive error handling.

## Architecture

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point with DI
├── internal/
│   ├── domain/                  # Business logic and entities
│   │   ├── user.go
│   │   └── weather.go
│   ├── dto/                     # Data Transfer Objects (API contracts)
│   │   ├── user.go
│   │   └── weather.go
│   ├── ports/                   # Interfaces (contracts)
│   │   ├── repository.go
│   │   └── weather_service.go
│   ├── application/             # Use cases/services
│   │   ├── user_service.go
│   │   └── weather_service.go
│   ├── adapters/                # External implementations
│   │   ├── http/               # REST API handlers
│   │   │   ├── handler.go
│   │   │   ├── user_handler.go
│   │   │   └── weather_handler.go
│   │   ├── repository/         # Database implementations
│   │   │   └── memory/
│   │   │       └── user_repository.go
│   │   └── api/                # External API clients
│   │       └── weather_client.go
│   └── errors/                  # Error handling utilities
│       └── errors.go
```

## Features

- **Hexagonal Architecture**: Clean separation between business logic and external concerns
- **DTO Layer**: Data Transfer Objects providing stable API contracts independent of domain models
- **REST API**: User CRUD operations and weather service endpoints
- **External API Integration**: Weather API client with proper error handling
- **Repository Pattern**: In-memory database with interface-based abstraction
- **Comprehensive Error Handling**: Typed errors with HTTP status mapping
- **Dependency Injection**: Clean wiring in main.go
- **Graceful Shutdown**: Proper server lifecycle management

## API Endpoints

### Users
- `POST /api/users` - Create user
- `GET /api/users` - List users (supports limit/offset)
- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

### Weather
- `GET /api/weather?city={city}` - Get weather for city
- `GET /api/users/{id}/weather?city={city}` - Get weather for authenticated user

### Health
- `GET /health` - Health check

## Data Transfer Objects (DTOs)

The application uses DTOs to define stable API contracts that are independent of domain models:

### User DTOs
- `CreateUserDTO` - Request contract for user creation with validation
- `UpdateUserDTO` - Request contract for user updates
- `UserResponseDTO` - Response contract for user data

### Weather DTOs
- `WeatherResponseDTO` - Response contract for weather data

### Benefits
- **API Stability**: Changes to domain models don't break API contracts
- **Validation**: DTOs include validation rules for input data
- **Separation of Concerns**: API layer is decoupled from domain layer
- **Type Safety**: Clear contracts for request/response data structures

### Mapping Functions
- `ToUserResponseDTO(user)` - Converts domain User to response DTO
- `ToUserResponseDTOs(users)` - Converts slice of Users to DTOs
- `ToWeatherResponseDTO(weather)` - Converts domain Weather to response DTO

## Running the Application

```bash
# Set environment variables (optional)
export PORT=8080
export WEATHER_API_KEY=your-openweather-api-key

# Run the server
go run cmd/server/main.go
```

## Example Usage

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","name":"John Doe"}'

# Get user
curl http://localhost:8080/api/users/user_1

# List users
curl http://localhost:8080/api/users?limit=10&offset=0

# Get weather
curl http://localhost:8080/api/weather?city=London

# Update user
curl -X PUT http://localhost:8080/api/users/user_1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe"}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/user_1
```

## Key Design Principles

1. **Dependency Inversion**: High-level modules don't depend on low-level modules
2. **Interface Segregation**: Small, focused interfaces
3. **Single Responsibility**: Each component has one reason to change
4. **API Contract Stability**: DTOs provide stable external contracts independent of domain changes
5. **Separation of Concerns**: Clear boundaries between domain, application, and presentation layers
6. **Testability**: Easy to mock dependencies and test in isolation
7. **Error Handling**: Consistent error types with proper HTTP status mapping