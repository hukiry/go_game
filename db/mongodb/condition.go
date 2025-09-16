package mongodb

import (
	"game/util"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type DBFilter struct {
	Array []bson.M
	or    []bson.M
	and   []bson.M
}

// = 等于
func (this *DBFilter) Eq(fieldName string, val any) *DBFilter {
	result := bson.M{strings.ToLower(fieldName): bson.M{"$eq": val}}
	this.Array = append(this.Array, result)
	return this
}

// > 大于
func (this *DBFilter) Gt(fieldName string, val any) *DBFilter {
	result := bson.M{strings.ToLower(fieldName): bson.M{"$gt": val}}
	this.Array = append(this.Array, result)
	return this
}

// >= 大于等于
func (this *DBFilter) Gte(fieldName string, val any) *DBFilter {
	result := bson.M{strings.ToLower(fieldName): bson.M{"$gte": val}}
	this.Array = append(this.Array, result)
	return this
}

// < 小于
func (this *DBFilter) Lt(fieldName string, val any) *DBFilter {
	result := bson.M{strings.ToLower(fieldName): bson.M{"$lt": val}}
	this.Array = append(this.Array, result)
	return this
}

// <= 小于等于
func (this *DBFilter) Lte(fieldName string, val any) *DBFilter {
	result := bson.M{strings.ToLower(fieldName): bson.M{"$lte": val}}
	this.Array = append(this.Array, result)
	return this
}

// 或 ||
func (this *DBFilter) Or() *DBFilter {
	this.or = []bson.M{}
	for _, value := range this.Array {
		this.or = append(this.or, value)
	}
	this.Array = []bson.M{}
	return this
}

// 并 &&
func (this *DBFilter) And() *DBFilter {
	this.and = []bson.M{}
	for _, value := range this.Array {
		this.and = append(this.and, value)
	}
	this.Array = []bson.M{}
	return this
}
func (this *DBFilter) toOnlyId() bson.M {
	result := bson.M{}
	if this.Array == nil || len(this.Array) == 0 {
		return nil
	} else if len(this.Array) == 1 {
		return this.Array[0]
	} else {
		this.And()
	}
	if this.and != nil {
		result["$and"] = this.and
	}
	return result
}
func (this *DBFilter) ToBson() bson.M {
	result := bson.M{}
	if this.or == nil && this.and == nil {
		if this.Array == nil || len(this.Array) == 0 {
			return bson.M{"_id": bson.M{"$not": bson.M{"$eq": ""}}}
		} else if len(this.Array) == 1 {
			return this.Array[0]
		} else {
			this.And()
		}
	}

	if this.or != nil {
		result["$or"] = this.or
	}
	if this.and != nil {
		result["$and"] = this.and
	}
	return result
}

type DBUpdate struct {
	Array bson.M
}

// 设置字段和值
func (this *DBUpdate) Set(fieldName string, val any) *DBUpdate {
	this.Array[strings.ToLower(fieldName)] = val
	return this
}

// bson.M
func (this *DBUpdate) ToBson() bson.M {
	if len(this.Array) == 0 {
		util.LogError("无任何数据更新")
	}
	result := bson.M{}
	result["$set"] = this.Array
	return result
}
