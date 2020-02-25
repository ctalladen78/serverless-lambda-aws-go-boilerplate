package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DbController struct {
	conn *dynamodb.DynamoDB
}

// uses localhost only
func InitLocalDbConnection(h string) *DbController {
	return &DbController{
		conn: dynamodb.New(session.New(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String(h),
		})),
	}
}

// get item by key attributes as per table schema
func (ctrl *DbController) GetTodoItem(t *TodoObject, table string) (interface{}, error) {
	// https://github.com/ace-teknologi/memzy
	// https://github.com/nullseed/lesshomeless-backend/blob/master/services/user/dynamodb/dynamodb.go
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go
	// building pkey for search query
	var pkey = map[string]*dynamodb.AttributeValue{
		"objectid": {
			S: aws.String(t.ObjectId),
		},
		// "todo": {
		// 	S: aws.String(t.Todo),
		// },
	}

	// TodoObject and table key attributes do not match because of extra "createdat" field
	// pkey, err := dynamodbattribute.MarshalMap(t)
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       pkey,
	}
	res, err := ctrl.conn.GetItem(input)
	log.Println("GET ITEM output", res)
	if err != nil {
		return nil, err
	}
	var out *TodoObject
	err = dynamodbattribute.UnmarshalMap(res.Item, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// upsert item if exists then replace, otherwise make new item
// ensure item follows attribute value schema
func (ctrl *DbController) PutItem(table string, item interface{}) (interface{}, error) {
	// https://stackoverflow.com/questions/38151687/dynamodb-adding-non-key-attributes/56177142
	newItemAV, err := dynamodbattribute.MarshalMap(item) // conver todo item to av map
	log.Printf("AV Map %v", newItemAV)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      newItemAV,
		TableName: aws.String(table),
	}
	log.Printf("Put Input %v", input)
	o, err := ctrl.conn.PutItem(input)
	if err != nil {
		log.Printf("PUT error", err)
		return nil, err
	}
	var out map[string]interface{}
	log.Printf("Put output %s", o.Attributes)
	dynamodbattribute.UnmarshalMap(o.Attributes, &out)
	return o.Attributes, err
}

// pass in an empty attribute value struct which will be populated as a result
func (ctrl *DbController) ScanUser(table string) (interface{}, error) {
	if ctrl.conn == nil {
		return nil, errors.New("db connection error")
	}
	// get all items in table
	scanOutput, err := ctrl.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	if err != nil {
		fmt.Println("SCAN error", err)
		return nil, err
	}
	var castTo []*UserObject
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &castTo)
	if err != nil {
		return nil, err
	}
	return castTo, nil
}
func (ctrl *DbController) QueryUser(qc QueryCondition, val string) ([]*UserObject, error) {
	condition := ""
	switch qc {
	case EMAIL:
		condition = "email = :val"
	default: // return all items
		condition = ""
	}
	qInput := &dynamodb.QueryInput{
		TableName: aws.String("usertable"),
		// AttributesToGet : , // select only certain attribute values
		KeyConditionExpression: aws.String(condition),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {S: aws.String(val)}, // use this value in an expression
		},
		// KeyConditionExpression: "",
	}
	fmt.Println("QUERY USER BY EMAIL ", qInput)
	res, err := ctrl.conn.Query(qInput)
	castTo := []*UserObject{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &castTo)
	if err != nil {
		return nil, err
	}
	return castTo, nil
}

// using dynamodb.Query as opposed to dynamodb.Scan
// query by enums CREATED_AT | CREATED_BY
// --key-condition-expression 'Artist = :a AND SongTitle BETWEEN :t1 AND :t2' \
func (ctrl *DbController) QueryTodo(table string, qc QueryCondition, user string) ([]*TodoObject, error) {
	condition := ""
	switch qc {
	case CREATED_AT:
		condition = "created_at = :val"
	case CREATED_BY: // return items created by
		condition = "created_by = :user and begins_with(objectid, :type)"
	default: // return all items
		condition = ""
	}
	// https://johnmackenzie.co.uk/post/setting-up-dynamodb-for-local-development/
	// https://www.dynamodbguide.com/expression-basics/
	qInput := &dynamodb.QueryInput{
		TableName: aws.String(table),
		// AttributesToGet : , // select only certain attribute values
		KeyConditionExpression: aws.String(condition),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {S: aws.String(user)},
			":type": {S: aws.String("TODO")},
		},
	}
	fmt.Printf("QUERY TODO BY %s %s ", qc, qInput)
	res, err := ctrl.conn.Query(qInput)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range res.Items {
		fmt.Println("QUERY RESULT ", k, v)
	}
	castTo := []*TodoObject{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &castTo)

	if err != nil {
		return nil, err
	}
	return castTo, nil

}

// update has reference to old data
// returns new updated object
func (ctrl *DbController) Update(table string, request *events.APIGatewayProxyRequest) (interface{}, error) {
	// get post form body
	fmt.Println("REQUEST METHOD", request.HTTPMethod)
	type formInput struct {
		Id string
		Ot string
		Nt string
	}
	form := &formInput{}
	json.Unmarshal([]byte(request.Body), form)
	var err error
	// var keyMapAV2 map[string]*dynamodb.AttributeValue
	// var toUpdate map[string]*dynamodb.AttributeValue
	// http://gist.github.com/doncicuto
	// keyMapAV, err := dynamodbattribute.MarshalMap(oldItem)
	oldItemKeys := map[string]*dynamodb.AttributeValue{
		"id":   {S: aws.String(form.Id)},
		"todo": {S: aws.String(form.Ot)},
	}
	if err != nil {
		return nil, errors.New("itemkey error")
	}
	if err != nil {
		return nil, errors.New("newItem error")
	}
	// https://aws.amazon.com/blogs/developer/introducing-amazon-dynamodb-expression-builder-in-the-aws-sdk-for-go/
	// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.UpdateExpressions.html
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go#L33
	itemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key:       oldItemKeys, // match key attributes per table definition
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {S: aws.String(form.Nt)}, // set new value
		},
		// attribute being updated must not be part of the key
		ExpressionAttributeNames: map[string]*string{
			"#T": aws.String("createdat"),
		},
		// ConditionExpression:	"attribute_exists("#T"),
		// https://gist.github.com/doncicuto/d623ec0e74bf6ea0db7c364d88507393#file-main-go-L63
		ReturnValues:     aws.String("ALL_NEW"),     // enum of ReturnValue class UPDATED_NEW ALL_NEW ALL_OLD
		UpdateExpression: aws.String("set #T = :t"), // SET,REMOVE the attribute to update

	}
	result, err := ctrl.conn.UpdateItem(itemInput)
	if err != nil {
		return nil, err
	}
	out := &TodoObject{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, out)
	if err != nil {
		return nil, err
	}

	log.Printf("UPDATE RESULT %s", result.Attributes)
	return out, nil
}
