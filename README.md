# **SkillCode Backend**

This repository contains the backend implementation for the SkillCode application. The backend is built using Go and the Gin framework, with MongoDB as the database.

---

### **Highlights**
- A basic Go server with CRUD operations for the `Question` entity.
- Supports testing solutions submitted for questions.
- Docker images created for the backend: `tehilathestudent/skillcode-backend`.

---

### **API Endpoints**
- **Create a Question**: POST `http://localhost:8080/questions`
- **Get a Question by ID**: GET `http://localhost:8080/questions/:id`
- **Get All Questions**: GET `http://localhost:8080/questions`
- **Update a Question**: PUT `http://localhost:8080/questions/:id`
- **Delete a Question**: DELETE `http://localhost:8080/questions/:id`
- **Submit and Test Solution**: POST `http://localhost:8080/questions/:id/test`

---

### **Running the Application**

To start the application with Docker Compose, run:
```bash
docker-compose up
```
This command launches both the MongoDB database and the Go backend service.

---

### **Testing**

The repository includes three levels of tests to ensure application reliability:

#### **Unit Tests**
  ```bash
  go test ./internal/... -v
  ```

#### **Integration Tests**
  ```bash
  docker-compose up -d mongo
  go test ./tests/integration -v
  ```

#### **End-to-End (E2E) Tests**
```bash
docker-compose -f docker-compose.test.yaml up --abort-on-container-exit
```
---



### TODOS:

- add logging
- add error handling

### Decisions:
- start with simple types- not composite inside composite
- no void
- not saving functions signatues in database, although there are reusable
- config exists twice- in the backend and frontend, they have to sync
