package main

import (
	"context"
	"errors"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/permission"
	"github.com/huyshop/permission/utils"
)

func (p *Permission) CreateRole(ctx context.Context, req *pb.Role) (*pb.Role, error) {
	if req == nil {
		return nil, errors.New(utils.E_not_found)
	}
	if req.GetName() == "" {
		return nil, errors.New("role name is empty")
	}
	role, err := p.Db.InsertRole(req)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (p *Permission) GetRole(ctx context.Context, req *pb.RoleRequest) (*pb.Role, error) {
	if req == nil {
		return nil, errors.New(utils.E_not_found)
	}
	if req.GetId() == 0 {
		return nil, errors.New("role id is empty")
	}
	role, err := p.Db.GetRole(&pb.Role{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (p *Permission) ListRoles(ctx context.Context, req *pb.RoleRequest) (*pb.Roles, error) {
	if req == nil {
		return nil, errors.New(utils.E_not_found)
	}
	roles, err := p.Db.ListRole(req)
	if err != nil {
		return nil, err
	}
	count, _ := p.Db.CountRoles(req)
	return &pb.Roles{Roles: roles, Total: count}, nil
}

func (p *Permission) UpdateRole(ctx context.Context, req *pb.Role) (*pb.Role, error) {
	if req == nil {
		return nil, errors.New(utils.E_not_found)
	}
	if req.GetId() == 0 {
		return nil, errors.New("role id is empty")
	}
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
	if req == nil {
		return nil, errors.New(utils.E_not_found)
	}
	if req.GetId() == 0 {
		return nil, errors.New("role id is empty")
	}
	if err := p.Db.DeleteRole(req); err != nil {
		return nil, err
	}
	return nil, nil
}
