package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Test 测试id
type Test struct {
	ID   *primitive.ObjectID `bson:"_id, omitempty"`
	Name *string             `bson:"name, omitempty"`
}

func main() {
	const (
		mongo = flag.String("m", "mongodb://root:root@192.168.101.68:27017,192.168.101.69:27017,192.168.101.70:27017/?authSource=admin&replicaSet=rs1", "mongod addr flag")
		//mongo = flag.String("m", "", "mongod addr flag")
		db = flag.String("db", "solitaire_way", "mongod addr flag")
	)
	flag.Parse()

	eng := NewDbEngine()
	err := eng.Open(*mongo, *db)
	defer eng.Close()

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
