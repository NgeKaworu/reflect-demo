package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Test 测试id
type Test struct {
	ID   interface{} `bson:"_id, omitempty, objectId"`
	Name *string     `bson:"name, omitempty"`
}

func main() {
	// var (
	// 	mongo = flag.String("m", "mongodb://root:root@192.168.101.68:27017,192.168.101.69:27017,192.168.101.70:27017/?authSource=admin&replicaSet=rs1", "mongod addr flag")
	// 	//mongo = flag.String("m", "", "mongod addr flag")
	// 	db = flag.String("db", "solitaire_way", "mongod addr flag")
	// )
	// flag.Parse()
	// ctx := context.Background()

	// eng := NewDbEngine()
	// err := eng.Open(*mongo, *db)
	// defer eng.Close()

	// c := eng.GetColl("t_test")

	s := "23333"
	id := "5e2400e9123bb7a386d2f18b"

	o := Test{
		Name: &s,
		ID:   &id,
	}

	// o.ID, _ = primitive.ObjectIDFromHex(*o.ID.(*string))

	// data, err := bson.Marshal(o)

	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)

	tn := t.NumField()

	for i := 0; i < tn; i++ {
		tag := t.Field(i).Tag.Get("bson")
		fmt.Println(tag)
		// for idx, str := range strings.Split(tag, ",") {
		// 	if idx == 0 && str != "" {
		// 		key = str
		// 	}
		// 	switch str {
		// 	case "omitempty":
		// 		st.OmitEmpty = true
		// 	case "minsize":
		// 		st.MinSize = true
		// 	case "truncate":
		// 		st.Truncate = true
		// 	case "inline":
		// 		st.Inline = true
		// 	}
		// }
	}

	fmt.Println(tn, t, v)
	// r, err := c.InsertOne(ctx, data)

	// println(r, err)

	// var ret Test

	// iid, _ := primitive.ObjectIDFromHex(id)
	// res := c.FindOne(ctx, bson.M{
	// 	"_id": iid,
	// })

	// res.Decode(&ret)

	// ret.ID = ret.ID.(primitive.ObjectID).Hex()

	// fmt.Println(err, ret, res)
	// fmt.Printf("%T %#v", ret, ret)
}

type DbEngine struct {
	MgEngine *mongo.Client //关系型数据库引擎
	Mdb      string
	RoleLock sync.RWMutex
}

func NewDbEngine() *DbEngine {
	return &DbEngine{}
}

func (d *DbEngine) Open(mg, mdb string) error {
	d.Mdb = mdb
	ops := options.Client().ApplyURI(mg)
	p := uint64(39000)
	ops.MaxPoolSize = &p
	ops.WriteConcern = writeconcern.New(writeconcern.J(true), writeconcern.W(1))
	ops.ReadPreference = readpref.PrimaryPreferred()
	db, err := mongo.NewClient(ops)

	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = db.Connect(ctx)
	if err != nil {
		return err
	}
	err = db.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		log.Println("ping err", err)
	}

	d.MgEngine = db

	return nil
}

func (d *DbEngine) GetColl(coll string) *mongo.Collection {
	col, _ := d.MgEngine.Database(d.Mdb).Collection(coll).Clone()
	return col
}

func (d *DbEngine) Close() {
	d.MgEngine.Disconnect(context.Background())
}
