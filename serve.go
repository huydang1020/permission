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

type IDatabase interface{}

func NewPermisssion(cf *Config) (*Permission, error) {
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
