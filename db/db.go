package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/huyshop/header/permission"
	"github.com/huyshop/permission/utils"
	"xorm.io/xorm"
)

type DB struct {
	engine *xorm.Engine
}

func (d *DB) ConnectDb(sqlPath, dbName string) error {
	sqlConnStr := fmt.Sprintf("%s/%s", sqlPath, dbName)
	engine, err := xorm.NewEngine("mysql", sqlConnStr)
	if err != nil {
		return err
	}
	tickPingSql := time.NewTicker(15 * time.Minute)
	go func() {
		for {
			select {
			case <-tickPingSql.C:
				if err := engine.Ping(); err != nil {
					log.Print("sql can not ping")
				}
			}
		}
	}()
	d.engine = engine
	d.engine.ShowSQL(false)
	return err
}

func (d *DB) InsertRole(req *pb.Role) error {
	count, err := d.engine.Insert(req)
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}

func (d *DB) GetRole(req *pb.Role) (*pb.Role, error) {
	role := &pb.Role{Id: req.Id}
	b, err := d.engine.Get(role)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, errors.New(utils.E_not_found_role)
	}
	return role, nil
}

func (d *DB) listRoleQuery(req *pb.RoleRequest) *xorm.Session {
	ss := d.engine.Table("role")
	if req.GetId() != "" {
		ss.And("id = ?", req.GetId())
	}
	if req.GetName() != "" {
		ss.And("name LIKE ?", "%"+req.GetName()+"%")
	}
	return ss
}

func (d *DB) ListRole(req *pb.RoleRequest) ([]*pb.Role, error) {
	roles := []*pb.Role{}
	ss := d.listRoleQuery(req)
	err := ss.Desc("created_at").Find(&roles)
	if err != nil {
		log.Println("get list:", err)
		return nil, err
	}
	return roles, nil
}

func (d *DB) CountRoles(rq *pb.RoleRequest) (int64, error) {
	ss := d.listRoleQuery(rq)
	return ss.Count()
}

func (d *DB) IsRoleExist(req *pb.Role) (bool, error) {
	b, err := d.engine.Exist(&pb.Role{Id: req.Id})
	if err != nil {
		return false, err
	}
	return b, err
}

func (d *DB) UpdateRole(req *pb.Role) error {
	b, err := d.IsRoleExist(req)
	if err != nil {
		return err
	}
	if !b {
		return errors.New(utils.E_not_found_role)
	}
	count, err := d.engine.Update(req, &pb.Role{Id: req.Id})
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_update)
	}
	return nil
}

func (d *DB) DeleteRole(req *pb.Role) error {
	b, err := d.IsRoleExist(req)
	if err != nil {
		return err
	}
	if !b {
		return errors.New(utils.E_not_found_role)
	}
	count, err := d.engine.Delete(req)
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}

func (d *DB) TranDeleteRole(req *pb.Role) error {
	ss := d.engine.NewSession()
	defer ss.Close()

	if err := ss.Begin(); err != nil {
		return err
	}
	_, err := ss.Table(tblGroup).Delete(&pb.Group{RoleId: req.GetId()})
	if err != nil {
		ss.Rollback()
		log.Println("delete error:", err)
		return err
	}
	_, err = ss.Table(tblRole).Delete(&pb.Role{Id: req.GetId()})
	if err != nil {
		ss.Rollback()
		log.Println("delete error:", err)
		return err
	}
	if err := ss.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *DB) InsertPage(req *pb.Page) (*pb.Page, error) {
	count, err := d.engine.Insert(req)
	if err != nil {
		return nil, err
	}
	if count < 1 {
		return nil, errors.New(utils.E_can_not_insert)
	}
	return req, nil
}

func (d *DB) GetPage(req *pb.Page) (*pb.Page, error) {
	b, err := d.engine.Get(req)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, errors.New(utils.E_not_found_page)
	}
	return req, nil
}

