package utils

import "github.com/rs/xid"

func MakeRoleId() string {
	return "role" + xid.New().String()
}

func MakePageId() string {
	return "page" + xid.New().String()
}
