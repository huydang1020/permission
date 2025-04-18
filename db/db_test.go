package db

import (
	"log"
	"testing"

	"github.com/huyshop/header/permission"
)

func TestXxx(t *testing.T) {
	d := DB{}
	d.ConnectDb("root:hanhuquynh@tcp(localhost:3306)", "permission")
	pages, err := d.ListPage(&permission.PageRequest{})
	if err != nil {
		t.Fail()
	}
	action := []string{"c", "r", "u", "d"}
	count := 0
	for _, page := range pages {
		if err := d.InsertPageRole(&permission.PageRole{
			RoleId:  "rolecvgi1navufdq3at642m0",
			PageId:  page.GetId(),
			Actions: action,
		}); err != nil {
			log.Println("err :", err, page.GetId())
			continue
		}
		count++
	}
	log.Println("done:", count)
}
