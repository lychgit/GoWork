package models

import (
	"github.com/astaxie/beego/orm"
		)

type User struct {
	Id        	int  		`orm:"auto"`
	UserName  	string
	Password  	string
	Salt      	string
	Email     	string
	LastLogin 	int64
	LastIp    	string
	Status    	int
	Level		int
}

func (u *User) TableName() string {
	return TableName("user")
}

func (u *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func (u *User) Add() error {
	if _, err := orm.NewOrm().Insert(u); err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error{
	if _, err := orm.NewOrm().Delete(u); err != nil {
		return err
	}
	return nil
}

func UserAdd(user *User) (int64, error) {
	return orm.NewOrm().Insert(user)
}

func UserUpdate(user *User, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}

func UserGetById(id int) (*User, error) {
	u := new(User)
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserGetByName(userName string) (*User, error) {
	u := new(User)
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("user_name", userName).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

