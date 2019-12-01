package tests

import (
    // "github.com/htmambo/leanote/app/db"
    "testing"
    //  . "github.com/htmambo/leanote/app/lea"
    // "github.com/htmambo/leanote/app/service"
    //  "gopkg.in/mgo.v2"
    //  "fmt"
    "regexp"
)


// 测试登录
func TestReg(t *testing.T) {
    a := `<a href="http://localhost:9000/api/file/getAttach?fileId=585e0e9c270a35609300000c" target="_blank" data-mce-href="http://localhost:9000/api/file/getAttach?fileId=585e0e9c270a35609300000c">proxy.go</a>`
    reg, _ := regexp.Compile(`"https*://[^/]*?/api/file/getAttach\?fileId=585e0e9c270a35609300000c`)

    a2 := reg.ReplaceAllString(a, `"`)
    t.Log(a2)
}





