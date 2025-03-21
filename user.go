package main

import (
	"context"
	"errors"

	pb "github.com/huyshop/header/permission"
)

func (p *Permission) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if req == nil {
		return nil, errors.New("page is nil")
	}
	if req.GetFullName() == "" {
		return nil, errors.New("page name is empty")
	}
	page, err := p.Db.InsertUser(req)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) GetUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if req == nil {
		return nil, errors.New("page is nil")
	}
	if req.GetId() == "" {
		return nil, errors.New("page id is empty")
	}
	page, err := p.Db.GetUser(&pb.User{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return page, nil
}

