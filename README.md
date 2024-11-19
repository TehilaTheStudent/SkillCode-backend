### note: .vecode/ here couse i develop in codespaces & locally

# SkillCode Platform

### remember Itai tip: 👍 work horizontally, implement the basic, first have a functioning product

##### the purpose of this readme:

the platform im building has to "know" what its capabilities and limitations are, of course, there are things i didnt implement and are not supported yet, i have to specify it, after the development stage ( my stage) its goint next in the pipeline- will the devops know what to do with the product? will the end user understand what i created for him? as a product not as a coding project

Itay: " if there is a feature you didnt implement, tell me, so i wont have to guess"

### Features

- no-user platform
- Manage collection of questions: CRUD methods
- test user submission for questions and give feedback
- ## more specific
  - supported languages: javascript and python
  - utils are provided for user ease (TODO)

---

## Installation -thats DevOps instructions:

### Requirements/ resources needed to run the platform

- enviroment variables??? configuration files?
- docker installed
- ports-
  - MongoDB: `27017`
  - Backend: `8080`
  - Frontend: `3000`
- **mounts**:
  - /var/run/docker.sock:/var/run/docker.sock
  - ~/.kube:/root/.kube
  - Persistent data stored in `mongo-data` volume.
  - mounts the logs

### Quick Start

1. Clone the repo:
   ```bash
   git clone https://github.com/TehilaTheStudent/SkillCode
   cd SkillCode
   ```
2. Start the app:
   ```bash
   docker-compose up - this will pull images...
   ```
3. Access:
   - Frontend: `http://localhost:3000`

---

## Debugging and Testing

- logs from the backend
- running the tests?

## frontend User Instructions

- pages: questions table, code editor for question, form for adding/updating quesion
- explanation about question entity ( rules when creating new question)- what are the validations, what is not checked yet but not allowed?
- instructions for using the code editor (like: you cant import modules and libraries, use only the built in features of the languages)
- help information in the frontend (?) (like in the form, or in the code editor- should i create that?)
-

## developer instructions- adding more languages functionality

- creeate in ./assets/new_language
- implement : generating function sugnature, testing

## backend urls

## Base URL

All endpoints are prefixed with:  
`/skillcode`

---

## Endpoints

### **Questions**

| Method   | URL              | Description              |
| -------- | ---------------- | ------------------------ |
| `POST`   | `/questions`     | Create a new question.   |
| `GET`    | `/questions/:id` | Get question by ID.      |
| `GET`    | `/questions`     | Get all questions.       |
| `PUT`    | `/questions/:id` | Update a question by ID. |
| `DELETE` | `/questions/:id` | Delete a question by ID. |

---

### **Testing**

| Method | URL                   | Description            |
| ------ | --------------------- | ---------------------- |
| `POST` | `/questions/:id/test` | Test a question by ID. |

---

### **Function Signature**

| Method | URL                        | Description                            |
| ------ | -------------------------- | -------------------------------------- |
| `GET`  | `/questions/:id/signature` | Get function signature for a question. |

#### my personal notes

### so what did i do here?

- mongo database
- design of the question entity
- go & gin backend
-
- nuxt frontend

## TODOS:

- google docs about the UI [[UI NOTES](https://docs.google.com/document/d/1ALAKcifoX5DRHbdMJkeR07SC64mj_ZiGxcPbDIpEtEw/edit?usp=sharing)]
