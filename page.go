package main

import (
	"context"
	"errors"
	"time"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/permission"
	"github.com/huyshop/permission/utils"
)

func (p *Permission) CreatePage(ctx context.Context, req *pb.Page) (*pb.Page, error) {
	if req.GetName() == "" {
		return nil, errors.New(utils.E_not_found_name)
	}
	req.Id = utils.MakePageId()
	req.CreatedAt = time.Now().Unix()
	req.Status = int32(pb.Page_active)
	page, err := p.Db.InsertPage(req)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) GetPage(ctx context.Context, req *pb.PageRequest) (*pb.Page, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	page, err := p.Db.GetPage(&pb.Page{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	// log.Println("page", page)
	if page.ParentId != "" {
		chid, err := p.Db.ListPage(&pb.PageRequest{ParentId: page.GetParentId()})
		if err != nil {
			return nil, err
		}
		page.Children = chid
	}
	return page, nil
}

func (p *Permission) ListPages(ctx context.Context, req *pb.PageRequest) (*pb.Pages, error) {

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
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	req.UpdatedAt = time.Now().Unix()
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
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	err := p.Db.DeletePage(&pb.Page{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
