package main

import (
	"context"
	"errors"
	"log"
	"slices"
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
	if len(req.GetGroups()) > 0 {
		for _, group := range req.GetGroups() {
			group.RoleId = req.GetId()
		}
		if err := p.Db.InsertGroup(req.GetGroups()...); err != nil {
			return nil, err
		}
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
	for _, role := range roles {
		groups, err := p.Db.ListGroup(&pb.Group{RoleId: role.GetId()})
		if err != nil {
			continue
		}
		role.Groups = groups
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
	mapOldGR := map[string]*pb.Group{}
	listGroup, err := p.Db.ListGroup(&pb.Group{RoleId: req.GetId()})
	if err != nil {
		log.Println("list page role err:", err)
		return nil, err
	}
	for _, gr := range listGroup {
		mapOldGR[gr.GetGroup()] = gr
	}
	mapNewGR := map[string]*pb.Group{}
	mapUpdateGR := map[string]*pb.Group{}
	for _, gr := range req.GetGroups() {
		gr.RoleId = req.GetId()
		if _, ok := mapOldGR[gr.GetGroup()]; ok {
			mapUpdateGR[gr.GetGroup()] = gr
			delete(mapOldGR, gr.GetGroup())
		} else {
			mapNewGR[gr.GetGroup()] = gr
		}
	}
	if len(mapNewGR) > 0 {
		for _, gr := range mapNewGR {
			if err := p.Db.InsertGroup(gr); err != nil {
				log.Println("insert page role err:", err)
				return nil, err
			}
		}
	}
	if len(mapOldGR) > 0 {
		for _, gr := range mapOldGR {
			if err := p.Db.DeleteGroup(gr); err != nil {
				log.Println("delete page role err:", err)
				return nil, err
			}
		}
	}
	if len(mapUpdateGR) > 0 {
		for _, pr := range mapUpdateGR {
			if err := p.Db.UpdateGroup(pr); err != nil {
				log.Println("update page role err:", err)
				return nil, err
			}
		}
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
	if err := p.Db.TranDeleteRole(req); err != nil {
		return nil, errors.New(utils.E_can_not_delete)
	}
	return nil, nil
}

func (p *Permission) CheckAccess(ctx context.Context, in *pb.PolicyRequest) (*common.Empty, error) {
	log.Println("check access:", in)
	if in.GetRoleId() == "" || in.GetGroup() == "" || in.GetAction() == "" {
		return nil, errors.New(utils.E_error_invalid_params)
	}
	group, err := p.Db.GetGroup(&pb.Group{RoleId: in.GetRoleId(), Group: in.GetGroup()})
	if err != nil {
		log.Println("get page err:", err)
		return nil, errors.New(utils.E_access_is_denied)
	}
	if !slices.Contains(group.GetActions(), in.GetAction()) {
		return nil, errors.New(utils.E_access_is_denied)
	}
	return &common.Empty{}, nil
}
