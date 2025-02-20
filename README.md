# CodePush Server

A Go-based implementation of a CodePush server for managing Over-The-Air (OTA) updates for mobile applications. This server allows mobile app developers to push updates directly to their users' devices without going through the app store review process.

## Features

- 🚀 Over-The-Air (OTA) updates for mobile applications
- 🔐 Secure authentication with JWT
- 📱 App management (create, delete, rename)
- 🔄 Deployment management with multiple environments (staging, production)
- 📦 Package management with versioning and rollback support
- 👥 Collaborator management for team workflows
- 📊 Metrics and status reporting for deployment tracking
- 🔑 Access key management for secure deployments
- 🌐 RESTful API with versioning support
- 💻 Web-based admin dashboard

## Tech Stack

- **Backend**: Go (Gin Framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Frontend**: HTML, CSS, JavaScript
- **Configuration**: Environment variables (.env)

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher

## Project Structure
```
.
├── config/             # Configuration and environment setup
├── controllers/        # HTTP request handlers
├── frontend/          # Frontend assets and code
│   ├── auth.js        # Authentication logic
│   ├── packages.js    # Package management
│   ├── styles.css     # Global styles
│   └── dashboard.html # Admin dashboard
├── middleware/        # Custom middleware (auth, logging)
├── models/           # Database models
├── routes/           # Route definitions
├── services/         # Business logic layer
├── sql/             # Database migrations and seeds
├── utils/           # Helper functions
├── .env             # Environment configuration
└── main.go          # Application entry point
```

## Setup and Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/code-push-server.git
cd code-push-server
```

2. Create a `.env` file in the root directory:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_NAME=codepushv1
DB_PASSWORD=postgres
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-should-be-very-long-and-secure
JWT_TOKEN_DURATION=24h

# Server Configuration
SERVER_PORT=8080
ENV=development
```

3. Initialize the database:
```bash
go run sql/main.go migrate
go run sql/main.go seed  # Optional: Add sample data
```

4. Start the server:
```bash
go run main.go
```

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/register` - User registration
- `POST /auth/logout` - User logout

### Apps
- `POST /apps` - Create new app
- `DELETE /apps/:appName` - Delete app
- `PATCH /apps/:appName` - Rename app
- `GET /apps/:appName/collaborators` - List collaborators

### Deployments
- `POST /apps/:appName/deployments` - Create deployment
- `POST /apps/:appName/deployments/:deploymentName/release` - Release update
- `POST /apps/:appName/deployments/promote` - Promote deployment
- `POST /apps/:appName/deployments/:deploymentName/rollback` - Rollback deployment

### Packages
- `GET /packages` - List all packages
- `PATCH /packages/:packageId` - Update package status

### Client SDK Endpoints
- `GET /v0.1/public/codepush/update_check` - Check for updates
- `POST /v0.1/public/codepush/report_status/download` - Report download status
- `POST /v0.1/public/codepush/report_status/deploy` - Report deployment status

## Authentication

The server uses JWT (JSON Web Tokens) for authentication. All protected routes require a valid JWT token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

## Development

### Running in Development Mode
```bash
ENV=development go run main.go
```

### Database Migrations
Database schema changes are managed through GORM's AutoMigrate feature and the SQL migration tool:
```bash
go run sql/main.go migrate
```

## Security Considerations

- All passwords are hashed using secure algorithms
- JWT tokens expire after the configured duration
- Environment variables for sensitive configuration
- CORS protection for API endpoints
- Database credentials are never exposed to clients

## Error Handling

The server implements standardized error responses:
```json
{
    "error": "Error message",
    "status": "ERROR",
    "code": 400
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [JWT Go](https://github.com/golang-jwt/jwt)
```

