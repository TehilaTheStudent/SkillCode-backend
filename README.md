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
- my-cluster can be present but container my-cluster-control-plane missing (deleted by  user couse its created by default)- server will exit immediately



### limitations:
- as server starts it checks for dependencies:
- connect to mongo, my-cluster exists, kind installed
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


### running user submissions locally for debugging:
MODE_ENV=development docker-compose up


- dont put nul in data types
- dont use the void
- i wait up to 15 sec to job, im sure i miss something couse if job fails it should get failed feedback immidiately 

- you can run in development mode- no jobs, but run the dev/dev.sh script before and after




if you dont have kubectl/kind
---

### **Install Kind**
```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
kind --version
```

---

### **Install kubectl**//i have to check if this is needed
```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version --client
```

---

These commands install the latest versions of **Kind** and **kubectl** and verify their installations.
