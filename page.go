package main

import (
	"context"
	"errors"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/permission"
)

func (p *Permission) CreatePage(ctx context.Context, req *pb.Page) (*pb.Page, error) {
	if req == nil {
		return nil, errors.New("page is nil")
	}
	if req.GetName() == "" {
		return nil, errors.New("page name is empty")
	}
	page, err := p.Db.InsertPage(req)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) GetPage(ctx context.Context, req *pb.PageRequest) (*pb.Page, error) {
	if req == nil {
		return nil, errors.New("page is nil")
	}
	if req.GetId() == "" {
		return nil, errors.New("page id is empty")
	}
	page, err := p.Db.GetPage(&pb.Page{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) ListPages(ctx context.Context, req *pb.PageRequest) (*pb.Pages, error) {
	if req == nil {
		return nil, errors.New("page is nil")
	}
	pages, err := p.Db.ListPage(req)
	if err != nil {
		return nil, err
	}
	count, _ := p.Db.CountPages(req)
	return &pb.Pages{Pages: pages, Total: count}, nil
}

func (p *Permission) ListPage(ctx context.Context, req *pb.PageRequest) ([]*pb.Page, error) {
	pages, err := p.Db.ListPage(req)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (p *Permission) UpdatePage(ctx context.Context, req *pb.Page) (*pb.Page, error) {
	if req == nil {
		return nil, errors.New("page is nil")
	}
	if req.GetId() == "" {
		return nil, errors.New("page id is empty")
	}
	err := p.Db.UpdatePage(req)
	if err != nil {
		return nil, err
	}
	page, err := p.Db.GetPage(req)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) DeletePage(ctx context.Context, req *pb.Page) (*common.Empty, error) {
	if req == nil {
		return nil, errors.New("page is nil")
	}
	if req.GetId() == "" {
		return nil, errors.New("page id is empty")
	}
	return nil, nil
}
