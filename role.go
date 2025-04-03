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
	if req.GetName() == "" {
		return nil, errors.New(utils.E_not_found_name)
	}
	req.Id = utils.MakeRoleId()
	req.CreatedAt = time.Now().Unix()
	req.State = pb.Page_active.String()
	if len(req.Permission) > 0 {
		for _, perm := range req.Permission {
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
	log.Println("req: ", req)
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
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if err := p.Db.DeleteRole(req); err != nil {
		return nil, err
	}
	return nil, nil
}
