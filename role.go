package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/permission"
	"github.com/huyshop/permission/utils"
)

func (p *Permission) CreateRole(ctx context.Context, req *pb.Role) (*common.Empty, error) {
	log.Println("create role:", req)
	if req.GetName() == "" {
		return nil, errors.New(utils.E_not_found_name)
	}
	req.Id = utils.MakeRoleId()
	req.CreatedAt = time.Now().Unix()
	req.State = pb.Page_active.String()
	if len(req.Page) > 0 {
		for _, perm := range req.Page {
			if err := p.Db.TransInsertPageRole(req, perm.PageRole); err != nil {
				log.Println("trans insert pager role err:", err)
				return nil, err
			}
		}
		return &common.Empty{}, nil
	}
	if err := p.Db.InsertRole(req); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Permission) GetRole(ctx context.Context, req *pb.RoleRequest) (*pb.Role, error) {
	log.Println("get role:", req)
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	role, err := p.Db.GetRole(&pb.Role{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (p *Permission) ListRoles(ctx context.Context, req *pb.RoleRequest) (*pb.Roles, error) {
	log.Println("list role:", req)
	roles, err := p.Db.ListRole(req)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return &pb.Roles{}, nil
	}
	count, _ := p.Db.CountRoles(req)
	return &pb.Roles{Roles: roles, Total: count}, nil
}

func (p *Permission) UpdateRole(ctx context.Context, req *pb.Role) (*pb.Role, error) {
	log.Println("update role:", req)
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	req.UpdatedAt = time.Now().Unix()
	err := p.Db.UpdateRole(req)
	if err != nil {
		return nil, err
	}
	role, err := p.Db.GetRole(req)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (p *Permission) DeleteRole(ctx context.Context, req *pb.Role) (*common.Empty, error) {
	log.Println("delete role:", req)
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if err := p.Db.DeleteRole(req); err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *Permission) CheckAccess(ctx context.Context, in *pb.PolicyRequest) (*common.Empty, error) {
	log.Println("check access:", in)
	if in.GetRoleId() == "" || in.GetGroup() == "" || in.GetAction() == "" {
		return nil, errors.New(utils.E_error_invalid_params)
	}
	page, err := p.Db.GetPage(&pb.Page{Group: in.GetGroup()})
	if err != nil {
		log.Println("get page err:", err)
		return nil, err
	}
	if page == nil {
		return nil, errors.New(utils.E_not_found_page)
	}
	pageRole, err := p.Db.GetPageRole(&pb.PageRole{PageId: page.GetId(), RoleId: in.GetRoleId()})
	if err != nil {
		log.Println("get page role err:", err)
		return nil, errors.New(utils.E_access_is_denied)
	}
	if pageRole == nil || pageRole.GetActions() == nil {
		return nil, errors.New(utils.E_access_is_denied)
	}
	for _, act := range pageRole.GetActions() {
		if act == in.GetAction() {
			return &common.Empty{}, nil
		}
	}
	return nil, errors.New(utils.E_access_is_denied)
}
