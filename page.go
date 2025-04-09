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

func (p *Permission) CreatePage(ctx context.Context, req *pb.Page) (*pb.Page, error) {
	if req.GetName() == "" {
		return nil, errors.New(utils.E_not_found_name)
	}
	req.Id = utils.MakePageId()
	req.CreatedAt = time.Now().Unix()
	req.State = pb.Page_active.String()
	page, err := p.Db.InsertPage(req)
	if err != nil {
		return nil, err
	}
	if req.GetRoleActions() != nil {
		for _, pr := range req.GetRoleActions() {
			pr.PageId = page.GetId()
			if pr.GetRoleId() == "" {
				return nil, errors.New(utils.E_not_found_role_id)
			}
			if err := p.Db.InsertPageRole(pr); err != nil {
				log.Println("insert page role err:", err)
				return nil, err
			}
		}
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
	log.Println("list page:", req)
	roles, err := p.Db.ListRole(&pb.RoleRequest{Id: req.GetRoleId()})
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	mapRoles := map[string]*pb.Role{}
	for _, r := range roles {
		mapRoles[r.GetId()] = r
	}
	pageRoles, err := p.Db.ListPageRole(&pb.PageRoleRequest{RoleId: req.GetRoleId()})
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	// pageIds := []string{}
	for _, pr := range pageRoles {
		// pageIds = append(pageIds, pr.GetPageId())
		if r, ok := mapRoles[pr.GetRoleId()]; ok {
			pr.Role = &pb.Role{
				Id:          r.GetId(),
				Name:        r.GetName(),
				Description: r.GetDescription(),
				State:       r.GetState(),
			}
		}
	}
	// req.Ids = pageIds
	pages, err := p.Db.ListPage(req)
	if err != nil {
		return nil, err
	}
	if len(pageRoles) > 0 {
		for _, page := range pages {
			for _, pr := range pageRoles {
				if page.GetId() == pr.GetPageId() {
					page.RoleActions = append(page.RoleActions, &pb.PageRole{RoleId: pr.GetRoleId(), Actions: pr.GetActions(), Role: pr.GetRole()})
				}
			}
		}
	}
	if len(pages) == 0 {
		return &pb.Pages{Pages: pages, Total: 0}, nil
	}
	count, _ := p.Db.CountPages(req)
	return &pb.Pages{Pages: pages, Total: count}, nil
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
	mapOldPR := map[string]*pb.PageRole{}
	listpr, err := p.Db.ListPageRole(&pb.PageRoleRequest{PageId: req.GetId()})
	if err != nil {
		log.Println("list page role err:", err)
		return nil, err
	}
	for _, pr := range listpr {
		mapOldPR[pr.GetRoleId()] = pr
	}
	mapNewPR := map[string]*pb.PageRole{}
	mapUpdatePR := map[string]*pb.PageRole{}
	for _, pr := range req.GetRoleActions() {
		pr.PageId = req.GetId()
		if pr.GetRoleId() == "" {
			return nil, errors.New(utils.E_not_found_role_id)
		}
		if _, ok := mapOldPR[pr.GetRoleId()]; ok {
			mapUpdatePR[pr.GetRoleId()] = pr
			delete(mapOldPR, pr.GetRoleId())
		} else {
			mapNewPR[pr.GetRoleId()] = pr
		}
	}
	if len(mapNewPR) > 0 {
		for _, pr := range mapNewPR {
			if err := p.Db.InsertPageRole(pr); err != nil {
				log.Println("insert page role err:", err)
				return nil, err
			}
		}
	}
	if len(mapOldPR) > 0 {
		for _, pr := range mapOldPR {
			if err := p.Db.TranDelPageRole([]*pb.PageRole{pr}); err != nil {
				log.Println("delete page role err:", err)
				return nil, err
			}
		}
	}
	if len(mapUpdatePR) > 0 {
		for _, pr := range mapUpdatePR {
			if err := p.Db.UpdatePageRole(pr); err != nil {
				log.Println("update page role err:", err)
				return nil, err
			}
		}
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
	err := p.Db.TranDeletePage(&pb.Page{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
