
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


Serverless: Stack update finished...
Service Information
service: aws-auth-cognito
stage: dev
region: eu-west-1
stack: aws-auth-cognito-dev
resources: 22
api keys:
  None
endpoints:
  GET - https://6sw3x5b3lf.execute-api.eu-west-1.amazonaws.com/dev/hello-world
  GET - https://6sw3x5b3lf.execute-api.eu-west-1.amazonaws.com/dev/todo
functions:
  authorize: aws-auth-cognito-dev-authorize
  hello-world: aws-auth-cognito-dev-hello-world
  todo: aws-auth-cognito-dev-todo
layers:
  None

Serverless: Stack update finished...
Service Information
service: aws-auth-cognito
stage: dev
region: us-east-1
stack: aws-auth-cognito-dev
resources: 22
api keys:
  None
endpoints:
  GET - https://ns47vf9fu0.execute-api.us-east-1.amazonaws.com/dev/hello-world
  GET - https://ns47vf9fu0.execute-api.us-east-1.amazonaws.com/dev/todo
functions:
  authorize: aws-auth-cognito-dev-authorize
  hello-world: aws-auth-cognito-dev-hello-world
  todo: aws-auth-cognito-dev-todo