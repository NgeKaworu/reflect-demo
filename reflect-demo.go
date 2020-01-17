package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test 测试id
type Test struct {
	ID primitive.ObjectID `bson:"_id"`
}

func main() {
	fmt.Println("hello world")
}
