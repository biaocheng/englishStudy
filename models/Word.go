package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func init() {
	orm.RegisterModel(new(Word))
}

type Word struct{
	Id int
	Word string `orm:"size(255)"`
	Phonetic string `orm:"size(255);null"`
	Interpretation string `orm:"size(255);null"`
	MusicUrl string `orm:"size(100);null"`
	WordType *WordType `orm:"rel(fk)"`
}
func InsertWord(word *Word){
	o := orm.NewOrm()
	_,err := o.Insert(word)
	if err!=nil{
		logs.Error(err)
	}

}