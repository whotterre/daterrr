# Dating App Backend API  
*Current Status: In Development*  

This is the backend service for a dating app, built with:  
- **Go** (Gin framework)  
- **PostgreSQL** (SQLC for type-safe queries)  
- **PASETO** authentication  
- **golang-migrate** for database schema management  

## What's Working Right Now  

### Database  
✅ **User accounts**  
- Secure password storage (bcrypt hashing)  
- Email uniqueness enforcement  
- Schema migrations via `golang-migrate`  

✅ **Profiles**  
- Name, age, gender, bio  
- Location (PostGIS point)  
- Interests (array of strings)  

✅ **Relationships**  
- Swipes (like/pass)  
- Mutual matches  

### API Endpoints  
*(Tested via Postman)*  

**Auth**  
- `POST /auth/signup` - Creates user + profile  
- `POST /auth/login` - Returns JWT token  

**Users**  
- `GET /users/me` - Gets full profile  
- `DELETE /users/me` - Deletes account (cascades)  

**Matching**  
- `POST /swipes` - Records a swipe  
- `GET /matches` - Lists mutual matches  

## Database Migrations

We manage schema changes using `golang-migrate`:

```bash
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_profiles_table.up.sql
└── 000002_create_profiles_table.down.sql
```

**Key commands:**
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq add_columns_to_profiles

# Apply migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback last migration
migrate -path migrations -database "$DATABASE_URL" down 1
```

## How to Run It  

1. **Database Setup**  
```bash
createdb datingapp
migrate -path migrations -database "postgres://localhost:5432/datingapp?sslmode=disable" up
```

2. **Start Server**  
```bash
go run cmd/api/main.go
```

## What's Next  
- [ ] Real-time messaging  
- [ ] Better matching algorithm  
- [ ] Photo uploads  

