package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test 测试id
type Test struct {
	ID   *primitive.ObjectID `bson:"_id, omitempty"`
	Name *string             `bson:"name, omitempty"`
}

func main() {
	s := "233"
	data, err := bson.Marshal(&Test{
		Name: &s,
	})

	var ret Test

	if err := bson.Unmarshal(data, &ret); err != nil {
		fmt.Println(err)
	}

	fmt.Println(data, err, ret)
}
