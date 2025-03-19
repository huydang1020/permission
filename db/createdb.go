package db

import (
	"log"

	pb "github.com/huyshop/header/permission"

	"xorm.io/xorm"
)

const (
	tblUser = "user"
	tblRole = "role"
	tblPage = "page"
)

func createTable(model interface{}, tblName string, engine *xorm.Engine) error {
	b, err := engine.IsTableExist(model)
	if err != nil {
		return err
	}
	log.Print(b, " ", tblName)
	if b {
		if err = engine.Sync2(model); err != nil {
			return err
		}
		return nil
	}
	if !b {
		if err := engine.CreateTables(model); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func (d *DB) CreateDb() error {
	if err := createTable(&pb.Role{}, tblUser, d.engine); err != nil {
		return err
	}
	if err := createTable(&pb.Page{}, tblPage, d.engine); err != nil {
		return err
	}
	return nil
}
