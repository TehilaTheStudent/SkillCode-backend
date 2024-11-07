# SkillCode Backend

## Overview

This repository contains the **BASIC** backend implementation for the SkillCode application. The backend is built using Go and the Gin framework, with MongoDB as the database.

## Highlights

- Implemented basic Go server with CRUD operations on `Question` entity.
- Docker images created: `tehilathestudent/skillcode-backend`.

## API Endpoints

- **Create a Question** (POST): `http://localhost:8080/questions`
- **Get a Question by ID** (GET): `http://localhost:8080/questions/:id`
- **Get All Questions** (GET): `http://localhost:8080/questions`
- **Update a Question** (PUT): `http://localhost:8080/questions/:id`
- **Delete a Question** (DELETE): `http://localhost:8080/questions/:id`
- **Submit and Test Solution** (POST): `http://localhost:8080/questions/:id/test`

### Running the Application

To start the application with Docker Compose, use the following command:

```bash
docker-compose up
```

This command will launch both the MongoDB database and the Go backend services.

