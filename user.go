package main

import (
	"context"
	"errors"

	pb "github.com/huyshop/header/permission"
	"github.com/huyshop/permission/utils"
)

func (p *Permission) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if req.GetFullName() == "" {
		return nil, errors.New(utils.E_not_found_name)
	}
	page, err := p.Db.InsertUser(req)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) GetUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	page, err := p.Db.GetUser(&pb.User{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return page, nil
}

