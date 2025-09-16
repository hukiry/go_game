package mongodb

import (
	"context"
	"game/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

const DATA_NAME_DB = "game"
const ERROR_Table = "LoadTable 请先加载数据表"

var mongoUtil *MongoUtil
var mu sync.Mutex

type MongoUtil struct {
	mgoCli *mongo.Client
	mgoDb  *mongo.Database
}

func init() {
	mongoUtil = &MongoUtil{}
	mongoUtil.enableMongoDb()
}
func (this *MongoUtil) initEngine() {
	var err error
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := "mongodb://localhost:27017" //os.Getenv("MONGODB_URI") //

	// 连接到MongoDB
	this.mgoCli, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI))
	if err != nil {
		util.LogError(err)
	}

	// 检查连接
	this.checkConnect()
}
func (this *MongoUtil) checkConnect() {
	// 检查连接
	err := this.mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		util.LogError(err)
	}
}
func (this *MongoUtil) loadDatabase(dataName string) *MongoUtil {
	if this.mgoDb == nil {
		this.mgoDb = this.mgoCli.Database(dataName)
	}
	return this
}
func (this *MongoUtil) enableMongoDb() {
	if this.mgoCli == nil {
		this.initEngine()
	}
	this.loadDatabase(DATA_NAME_DB)
}

type DBConfig[T any] struct {
	filter     *DBFilter
	update     *DBUpdate
	collection *mongo.Collection
	results    []T
}

// 清除并过滤
func (this *DBConfig[T]) BuildFilter() *DBFilter {
	this.filter = &DBFilter{Array: []bson.M{}}
	return this.filter
}

// 清除并设置
func (this *DBConfig[T]) BuildUpdate() *DBUpdate {
	this.update = &DBUpdate{Array: bson.M{}}
	return this.update
}

func (this *DBConfig[T]) LoadTable(tableName string) *DBConfig[T] {
	mu.Lock()
	defer mu.Unlock()
	this.collection = mongoUtil.mgoDb.Collection(tableName)
	return this
}

// 是否写入成功
func (this *DBConfig[T]) Insert(doc T) (isUpdate bool) {
	mu.Lock()
	defer mu.Unlock()
	filter := this.filter.toOnlyId()
	if filter == nil {
		util.LogError("缺少过滤器 filter")
		return false
	}

	if this.isExit() {
		return true
	}

	if this.collection == nil {
		util.LogError(ERROR_Table)
		return false
	}
	_, err := this.collection.InsertOne(context.TODO(), doc)
	if err != nil {
		util.LogError(err)
		return false
	}
	return false
}

func (this *DBConfig[T]) isExit() bool {
	cursor, _ := this.collection.Find(context.TODO(), this.filter.toOnlyId())
	if cursor != nil {
		if err := cursor.All(context.TODO(), &this.results); err != nil {
			util.LogError(err.Error())
		}
	}
	if this.results != nil && len(this.results) > 0 {
		return true
	}
	return false
}

func (this *DBConfig[T]) InsertMany(docs []T) {
	for _, doc := range docs {
		this.Insert(doc)
	}
}

func (this *DBConfig[T]) Update() {
	mu.Lock()
	defer mu.Unlock()
	if this.collection == nil {
		util.LogError(ERROR_Table)
		return
	}
	if this.filter == nil {
		this.BuildFilter()
	}

	if this.update == nil {
		this.BuildUpdate()
	}
	_, err := this.collection.UpdateOne(context.TODO(), this.filter.ToBson(), this.update.ToBson())
	if err != nil {
		util.LogError(err)
	}
}

func (this *DBConfig[T]) Find() *DBConfig[T] {
	mu.Lock()
	defer mu.Unlock()
	if this.collection == nil {
		util.LogError(ERROR_Table)
		return this
	}

	if this.filter == nil {
		this.BuildFilter()
	}
	cursor, err := this.collection.Find(context.TODO(), this.filter.ToBson())

	if err != nil {
		util.LogError(err)
	}

	if cursor != nil {
		if err = cursor.All(context.TODO(), &this.results); err != nil {
			util.LogError(err)
		}
	}

	return this
}

func (this *DBConfig[T]) One() *T {
	var t *T
	if len(this.results) > 0 {
		t = &this.results[0]
		return t
	}
	return nil
}
func (this *DBConfig[T]) All() []T { return this.results }

func (this *DBConfig[T]) Delete() {
	mu.Lock()
	defer mu.Unlock()
	if this.collection == nil {
		util.LogError(ERROR_Table)
		return
	}

	if this.filter == nil {
		util.LogError("缺少过滤器 filter")
		return
	}
	_, err := this.collection.DeleteOne(context.TODO(), this.filter.ToBson())
	if err != nil {
		util.LogError(err)
	}
}
func (this *DBConfig[T]) DeleteMany() {
	mu.Lock()
	defer mu.Unlock()
	if this.collection == nil {
		util.LogError(ERROR_Table)
		return
	}
	if this.filter == nil {
		util.LogError("缺少过滤器 filter")
		return
	}
	_, err := this.collection.DeleteMany(context.TODO(), this.filter.ToBson())
	if err != nil {
		util.LogError(err)
	}
}
func Disconnect() {
	defer func() {
		if err := mongoUtil.mgoCli.Disconnect(context.TODO()); err != nil {
			util.LogError(err)
		}
	}()
}

func (this *DBConfig[T]) Test(doc T) {
	////更新数据
	//res := DBConfig[data.User]{}
	//res.BuildFilter().Eq("RoleID", 0)
	//res.BuildUpdate().Set("item", map[byte]uint32{1: 0, 2: 90, 3: 12})
	//res.LoadTable("user").Update()
	//
	////查找数据
	//res.BuildFilter().Eq("RoleID", 0)
	//result := res.LoadTable("user").Find().All()
	//resultOne := res.LoadTable("user").Find().One()
	//util.Log(result, resultOne)
	//
	////添加数据
	//newRestaurant := data.User{RoleID: 1, ItemMoney: map[byte]uint32{1: 8, 2: 13, 3: 222}}
	//res.BuildFilter().Eq("RoleID", 1)
	//if res.LoadTable("user").Insert(newRestaurant) {
	//	build := res.BuildUpdate()
	//	build.Set("item", newRestaurant.ItemMoney)
	//	res.LoadTable("user").Update()
	//}
	//
	////删除数据
	//res.BuildFilter().Eq("RoleID", 0)
	//res.LoadTable("user").Delete()
	//util.Log(result)
}
