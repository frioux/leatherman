package now

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/frioux/leatherman/internal/notes"
)

//go:embed templates/*
var templates embed.FS

var tpl *template.Template

func init() {
	var err error
	tpl, err = template.ParseFS(templates, "templates/*")
	if err != nil {
		panic(err)
	}
}

type HTMLVars struct {
	*notes.Zine

	Title string
	body  []byte
}

func (v *HTMLVars) Body() template.HTML {
	return template.HTML(string(v.body))
}

func (v *HTMLVars) TODOCount() (int, error) {
	const sql = `SELECT COUNT(*) FROM _ WHERE tag = 'inbox' AND title != '000 IN'`
	var c int
	if err := v.Zine.DB.Get(&c, sql); err != nil {
		return 0, err
	}
	return c, nil
}

func (v *HTMLVars) Write(b []byte) (int, error) {
	v.body = append(v.body, b...)
	return len(b), nil
}

type articleVars struct {
	*HTMLVars
	ArticleTitle, Filename string
}

func (v articleVars) Title() string { return "now: " + v.ArticleTitle }

type listVars struct {
	*HTMLVars
	Articles []struct {
		Title, URL string
	}
}

func (v listVars) Title() string { return "now: list" }

type qVars struct {
	*HTMLVars
	Records []map[string]interface{}
}

func (v qVars) Title() string { return "now: q" }

type updateVars struct {
	*HTMLVars
	File, Content string
}

func (v updateVars) Title() string { return "now: update " + v.File }

type option struct {
	error
	value interface{}
}

func (o option) HTML() template.HTML {
	if err := o.error; err != nil {
		return template.HTML(`<span style="color: red">` + err.Error() + `</span>`)
	}

	return template.HTML(fmt.Sprint(o.value))
}

type supVars struct {
	*HTMLVars
	Versions                 []option
	retroPie, steamOS, pi400 option
}

func (v supVars) RetroPie() template.HTML { return v.retroPie.HTML() }

func (v supVars) SteamOS() template.HTML { return v.steamOS.HTML() }

func (v supVars) Pi400() template.HTML { return v.pi400.HTML() }

func (v supVars) Title() string { return "now: sup" }
