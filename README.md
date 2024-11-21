# **SkillCode Platform Documentation**

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
- `/var/run/docker.sock -> /var/run/docker.sock`: Enables backend to interact with Docker (e.g., Kind).
- `~/.kube -> /root/.kube`: Shares Kubernetes configuration.
- `./logs -> /app/logs`: Stores backend logs.
- **Named Volume**:
  - `mongo-data -> /data/db`: Persistent MongoDB data storage.

- **Ports**:
    - Frontend: `3000`
    - Backend: `8080`
    - MongoDB: `27017`


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
  docker-compose up --pull always
  ```
  - Description: This command pulls the required images, starts the services, and makes the platform accessible at `http://localhost:3000`. a dataset of questions will be there

- check google docs about the UI [[UI NOTES](https://docs.google.com/document/d/1ALAKcifoX5DRHbdMJkeR07SC64mj_ZiGxcPbDIpEtEw/edit?usp=sharing)]
- check backedn logs at logs/app.log
---

### limitations:
- as server starts it checks for dependencies:
- mongo, my-cluster kind cluster, images
- there might be name collision whith my-cluster 
- if something goes wrong it exists, check logs for details
  



## endpoints

| **HTTP Method** | **Endpoint**                         |
| --------------- | ------------------------------------ |
| `POST`          | `/skillcode/questions`               |
| `GET`           | `/skillcode/questions/:id`           |
| `GET`           | `/skillcode/questions`               |
| `PUT`           | `/skillcode/questions/:id`           |
| `DELETE`        | `/skillcode/questions/:id`           |
| `POST`          | `/skillcode/questions/:id/test`      |
| `GET`           | `/skillcode/questions/:id/signature` |
| `GET`           | `/skillcode/ds_utils`                |
| `GET`           | `/skillcode/configs`                 |
