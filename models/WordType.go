package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func init() {
	orm.RegisterModel(new(WordType))
}

type WordType struct{
	Id int
	Name string `orm:"size(255);null"`
	Parent *WordType `orm:"rel(fk);null"`
}

func InsertWordType(wordType *WordType){
	o := orm.NewOrm()
	_,err := o.Insert(wordType)
	if err!=nil{
		logs.Error(err)
	}

}