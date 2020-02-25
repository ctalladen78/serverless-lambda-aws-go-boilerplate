
## TEST send this jwt to localhost:3000/todo_list
## using postman with bearer token auth header
## manually verify this token by signing at jwt.io with private key `PRIVATE_KEY`
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RlbWFpbCJ9.aVWLjUhHk_6BjlbUT_E-F34Tt31eMa_54_PJJnnJNwk
```

```
Header:
"Authorization":"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RlbWFpbCJ9.aVWLjUhHk_6BjlbUT_E-F34Tt31eMa_54_PJJnnJNwk"
```

### TESTING REQUESTS

* POST endpoint `signup`
```
  header : "x-username" "x-password" "x-email"
  body
  ```
* POST endpoint `signin`
```
  header : "x-username" "x-password" "x-email" "x-access-token"
  body
  ```
* GET endpoint `todo_list`
uses authorization endpoint
```
  header
  body
  ```
* POST endpoint `todo_new`
uses authorization endpoint
```
  header : "Authorization : Bearer <TOKEN>"
  body : "{"data" : "<DATA>"}"
  ```

### TABLE SCHEMA

* todotable:
hash key: "created_by" 
range key: "objectid" // "PREFIX-objectId"

* usertable:
hash key: "email"
range key: "password"