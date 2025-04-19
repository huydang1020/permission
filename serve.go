package main

import (
	"log"
	"net"
	"time"

	pb "github.com/huyshop/header/permission"
	"github.com/huyshop/permission/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Permission struct {
	Db IDatabase
}

type IDatabase interface {
	InsertRole(req *pb.Role) error
	GetRole(req *pb.Role) (*pb.Role, error)
	ListRole(req *pb.RoleRequest) ([]*pb.Role, error)
	UpdateRole(req *pb.Role) error
	DeleteRole(req *pb.Role) error
	IsRoleExist(req *pb.Role) (bool, error)
	CountRoles(rq *pb.RoleRequest) (int64, error)
	TranDeleteRole(req *pb.Role) error

	InsertPage(req *pb.Page) (*pb.Page, error)
	GetPage(req *pb.Page) (*pb.Page, error)
	ListPage(req *pb.PageRequest) ([]*pb.Page, error)
	UpdatePage(req *pb.Page) error
	IsPageExist(req *pb.Page) (bool, error)
	DeletePage(req *pb.Page) error
	CountPages(rq *pb.PageRequest) (int64, error)

	InsertGroup(req ...*pb.Group) error
	GetGroup(req *pb.Group) (*pb.Group, error)
	ListGroup(req *pb.Group) ([]*pb.Group, error)
	DeleteGroup(req *pb.Group) error
	UpdateGroup(req ...*pb.Group) error
}

func NewPermisssion(cf *Configs) (*Permission, error) {
	dbase := &db.DB{}
	if err := dbase.ConnectDb(cf.DBPath, cf.DBName); err != nil {
		return nil, err
	}
	log.Println("Connect db successful")
	return &Permission{
		Db: dbase,
	}, nil
}

func startGRPCServe(port string, p *Permission) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 15 * time.Second,
		}),
	}
	serve := grpc.NewServer(opts...)
	pb.RegisterPermissionServiceServer(serve, p)
	reflection.Register(serve)
	return serve.Serve(listen)
}
