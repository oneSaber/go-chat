// 链接数据库,创建模型, 新建数据库表

package data

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Id        int    `gorm:"primary_key"`
	Account   string `gorm:"index"`
	Password  string
	Email     string
	Name      string
	Signature string // 个性签名
	Avatar    string // 头像 上传url
}

// 主贴
type Aritcle struct {
	gorm.Model
	Name     string `gorm:"priamry_key"`
	content  string
	AuthorId int
	Author   User `gorm:"ForeignKey:AuthorId;AssociationForeignKey:Id"`
}

// 更贴
type Post struct {
	gorm.Model
	content    string
	AuthorId   int
	Author     User `gorm:"ForeignKey:AuthorId;AssociationForeignKey:Id"`
	AritcleId  int
	MainAritce Aritcle `gorm:"ForeignKey:AritcleId;AssociationForeignKey:Id"`
}

var database *gorm.DB
var redisClient *redis.Client

func init() {

	// 配置数据库
	var err error
	database, err = gorm.Open("mysql", "root:123456@(39.105.64.7:3306)/Myblog?charset=utf8&parserTime&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	// 配置自动迁移
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Aritcle{})
	database.AutoMigrate(&Post{})

	// 创建表
	if !database.HasTable(&User{}) {
		database.CreateTable(&User{})
	}

	if !database.HasTable(&Aritcle{}) {
		database.CreateTable(&Aritcle{})
	}

	if !database.HasTable(&Post{}) {
		database.CreateTable(&Post{})
	}

	// 配置redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "39.105.64.7: 6379",
		Password: "",
		DB:       0,
	})
	return

}
