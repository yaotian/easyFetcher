package main

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beedb"
	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var orm beedb.Model

type Userinfo struct {
	Uid        int `PK`
	Username   string
	Departname string
	Created    string
}

func init_db() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db)
}

func insert() {
	//save data
	var saveone Userinfo
	saveone.Username = "Test Add User"
	saveone.Departname = "Test Add Departname"
	saveone.Created = time.Now().Format("2006-01-02 15:04:05")
	orm.Save(&saveone)
	fmt.Println(saveone)
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	insert()
	this.Ctx.WriteString("hello world")
}

func main() {
	init_db()
	beego.RegisterController("/", &MainController{})
	beego.HttpPort = 9008
	beego.Run()
}
