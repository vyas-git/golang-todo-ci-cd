# go-todo-app-rpc
Golang Todo App with rpc protocol ,k8s, github actions


### using docker-compose locally
```
docker-compose build
```
```
docker-compose up
```
* Open http://localhost:3000 - Frontend
* RPC server runs on :9090 


### CI/CD Deploy Production GKE
* Git push to main branch
![alt text](https://raw.githubusercontent.com/saivyas/golang-todo-ci-cd/main/assets/screenshots/cicd_workflow.png)
* Open http://34.79.44.24:3000/ - Frontend
* RPC server runs on :9090 


