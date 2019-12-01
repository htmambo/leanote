package tests

import (
	"github.com/htmambo/leanote/app/db"
	"testing"
	//	. "github.com/htmambo/leanote/app/lea"
	//	"github.com/htmambo/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func TestDBConnect(t *testing.T) {
	db.Init("mongodb://localhost:27017/leanote", "leanote")
}
