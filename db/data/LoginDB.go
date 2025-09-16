package data

import (
	"game/db/mongodb"
	"game/util"
	"game/util/timer"
)

const (
	Login_LangCode  = iota //语言切换代码
	Login_Platform         // 登录平台
	Login_Nick             // 昵称
	Login_IpAddress        // ip地址
	Login_Type             // 登录类型：游客或者账号
	Login_HeadID           // 头像ID
	Login_GassID           // 社团或者组队 id
)

type User struct {
	RoleID   uint32 `roleid` //角色ID
	DeviceId string //首次登录的设备id
	Token    string //账号token 或者 邮箱名

	SaveTime   uint32          //最后保存时间
	CreateTime uint32          //创建的时间
	ItemMoney  map[byte]uint32 //货币集合
	ItemPack   map[byte]string //0语言，1平台，2昵称，3地址, 4登录类型，5头像ID，6组队id
}

func (this *User) ReadDB() {
	/*切换账号
	1，根据设备ID，查找数据。
		<1>，如果没有找到，那么根据token，查找数据。
			1，如没有找到，那么就是新号。
				创建新号
			2, 如果找到, 那么就是切换账号.
				切换账号：更新设备ID，通知另外一个在线设备
	*/
	db := mongodb.DBConfig[User]{}
	db.BuildFilter().Eq("DeviceId", this.DeviceId)
	table := db.LoadTable(LoginTable)
	oldData := table.Find().One()
	if oldData == nil {
		if this.Token == "" {
			this.createData()
			return
		}

		db.BuildFilter().Eq("Token", this.Token)
		oldData = table.Find().One()
		if oldData == nil {
			this.createData()
			return
		} else {
			oldData.DeviceId = this.DeviceId
			build := db.BuildUpdate()
			build.Set("DeviceId", this.DeviceId)
			build.Set("SaveTime", uint32(timer.GetLocalTime().Unix()))
			table.Update()
			//todo 通知另外一个在线设备
		}
	}

	oldData.ItemPack[Login_IpAddress] = this.ItemPack[Login_IpAddress]
	oldData.ItemPack[Login_LangCode] = this.ItemPack[Login_LangCode]
	oldData.ItemPack[Login_Type] = this.ItemPack[Login_Type]
	oldData.ItemPack[Login_Platform] = this.ItemPack[Login_Platform]

	*this = *oldData
}

func (this *User) BindDB() {
	db := mongodb.DBConfig[User]{}
	db.BuildFilter().Eq("RoleID", this.RoleID)
	build := db.BuildUpdate()
	build.Set("Token", this.Token)
	build.Set("SaveTime", uint32(timer.GetLocalTime().Unix()))
	db.LoadTable(LoginTable).Update()
}

// 创建新号
func (this *User) createData() {
	//创建角色id
	global := &GlobalGame{}
	global.ReadDB()
	this.RoleID = global.CreateUser()
	global.saveDB()

	//创建基本数据
	this.ItemPack[Login_HeadID] = "1"
	for i := 1; i <= 24; i++ {
		this.ItemMoney[byte(i)] = 0
	}
	this.CreateTime = uint32(timer.GetLocalTime().Unix())
	//添加到数据库
	db := mongodb.DBConfig[User]{}
	db.BuildFilter().Eq("RoleID", this.RoleID)
	db.LoadTable(LoginTable).Insert(*this)
}

func (this *User) SaveDB() {
	if this.RoleID == 0 {
		util.Log("RoleID is not 0")
	}
	db := mongodb.DBConfig[User]{}
	db.BuildFilter().Eq("RoleID", this.RoleID)
	build := db.BuildUpdate()
	build.Set("ItemMoney", this.ItemMoney)
	build.Set("ItemPack", this.ItemPack)
	build.Set("SaveTime", this.SaveTime)
	db.LoadTable(LoginTable).Update()
}