func (d *DB) listPageQuery(req *pb.PageRequest) *xorm.Session {
	ss := d.engine.Table("page")
	if len(req.GetIds()) > 0 {
		ss.In("id", req.GetIds())
	} else if req.GetId() != "" {
		ss.And("id = ?", req.GetId())
	}
	if req.GetName() != "" {
		ss.And("name LIKE ?", "%"+req.GetName()+"%")
	}
	if req.GetRoute() != "" {
		ss.And("route LIKE ?", "%"+req.GetRoute()+"%")
	}
	if req.GetType() != "" {
		ss.And("type = ?", req.GetType())
	}
	if req.GetParentId() != "" {
		ss.And("parent_id = ?", req.GetParentId())
	}
	return ss
}

func (d *DB) ListPage(req *pb.PageRequest) ([]*pb.Page, error) {
	pages := []*pb.Page{}
	log.Println("list page:", req)
	ss := d.listPageQuery(req)
	if req.GetLimit() != 0 {
		ss.Limit(int(req.GetLimit()), int(req.GetLimit())*int(req.GetSkip()))
	}
	err := ss.Desc("created_at").Find(&pages)
	if err != nil {
		log.Println("get list:", err)
		return nil, err
	}
	return pages, nil
}

func (d *DB) CountPages(rq *pb.PageRequest) (int64, error) {
	ss := d.listPageQuery(rq)
	return ss.Count()
}

func (d *DB) IsPageExist(req *pb.Page) (bool, error) {
	return d.engine.Exist(&pb.Page{Id: req.Id})
}

func (d *DB) UpdatePage(req *pb.Page) error {
	b, err := d.IsPageExist(req)
	if err != nil {
		return err
	}
	if !b {
		return errors.New(utils.E_not_found_page)
	}
	count, err := d.engine.Update(req, &pb.Page{Id: req.Id})
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_update)
	}
	return nil
}

func (d *DB) DeletePage(req *pb.Page) error {
	b, err := d.IsPageExist(req)
	if err != nil {
		return err
	}
	if !b {
		return errors.New(utils.E_not_found_page)
	}
	count, err := d.engine.Delete(req)
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}

func (d *DB) InsertGroup(req ...*pb.Group) error {
	count, err := d.engine.Insert(req)
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}

func (d *DB) GetGroup(req *pb.Group) (*pb.Group, error) {
	b, err := d.engine.Get(req)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, errors.New(utils.E_not_found_group)
	}
	return req, nil
}

func (d *DB) listGroupQuery(req *pb.Group) *xorm.Session {
	ss := d.engine.Table("group")
	if req.GetRoleId() != "" {
		ss.And("role_id = ?", req.GetRoleId())
	}
	if req.GetGroup() != "" {
		ss.And("group = ?", req.GetGroup())
	}
	return ss
}

func (d *DB) ListGroup(req *pb.Group) ([]*pb.Group, error) {
	groups := []*pb.Group{}
	ss := d.listGroupQuery(req)
	err := ss.Find(&groups)
	if err != nil {
		log.Println("get list err:", err)
		return nil, err
	}
	return groups, nil
}

func (d *DB) DeleteGroup(req *pb.Group) error {
	count, err := d.engine.Delete(req)
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}

func (d *DB) UpdateGroup(req ...*pb.Group) error {
	for _, gr := range req {
		count, err := d.engine.Update(gr, &pb.Group{RoleId: gr.GetRoleId(), Group: gr.GetGroup()})
		if err != nil {
			return err
		}
		if count < 1 {
			log.Println("update group err:", utils.E_can_not_update)
			return nil
		}
	}
	return nil
}

func (d *DB) IsGroupExist(req *pb.Group) (bool, error) {
	b, err := d.engine.Exist(&pb.Group{RoleId: req.GetRoleId()})
	if err != nil {
		return false, err
	}
	return b, nil
}
