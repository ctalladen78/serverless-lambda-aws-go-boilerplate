
### Environment 

* serverless framework
* npm serverless-offline
* aws dynamodb-local

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