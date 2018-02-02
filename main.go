package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/buger/jsonparser"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// CheckErr to handle errors
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetNums gets the random numbers to determine class and ascendency
func GetNums(num1 int) (classNum int, ascenNum int) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	choice := r1.Intn(num1)
	s2 := rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s2)
	ascen := r2.Intn(3)
	return choice, ascen
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
// Classes should contain a list of all the classes and ascendencies
type Classes []struct {
	Class     string   `json:"Class"`
	Ascension []string `json:"Ascension"`
}

// Result contains the result of choosing a random class / ascendency
type Result struct {
	Class      string `json:"class" xml:"class"`
	Ascendency string `json:"ascendency" xml:"ascendency"`
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	e := echo.New()
	e.Static("/static", "static")
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	var b []byte
	var Classes Classes

	c := []byte(`[{"Class":"Duelist","Ascension":["Slayer","Gladiator","Champion"]},{"Class":"Shadow","Ascension":["Assassin","Saboteur","Trickster"]},{"Class":"Marauder","Ascension":["Juggernaut","Berserker","Chieftain"]},{"Class":"Witch","Ascension":["Necromancer","Occultist","Elementalist"]},{"Class":"Ranger","Ascension":["Deadeye","Raider","Pathfinder"]},{"Class":"Templar","Ascension":["Inquisitor","Hierophant","Guardian"]},{"Class":"Scion","Ascension":["Ascendant","Ascendant","Ascendant"]}]`)
	err := json.Unmarshal(c, &Classes)

	//log.Printf("Class: %v Ascendency: %v", classes[choice].Class, classes[choice].Ascension[ascen])

	log.Println("Starting Stuff")
	// classes.Class = "Duelist"
	// classes.Ascendencies = ["1", "2", ]
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
		//as, _ := jsonparser.GetBoolean(value, "active_skill", "is_skill_totem")
		tags, _ := jsonparser.GetString(value, "tags")
		//CheckErr(err)
		if rs == "released" && sup == false && mc == true {
			log.Printf("%v %v", dn, tags)
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

	// rndhtml := "stuff"
	// e.GET("/rnd", func(c echo.Context) error {
	// 	num1 := 7
	// 	scion := c.FormValue("scion")
	// 	if scion == "true" {
	// 		num1 = 6
	// 	}
	// 	choice, ascen := GetNums(num1)
	// 	class := Classes[choice].Class
	// 	ascendency := Classes[choice].Ascension[ascen]

	// 	return c.HTML(http.StatusOK, rndhtml)
	// })

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer

	// Named route "foobar"
	e.GET("/rnd", func(c echo.Context) error {
		num1 := 6
		scion := c.FormValue("scion")
		if scion == "true" {
			num1 = 7
		}
		temp := "rnd.html"
		league := c.FormValue("league")
		if league == "true" {
			temp = "rndl.html"
		}
		choice, ascen := GetNums(num1)
		return c.Render(http.StatusOK, temp, map[string]interface{}{
			"class":      Classes[choice].Class,
			"ascendency": Classes[choice].Ascension[ascen],
			"name":       randomdata.SillyName(),
			"gem":        "",
			"league":     randomdata.StringSample("Abyss", "Standard"),
			"ssf":        randomdata.StringSample("SSF", "Normal"),
			"hc":         randomdata.StringSample("HardCore", "SoftCore"),
		})
	}).Name = "rnd"

	e.GET("/api", func(c echo.Context) error {
		num1 := 6
		scion := c.FormValue("scion")
		if scion == "true" {
			num1 = 7
		}
		choice, ascen := GetNums(num1)
		r := map[string]interface{}{
			"class":      Classes[choice].Class,
			"ascendency": Classes[choice].Ascension[ascen],
			"name":       randomdata.SillyName(),
			"gem":        "",
			"league":     randomdata.StringSample("Abyss", "Standard"),
			"ssf":        randomdata.StringSample("SSF", "Normal"),
			"hc":         randomdata.StringSample("HardCore", "SoftCore"),
		}
		return c.JSON(http.StatusOK, r)
	})

	e.GET("/", func(c echo.Context) error {

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
	}).Name = "index"

	e.Logger.Fatal(e.Start(":1323"))

}
