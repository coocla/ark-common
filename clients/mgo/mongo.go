package mgo

import (
	"context"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// col mongodb集合
type col struct {
	collection *mongo.Collection
}

// Client mongo查询客户端
type Client struct {
	db     *mongo.Database
	conn   *mongo.Client
	dbName string
}

// NewMgo 返回一个新的Mgo连接
func NewMgo(urn string) *Client {
	if urn == "" {
		urn = os.Getenv("MGO_URL")
		if urn == "" {
			log.Fatalf("please set ENV: MGO_URL")
		}
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(urn))
	if err != nil {
		log.Fatalf("Error while init mongodb database, Error: %v", err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to mongodb database, Error: %v", err)
		return nil
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("connect mongodb failed may be timeout, err: %v", err)
		return nil
	}
	dbURL, _ := url.Parse(urn)
	dbName := dbURL.Path[1:]
	dbs, err := client.ListDatabaseNames(ctx, bson.M{
		"name": dbName,
	})
	if err != nil || len(dbs) == 0 {
		log.Fatalf("database [%s] not found or connect mongodb failed: %v", dbName, err)
	}
	log.Info("mongodb connect success")
	return &Client{
		db:     client.Database(dbName),
		conn:   client,
		dbName: dbName,
	}
}

// Table 基于集合的抽象操作
func (c *Client) Table(collName string) *Collection {
	col := &Collection{
		collName:   collName,
		collection: c.db.Collection(collName),
		Client:     c,
	}
	return col
}

// WithTransaction 开启一个事务
func (c *Client) WithTransaction(fn func(sc mongo.SessionContext) error) (err error) {
	opt := options.Session()
	ctx := context.Background()
	session, err := c.conn.StartSession(opt)
	if err != nil {
		return
	}
	defer session.EndSession(ctx)
	err = mongo.WithSession(ctx, session, fn)
	return
}

// Collection 集合对象
type Collection struct {
	collName   string // 集合名
	collection *mongo.Collection
	*Client
}

// IsNotFoundError 不存在的错误
func IsNotFoundError(err error) bool {
	return err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument
}

// IsInvalidHex 判断是否ID格式错误
func IsInvalidHex(err error) bool {
	return err == primitive.ErrInvalidHex
}

func (c *Collection) Insert(data interface{}) (*mongo.InsertOneResult, error) {
	r, err := c.collection.InsertOne(context.Background(), data)
	return r, err
}

func (c *Collection) InsertMany(documents []interface{}, opt *options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := c.collection.InsertMany(context.Background(), documents, opt)
	return result, err
}

func (c *Collection) DeleteMany(filter interface{}) (delCount int64, err error) {
	delRet, err := c.collection.DeleteMany(context.Background(), filter)
	return delRet.DeletedCount, err
}

func (c *Collection) DeleteManyWithResult(filter interface{}) (*mongo.DeleteResult, error) {
	return c.collection.DeleteMany(context.Background(), filter)
}

func (c *Collection) QueryAndUpdate(filter interface{}, update interface{}, opt *options.FindOneAndUpdateOptions) *mongo.SingleResult {
	if opt == nil {
		opt = options.FindOneAndUpdate()
	}
	opt.SetReturnDocument(options.ReturnDocument(1))
	return c.collection.FindOneAndUpdate(context.Background(), filter, update, opt)
}

func (c *Collection) Update(filter interface{}, update interface{}, opt *options.UpdateOptions) error {
	_, err := c.collection.UpdateMany(context.Background(), filter, update, opt)
	return err
}

// Query 查询列表
func (c *Collection) Query(filter interface{}, pageSize, currentPage int, opt *options.FindOptions) (cur *mongo.Cursor, err error) {
	if pageSize > 0 || currentPage > 0 {
		if opt == nil {
			opt = options.Find()
		}
		opt = opt.SetLimit(int64(pageSize))
		opt = opt.SetSkip(int64((currentPage - 1) * pageSize))
	}
	ctx := context.Background()
	cur, err = c.collection.Find(ctx, filter, opt)
	return
}

// Count 计算总数
func (c *Collection) Count(filter interface{}, opt *options.CountOptions) (count int64, err error) {
	ctx := context.Background()
	if opt == nil {
		opt = options.Count()
	}
	count, err = c.collection.CountDocuments(ctx, filter, opt)
	return
}

// QueryOne 查询一个
func (c *Collection) QueryOne(filter interface{}, target interface{}, opt *options.FindOneOptions) error {
	if opt == nil {
		opt = options.FindOne()
	}
	return c.collection.FindOne(context.Background(), filter, opt).Decode(target)
}

// Replace 插入或更新
func (c *Collection) Replace(filter interface{}, target interface{}, opt *options.FindOneAndReplaceOptions) error {
	if opt == nil {
		opt = options.FindOneAndReplace()
		opt.SetUpsert(true)
	}
	return c.collection.FindOneAndReplace(context.Background(), filter, target, opt).Err()
}
