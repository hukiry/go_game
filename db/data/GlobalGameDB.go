package data

import (
	"game/db/mongodb"
)

type GlobalGame struct {
	Flag      uint32 `flag`
	UserCount uint32
}

const FLAG = 1001

func (this *GlobalGame) ReadDB() {
	db := mongodb.DBConfig[GlobalGame]{}
	db.BuildFilter().Eq("Flag", FLAG)
	oldData := db.LoadTable(GlobalGameTable).Find().One()
	if oldData == nil {
		this.Flag = FLAG
		db.LoadTable(GlobalGameTable).Insert(*this)
	} else {
		this = oldData
	}
}

func (this *GlobalGame) CreateUser() uint32 {
	this.UserCount++
	return this.UserCount
}

func (this *GlobalGame) saveDB() {
	db := mongodb.DBConfig[GlobalGame]{}
	db.BuildFilter().Eq("Flag", FLAG)
	db.BuildUpdate().Set("UserCount", this.UserCount)
	db.LoadTable(GlobalGameTable).Update()
}
