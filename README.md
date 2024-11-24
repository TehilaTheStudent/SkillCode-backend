**Skill Code Platform Documentation**

## Backend: [https://github.com/TehilaTheStudent/SkillCode-backend](https://github.com/TehilaTheStudent/SkillCode-backend)
## Frontend: [https://github.com/TehilaTheStudent/SkillCode-frontend](https://github.com/TehilaTheStudent/SkillCode-frontend)

- user-less platform
- Manage collection of questions: CRUD methods
- test user submission in javascript & python and give feedback
- supports serving request concurrently
- Question are managed regardless of the languages, with general data types, and function signature generated based on that, making it easy to add new language support, and safer Questions Add/Edit actions
---

## **Installation**
- **System Requirements**:
  - Docker compose,
  - kind

 **Mounts**
- `./seed:/seed` for populating the backend
- `/var/run/docker.sock -> /var/run/docker.sock`: Enables backend to interact with Docker (e.g., Kind).
- `~/.kube -> /root/.kube`: Shares Kubernetes configuration.
- `./logs -> /app/logs`: Stores backend logs.
- **Named Volume**:
  - `mongo-data -> /data/db`: Persistent MongoDB data storage.
  

- **Ports**:
    - Frontend: `3000`
    - Backend: `8080`
    - MongoDB: `27017`
    - kind cluster `37000`


---
## **Getting Started**
- **Clone the Repository**:
  - Command to clone the repo:
   ```bash
   git clone https://github.com/TehilaTheStudent/SkillCode
   cd SkillCode
   ```

- **Run with Docker Compose**:
  ```bash
  docker-compose up 
  ```
  - Description: This command pulls the required images, starts the services, and makes the platform accessible at `http://localhost:3000`. a dataset of questions will be there, 
- ```bash
  docker-compose up --scale seed-db=0
  ```
  - to start without the seed db

- check google docs about the UI [[UI NOTES](https://docs.google.com/document/d/1ALAKcifoX5DRHbdMJkeR07SC64mj_ZiGxcPbDIpEtEw/edit?usp=sharing)]
- check backedn logs at logs dir
---
### issue: 
- skillcode-cluster can be present but container skillcode-cluster-control-plane missing (deleted by  user couse its created by default)- server will exit immediately



### limitations:
- as server starts it checks for dependencies:
- connect to mongo, skillcode-cluster exists, kind installed
- there might be name collision whith skillcode-cluster 
- if something goes wrong it exists, check logs for details
  



## endpoints



| **Method** | **Endpoint**                            | **Description**                                   |
|------------|-----------------------------------------|---------------------------------------------------|
| POST       | `/skillcode/questions`                 | Create a new question.                           |
| GET        | `/skillcode/questions/:id`             | Retrieve a question by its ID.                   |
| GET        | `/skillcode/questions`                 | Retrieve all questions.                          |
| PUT        | `/skillcode/questions/:id`             | Update a specific question by its ID.            |
| DELETE     | `/skillcode/questions/:id`             | Delete a specific question by its ID.            |
| POST       | `/skillcode/questions/:id/test`        | Test a question with provided inputs.            |
| GET        | `/skillcode/questions/:id/signature`   | Get the function signature of a specific question.|
| GET        | `/skillcode/ds_utils`                  | Serve utility functions/data structures.          |
| POST       | `/skillcode/ds_utils/examples`         | Generate examples for data structures.            |

This version simplifies the view while retaining the key details about each endpoint.


### running user submissions locally for debugging:
MODE_ENV=development docker-compose up


- dont put nul in data types
- dont use the void
- i wait up to 15 sec to job, im sure i miss something couse if job fails it should get failed feedback immidiately 

- you can run in development mode- no jobs, but run the dev/dev.sh script before and after

### cleanup
-delete the skillcode-cluster 


### testing in parallel
- run k6 run tests/load-test.js

if you dont have kind
---

### **Install Kind**
```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.25.0/kind-linux-amd64

chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
kind --version
```
