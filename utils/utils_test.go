package utils

import (
	"log"
	"testing"
)

func TestXxx(t *testing.T) {
	for i := 1; i <= 16; i++ {
		log.Println(MakePageId())
	}
}
