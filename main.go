package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/buger/jsonparser"
)

// CheckErr to handle errors
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetGems makes a GET request to the API
func GetGems(uri string) (b []byte) {
	// var a apikeys
	// var secret string
	// var access string
	// a.getAPIKeys(".secrets.yaml")
	// secret = a.Secret
	// access = a.Access
	//log.Printf("Requesting %v", uri)
	// keys := "accessKey=" + access + "; secretKey=" + secret + ";"
	c := &http.Client{}
	r, err := http.NewRequest("GET", "https://raw.githubusercontent.com/brather1ng/RePoE/master/data/"+uri, nil)
	CheckErr(err)
	// r.Header.Add("X-ApiKeys", keys)
	// r.Header.Add("X-Impersonate", "username=zberg@indeed.com")
	resp, err := c.Do(r)
	CheckErr(err)
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	CheckErr(err)
	return b
}

// // APIPost sends a POST to the API
// func APIPost(uri string, bodyContent *bytes.Reader) (b []byte) {
// 	var a apikeys
// 	var secret string
// 	var access string
// 	a.getAPIKeys(".secrets.yaml")
// 	secret = a.Secret
// 	access = a.Access
// 	//log.Printf("Requesting %v", uri)
// 	keys := "accessKey=" + access + "; secretKey=" + secret + ";"
// 	c := &http.Client{}
// 	r, err := http.NewRequest("POST", "https://cloud.tenable.com"+uri, bodyContent)
// 	CheckErr(err)
// 	r.Header.Add("X-ApiKeys", keys)
// 	r.Header.Add("X-Impersonate", "username=zberg@indeed.com")
// 	resp, err := c.Do(r)
// 	CheckErr(err)
// 	defer resp.Body.Close()
// 	b, err = ioutil.ReadAll(resp.Body)
// 	CheckErr(err)
// 	return b
// }

type Classes []struct {
	Class        string   `json:"class"`
	Ascendencies []string `json:"ascendencies"`
}

func main() {
	var b []byte
	var classes Classes

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.Println("Starting Stuff")
	classes.Class = "Duelist"
	// classes.Class.Ascendencies = {"Slayer", "stuff"}
	// log.Println(classes.Class[0].Ascendencies[0])
	b = GetGems("gems.json")
	//log.Println(string(b))
	buf := new(bytes.Buffer)
	json.Indent(buf, []byte(b), "", "  ")
	file, err := os.Create("gems.json")
	CheckErr(err)
	defer file.Close()
	fmt.Fprintf(file, buf.String())

	var handler func([]byte, []byte, jsonparser.ValueType, int) error
	handler = func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		//log.Printf("Key: '%v'\n", offset) // Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
		rs, _ := jsonparser.GetString(value, "base_item", "release_state")
		dn, _ := jsonparser.GetString(value, "base_item", "display_name")
		sup, _ := jsonparser.GetBoolean(value, "is_support")
		mc, _ := jsonparser.GetBoolean(value, "active_skill", "is_manually_casted")
		as, _ := jsonparser.GetBoolean(value, "active_skill", "is_skill_totem")
		//CheckErr(err)
		if rs == "released" && sup == false && mc == true {
			log.Printf("%v %v", dn, as)
		}

		return nil
	}
	jsonparser.ObjectEach(b, handler)

	//jsonparser.ArrayEach(b, func(thing []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 	// "owner_id": 3,
	// 	oid, err := jsonparser.GetInt(thing, "owner_id")
	// 	CheckErr(err)
	// 	// "last_modification_date": 1511182840,
	// 	modified, err := jsonparser.GetInt(thing, "last_modification_date")
	// 	CheckErr(err)
	// 	// "status": "running",
	// 	status, err := jsonparser.GetString(thing, "status")
	// 	CheckErr(err)
	// 	// "history_id": 10091321,
	// 	hid, err := jsonparser.GetInt(thing, "history_id")
	// 	CheckErr(err)
	// 	// "type": "agent",
	// 	scanType, err := jsonparser.GetString(thing, "type")
	// 	CheckErr(err)
	// 	// "alt_targets_used": false,
	// 	// altTargets, err := jsonparser.Get(thing, "alt_targets_used")
	// 	// CheckErr(err)
	// 	// "scheduler": 0,
	// 	// sched, err := jsonparser.Get(thing, "scheduler")
	// 	// CheckErr(err)
	// 	// "uuid": "e2ef32a8-bce1-4bef-9bbb-223a95201e17",
	// 	uuid, err := jsonparser.GetString(thing, "uuid")
	// 	CheckErr(err)
	// 	// "creation_date": 1511182835
	// 	created, err := jsonparser.GetInt(thing, "creation_date")
	// 	CheckErr(err)

	//}, "history")

}
