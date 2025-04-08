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
		return nil, errors.New(utils.E_not_found)
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
		return errors.New(utils.E_not_found)
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
		return errors.New(utils.E_not_found)
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
	page := &pb.Page{Id: req.Id}
	b, err := d.engine.Get(page)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, errors.New(utils.E_not_found)
	}
	return page, nil
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
	b, err := d.engine.Exist(&pb.Page{Id: req.Id})
	if err != nil {
		return false, err
	}
	return b, err
}

func (d *DB) UpdatePage(req *pb.Page) error {
	b, err := d.IsPageExist(req)
	if err != nil {
		return err
	}
	if !b {
		return errors.New(utils.E_not_found)
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
		return errors.New(utils.E_not_found)
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

func (d *DB) TranDeletePage(req *pb.Page) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	count, err := sess.Delete(req)
	if err != nil {
		sess.Rollback()
		return errors.New(utils.E_can_not_delete)
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_delete)
	}
	count, err = sess.Delete(&pb.PageRole{PageId: req.Id})
	if err != nil {
		sess.Rollback()
		return errors.New(utils.E_can_not_delete)
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_delete)
	}
	return sess.Commit()
}

func (d *DB) TransInsertPageRole(role *pb.Role, pg *pb.PageRole) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	count, err := sess.Insert(role)
	if err != nil {
		sess.Rollback()
		return errors.New(utils.E_can_not_insert)
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_insert)
	}
	count, err = sess.Insert(pg)
	if err != nil {
		sess.Rollback()
		return errors.New(utils.E_can_not_insert)
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_insert)
	}
	return sess.Commit()
}

func (d *DB) InsertPageRole(req *pb.PageRole) error {
	count, err := d.engine.Insert(req)
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}

func (d *DB) IsPageRoleExist(req *pb.PageRole) (bool, error) {
	b, err := d.engine.Exist(&pb.PageRole{RoleId: req.RoleId, PageId: req.PageId})
	if err != nil {
		return false, err
	}
	return b, err
}

func (d *DB) UpdatePageRole(req *pb.PageRole) error {
	b, err := d.IsPageRoleExist(req)
	if err != nil {
		return err
	}
	if !b {
		return errors.New(utils.E_not_found)
	}
	count, err := d.engine.Update(req, &pb.PageRole{RoleId: req.RoleId, PageId: req.PageId})
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New(utils.E_can_not_update)
	}
	return nil
}

func (d *DB) listPageRoleQuery(req *pb.PageRoleRequest) *xorm.Session {
	ss := d.engine.Table("page_role")
	if req.RoleId != "" {
		ss.And("role_id = ?", req.RoleId)
	}
	if req.PageId != "" {
		ss.And("page_id = ?", req.PageId)
	}
	return ss
}

func (d *DB) ListPageRole(req *pb.PageRoleRequest) ([]*pb.PageRole, error) {
	pr := []*pb.PageRole{}
	ss := d.listPageRoleQuery(req)
	err := ss.Find(&pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (d *DB) GetPageRole(req *pb.PageRole) (*pb.PageRole, error) {
	pageRole := &pb.PageRole{RoleId: req.RoleId, PageId: req.PageId}
	b, err := d.engine.Get(pageRole)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, errors.New(utils.E_not_found)
	}
	return pageRole, nil
}

func (d *DB) TranDelPageRole(req []*pb.PageRole) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	for _, r := range req {
		count, err := sess.Delete(r)
		if err != nil {
			sess.Rollback()
			return errors.New(utils.E_can_not_delete)
		}
		if count < 1 {
			sess.Rollback()
			return errors.New(utils.E_can_not_delete)
		}
	}
	return sess.Commit()
}
