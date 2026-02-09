# MedInfoAssistant Backend

A modern REST API backend for managing medical prescriptions, doctors, and patient information. Built with Go and PostgreSQL, offering a comprehensive solution for medical data management and prescription tracking.

## Features

- **User Management** - Create and manage patient profiles
- **Doctor Management** - Register and manage doctor profiles with accuracy ratings
- **Prescription Management** - Create and track medical prescriptions with symptoms and links
- **Medical Items Tracking** - Track medicines and tests with AI-generated and doctor-provided reasons
- **RESTful API** - Clean, simple REST endpoints for all operations
- **Database Migrations** - Automatic schema creation on startup
- **UUID-based IDs** - Secure unique identifiers for all entities

## Prerequisites

- Go 1.18+ 
- PostgreSQL 12+
- Git

## Installation

### 1. Clone the Repository
```bash
git clone https://github.com/Debojyoti1915001/MedInfoAssistant-Backend.git
cd MedInfoAssistant-Backend
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Setup Environment Variables

Create a `.env` file in the root directory:
```env
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.kmifimmerhntwakrobyg.supabase.co:5432/postgres
PORT=8080
```

Replace `[YOUR-PASSWORD]` with your actual Supabase password.

### 4. Run the Application
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Project Structure

```
MedInfoAssistant-Backend/
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── .env                         # Environment configuration
├── .gitignore                   # Git ignore rules
├── README.md                    # This file
├── config/                      # Configuration files
├── database/
│   └── db.go                   # Database connection & migrations
├── models/
│   ├── user.go                 # User model
│   ├── doctor.go               # Doctor model
│   ├── prescription.go         # Prescription model
│   └── items.go                # Medical items model
├── services/
│   ├── user_service.go         # User business logic
│   ├── doctor_service.go       # Doctor business logic
│   ├── prescription_service.go # Prescription business logic
│   └── items_service.go        # Items business logic
├── handlers/
│   ├── health.go               # Health check handler
│   ├── user_handler.go         # User HTTP handlers
│   ├── doctor_handler.go       # Doctor HTTP handlers
│   ├── prescription_handler.go # Prescription HTTP handlers
│   └── items_handler.go        # Items HTTP handlers
└── routes/
    └── routes.go               # Route registration
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    phnNumber TEXT NOT NULL
);
```

### Doctors Table
```sql
CREATE TABLE doctors (
    docId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    accuracy FLOAT NOT NULL,
    name TEXT NOT NULL,
    phnNumber TEXT NOT NULL,
    spec TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL
);
```

### Prescriptions Table
```sql
CREATE TABLE prescriptions (
    presId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    docId UUID NOT NULL REFERENCES doctors(docId),
    userId UUID NOT NULL REFERENCES users(id),
    symptoms TEXT NOT NULL,
    createdDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    link TEXT
);
```

### Items Table
```sql
CREATE TABLE items (
    itemId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    presId UUID NOT NULL REFERENCES prescriptions(presId),
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('med', 'test')),
    aiReasons JSONB,
    docReason TEXT
);
```

## API Endpoints

### Health Check
```
GET /health
```
Response: `{ "status": "ok", "message": "Server is running" }`

### Users
```
POST /api/users/create
Content-Type: application/json

{
  "name": "John Doe",
  "phnNumber": "9876543210"
}
```

```
GET /api/users
```
Returns all users in the system.

### Doctors
```
POST /api/doctors/create
Content-Type: application/json

{
  "name": "Dr. Smith",
  "phnNumber": "9876543210",
  "spec": "Cardiology",
  "username": "dr_smith",
  "accuracy": 95.5
}
```

```
GET /api/doctors
```
Returns all doctors sorted by accuracy.

```
GET /api/doctors/get?id={docId}
```
Get a specific doctor by ID.

### Prescriptions
```
POST /api/prescriptions/create
Content-Type: application/json

{
  "docId": "uuid-here",
  "userId": "uuid-here",
  "symptoms": "High fever, cough",
  "link": "https://prescription-link.com"
}
```

```
GET /api/prescriptions?userId={userId}
```
Get all prescriptions for a user.

```
GET /api/prescriptions/get?id={presId}
```
Get a specific prescription.

### Medical Items
```
POST /api/items/create
Content-Type: application/json

{
  "presId": "uuid-here",
  "name": "Paracetamol 500mg",
  "type": "med",
  "aiReasons": {
    "fever": "reduces temperature",
    "pain": "pain reliever"
  },
  "docReason": "For symptomatic relief"
}
```

```
GET /api/items?presId={presId}
```
Get all items for a prescription.

```
GET /api/items/get?id={itemId}
```
Get a specific item.

## Usage Example

### 1. Create a User
```bash
curl -X POST http://localhost:8080/api/users/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "phnNumber": "9876543210"
  }'
```

### 2. Create a Doctor
```bash
curl -X POST http://localhost:8080/api/doctors/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dr. Smith",
    "phnNumber": "9876543210",
    "spec": "Cardiology",
    "username": "dr_smith",
    "accuracy": 95.5
  }'
```

### 3. Create a Prescription
```bash
curl -X POST http://localhost:8080/api/prescriptions/create \
  -H "Content-Type: application/json" \
  -d '{
    "docId": "doctor-uuid",
    "userId": "user-uuid",
    "symptoms": "High fever, cough",
    "link": "https://prescription-link.com"
  }'
```

### 4. Add Items to Prescription
```bash
curl -X POST http://localhost:8080/api/items/create \
  -H "Content-Type: application/json" \
  -d '{
    "presId": "prescription-uuid",
    "name": "Paracetamol 500mg",
    "type": "med",
    "aiReasons": {"fever": "reduces temperature"},
    "docReason": "For symptomatic relief"
  }'
```

## Building for Production

```bash
go build -o medinfo-assistant main.go
./medinfo-assistant
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Future Enhancements

- [ ] Authentication & Authorization (JWT)
- [ ] Input validation middleware
- [ ] Rate limiting
- [ ] Pagination for list endpoints
- [ ] Search and filter functionality
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Unit and integration tests
- [ ] Docker containerization
- [ ] CI/CD pipeline

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or suggestions, please open an issue on GitHub.

## Author

Debojyoti1915001

---

**Last Updated:** February 9, 2026