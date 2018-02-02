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
	"strings"
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
func GetNums(num1 int, gems int) (classNum int, ascenNum int, gemNum int) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	choice := r1.Intn(num1)
	s2 := rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s2)
	ascen := r2.Intn(3)
	s3 := rand.NewSource(time.Now().UnixNano())
	r3 := rand.New(s3)
	gem := r3.Intn(gems)
	return choice, ascen, gem
}

// GetGems makes a GET request to the API
func GetGems(uri string) (b []byte) {
	c := &http.Client{}
	r, err := http.NewRequest("GET", "https://raw.githubusercontent.com/brather1ng/RePoE/master/data/"+uri, nil)
	CheckErr(err)
	resp, err := c.Do(r)
	CheckErr(err)
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	CheckErr(err)
	return b
}

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
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
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
	var tags []string
	//var gems []Gems

	//log.Println(gems)

	c := []byte(`[{"Class":"Duelist","Ascension":["Slayer","Gladiator","Champion"]},{"Class":"Shadow","Ascension":["Assassin","Saboteur","Trickster"]},{"Class":"Marauder","Ascension":["Juggernaut","Berserker","Chieftain"]},{"Class":"Witch","Ascension":["Necromancer","Occultist","Elementalist"]},{"Class":"Ranger","Ascension":["Deadeye","Raider","Pathfinder"]},{"Class":"Templar","Ascension":["Inquisitor","Hierophant","Guardian"]},{"Class":"Scion","Ascension":["Ascendant","Ascendant","Ascendant"]}]`)
	err := json.Unmarshal(c, &Classes)
	b = GetGems("gems.json")
	buf := new(bytes.Buffer)
	json.Indent(buf, []byte(b), "", "  ")
	file, err := os.Create("gems.json")
	CheckErr(err)
	defer file.Close()
	fmt.Fprintf(file, buf.String())
	gems := []string{}
	var handler func([]byte, []byte, jsonparser.ValueType, int) error
	handler = func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		skip := false
		rs, _ := jsonparser.GetString(value, "base_item", "release_state")
		dn, _ := jsonparser.GetString(value, "base_item", "display_name")
		sup, _ := jsonparser.GetBoolean(value, "is_support")
		mc, _ := jsonparser.GetBoolean(value, "active_skill", "is_manually_casted")
		tag, _, _, _ := jsonparser.Get(value, "active_skill", "types")
		_ = json.Unmarshal(tag, &tags)
		for _, t := range tags {
			//log.Println(t)
			switch t {
			case "vaal", "aura", "curse", "movement", "buff":
				skip = true
			}
		}
		if rs == "released" && sup == false && mc == true && skip == false {
			//log.Printf("%v %v", dn, tags[0])
			gems = append(gems, dn)
		}

		return nil
	}
	jsonparser.ObjectEach(b, handler)

	numGems := len(gems)

	log.Println(numGems)

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer

	e.GET("/api", func(c echo.Context) error {
		num1 := 7
		scion := c.FormValue("scion")
		if scion == "no" {
			num1 = 6
		}
		choice, ascen, rndgem := GetNums(num1, numGems)
		gemImg := strings.Replace(gems[rndgem], " ", "_", -1)

		r := map[string]interface{}{
			"class":      Classes[choice].Class,
			"ascendency": Classes[choice].Ascension[ascen],
			"name":       randomdata.SillyName(),
			"gemImg":     gemImg,
			"gem":        gems[rndgem],
			"league":     randomdata.StringSample("SSF Hardcore Abyss", "Normal Hardcore Abyss", "SSF Softcore Abyss", "Normal Softcore Abyss", "SSF Hardcore Standard", "Normal Hardcore Standard", "SSF Softcore Standard", "Normal Softcore Standard"),
			// "ssf":        randomdata.StringSample("SSF", "Normal"),
			// "hc":         randomdata.StringSample("HardCore", "SoftCore"),
		}
		return c.JSON(http.StatusOK, r)
	})

	e.GET("/", func(c echo.Context) error {
		num1 := 7
		scion := c.FormValue("scion")
		if scion == "no" {
			num1 = 6
		}
		temp := "rnd.html"
		league := c.FormValue("league")
		if league == "true" {
			temp = "rndl.html"
		}
		choice, ascen, rndgem := GetNums(num1, numGems)
		gemImg := strings.Replace(gems[rndgem], " ", "_", -1)

		return c.Render(http.StatusOK, temp, map[string]interface{}{
			"class":      Classes[choice].Class,
			"ascendency": Classes[choice].Ascension[ascen],
			"name":       randomdata.SillyName(),
			"gemImg":     gemImg,
			"gem":        gems[rndgem],
			"league":     randomdata.StringSample("SSF Hardcore Abyss", "Normal Hardcore Abyss", "SSF Softcore Abyss", "Normal Softcore Abyss", "SSF Hardcore Standard", "Normal Hardcore Standard", "SSF Softcore Standard", "Normal Softcore Standard"),
			// "ssf":        randomdata.StringSample("SSF", "Normal"),
			// "hc":         randomdata.StringSample("HardCore", "SoftCore"),
		})
	}).Name = "index"

	e.Logger.Fatal(e.Start(":2086"))

}
