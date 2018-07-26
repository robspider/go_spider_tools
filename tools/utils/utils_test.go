package utils

import (
	"testing"
	"io/ioutil"
	"log"
	"strings"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
)



func TestFile(t *testing.T)  {
	type vvv struct {
		_id string `json:"_id"`
		totalValue string `json:"totalValue"`
	}
	fileName := "C:/Users/87560/Desktop/result.txt"
	data,_ := ioutil.ReadFile(fileName)
	content := string(data)
	ss  := strings.Split(content,"\r\n")
	log.Println(len(ss))
	for _,c := range ss{
		log.Println(gjson.Get(c,"_id" ).String())
		log.Println(gjson.Get(c,"totalValue" ).String())
	}

	session,_ := mgo.Dial("mongodb://127.0.0.1:27017")
	db := session.DB("")
	query := db.C("").Find("")
}
