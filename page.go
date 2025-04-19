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

func (p *Permission) CreatePage(ctx context.Context, req *pb.Page) (*pb.Page, error) {
	log.Println("create page:", req)
	if req.GetPath() == "" {
		return nil, errors.New(utils.E_not_found_path)
	}
	req.Id = utils.MakePageId()
	req.CreatedAt = time.Now().Unix()
	req.State = pb.Page_active.String()
	page, err := p.Db.InsertPage(req)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) GetPage(ctx context.Context, req *pb.PageRequest) (*pb.Page, error) {
	log.Println("get page:", req)
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
	log.Println("list page:", req)
	var pages []*pb.Page
	var err error
	pages, err = p.Db.ListPage(req)
	if err != nil {
		return nil, err
	}
	if req.GetRoleId() != "" {
		newPages := make([]*pb.Page, 0, len(pages))
		for _, page := range pages {
			if slices.Contains(page.GetRoleIds(), req.GetRoleId()) {
				newPages = append(newPages, page)
			}
		}
		pages = newPages
	}
	count, err := p.Db.CountPages(req)
	if err != nil {
		return nil, err
	}
	return &pb.Pages{Pages: pages, Total: count}, nil
}

func (p *Permission) UpdatePage(ctx context.Context, req *pb.Page) (*pb.Page, error) {
	log.Println("update page:", req)
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	req.UpdatedAt = time.Now().Unix()
	err := p.Db.UpdatePage(req)
	if err != nil {
		return nil, err
	}
	// mapOldPR := map[string]*pb.PageRole{}
	// listpr, err := p.Db.ListPageRole(&pb.PageRoleRequest{PageId: req.GetId()})
	// if err != nil {
	// 	log.Println("list page role err:", err)
	// 	return nil, err
	// }
	// for _, pr := range listpr {
	// 	mapOldPR[pr.GetRoleId()] = pr
	// }
	// mapNewPR := map[string]*pb.PageRole{}
	// mapUpdatePR := map[string]*pb.PageRole{}
	// for _, pr := range req.GetRoleActions() {
	// 	pr.PageId = req.GetId()
	// 	if pr.GetRoleId() == "" {
	// 		return nil, errors.New(utils.E_not_found_role_id)
	// 	}
	// 	if _, ok := mapOldPR[pr.GetRoleId()]; ok {
	// 		mapUpdatePR[pr.GetRoleId()] = pr
	// 		delete(mapOldPR, pr.GetRoleId())
	// 	} else {
	// 		mapNewPR[pr.GetRoleId()] = pr
	// 	}
	// }
	// if len(mapNewPR) > 0 {
	// 	for _, pr := range mapNewPR {
	// 		if err := p.Db.InsertPageRole(pr); err != nil {
	// 			log.Println("insert page role err:", err)
	// 			return nil, err
	// 		}
	// 	}
	// }
	// if len(mapOldPR) > 0 {
	// 	for _, pr := range mapOldPR {
	// 		if err := p.Db.TranDelPageRole([]*pb.PageRole{pr}); err != nil {
	// 			log.Println("delete page role err:", err)
	// 			return nil, err
	// 		}
	// 	}
	// }
	// if len(mapUpdatePR) > 0 {
	// 	for _, pr := range mapUpdatePR {
	// 		if err := p.Db.UpdatePageRole(pr); err != nil {
	// 			log.Println("update page role err:", err)
	// 			return nil, err
	// 		}
	// 	}
	// }
	page, err := p.Db.GetPage(req)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p *Permission) DeletePage(ctx context.Context, req *pb.Page) (*common.Empty, error) {
	log.Println("delete page:", req)
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	err := p.Db.DeletePage(&pb.Page{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
