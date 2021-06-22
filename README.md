# Readme

1. Setup MySQL
```
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=dbpass -e MYSQL_DATABASE=hackernews -d mysql:latest
```
2. Create MySQL database
```
docker exec -it mysql bash
mysql -u root -p
dbpass
CREATE DATABASE hackernews;
```
3. Run
````
go run ./server.go
````
---
4. open gql playground
http://localhost:8080/

5. createUser and get jwt token then add token to http headers
```json=
mutation {
  createUser(input: {username: "user1", password: "123", role: EDITOR})
}
```
6. createTask (Auth)
```json=
mutation {
  createTask(input: { title: "task x", content: "2222222t" }) {
    creater {
      name
    }
  	editor {
      name
    }
  }
}
```
7. modifyTask (Auth)
```json=
mutation {
  modifyTask(input: {id: "1", createrID: "1", title: "new title", content: "new-content",editorID: "1",status: COMPLETE})
}
```
8. deleteTask (Auth)
```json=
mutation {
  deleteTask(input: {id: "1", createrID: "1"})
}
```
9. Query Task
```json=
query {
  Task {
    title
    content
    id
    creater{
      id
    }
    status
  }
}
```
10. Query Staff
```json=
query{
  Staff{
    id
    name
    role
  }
}
```
11. Query User
```json=
query{
  User{
    id
    name
    role
  }
}
```
---
* schema
```json=
type Task {
    id: ID!
    title: String!
    content: String!
    creater: User!
    editor: User!
    status: Status!
}

type User {
    id: ID!
    name: String!
    role: Role! 
}

type Query {
    links: [Link!]!
    Task: [Task!]!
    User: [User!]!
    Staff: [User!]!
}

input NewTask {
    title: String!
    content: String!
}

input RefreshTokenInput{
    token: String!
}

input NewUser {
    username: String!
    password: String!
    role: Role! 
}

input DeleteTask {
    id: ID!
    createrID: String!
}

input ModifyTask {
    id: ID!
    createrID: String!
    title: String!
    content: String!
    editorID: String!
    status: Status!
}

input Login {
    username: String!
    password: String!
}

type Mutation {
    createTask(input: NewTask!): Task!
    createUser(input: NewUser!): String!
    deleteTask(input: DeleteTask!): String!
    modifyTask(input: ModifyTask!): String!
    login(input: Login!): String!
    # we'll talk about this in authentication section
    refreshToken(input: RefreshTokenInput!): String!
}

enum Status {
    PENDING
    PROCESSING
    COMPLETE
}

enum Role {
    ADMIN
    STAFF
    EDITOR
}
```