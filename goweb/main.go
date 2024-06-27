package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strings"
)

func main() {
	gin.SetMode(gin.DebugMode)
	s := Routers()
	s.Run(":18080")
}

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

var HtmlPath = "../views/blog/"

func Routers() *gin.Engine {
	Router.StaticFS("/static", http.Dir("../static"))

	//模板
	// 自定义模板方法
	Router.SetFuncMap(template.FuncMap{
		"echo":   Echo,
		"h":      H,
		"hc":     HFunc,
		"btn":    ButtonFunc,
		"harder": HarderHtml,
		"input":  FormControlFunc,
		"modal":  ModalFunc,
	})

	Router.Delims("{[", "]}")

	Router.LoadHTMLGlob(HtmlPath + "*")

	Router.GET("/", Case1)
	Router.GET("/case1", Index)
	Router.GET("/case2", Case2)
	Router.GET("/case3", Index3)

	return Router
}

func List(c *gin.Context) {
	c.HTML(http.StatusOK, "list.html", gin.H{})
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func Index2(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case2.html", gin.H{})
}

func Index3(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case3.html", gin.H{})
}

func Case1(c *gin.Context) {

	data1 := &HValue{
		Level: 1,
		Data:  "bbbb",
		Class: "bbb",
	}

	btn1 := &Button{
		Type:     "success",
		Text:     "成功",
		Outline:  true,
		Size:     "lg",
		Disabled: true,
	}

	input1 := &FormControl{
		Id:    "a1",
		Label: "aaa1",
	}

	modal1 := &Modal{
		Title:                   "aaaaa-sadsad",
		ContentTemplateFileName: "t1.html",
	}

	c.HTML(http.StatusOK, "goweb_case1.html", gin.H{
		"data1":  data1,
		"btn1":   btn1,
		"input1": input1,
		"modal1": modal1,
	})
}

func Case2(c *gin.Context) {
	c.HTML(http.StatusOK, "case2.html", gin.H{})
}

func Echo(data string) template.HTML {
	return template.HTML(fmt.Sprintf("<h1>%s</h1>", data))
}

func H(level int, data string, class string, id string) template.HTML {
	if class != "" {
		class = " class=\"" + class + "\""
	}
	if id != "" {
		id = " id=\"" + id + "\""
	}
	return template.HTML(fmt.Sprintf("<h%v %v%v>%v</h1>", level, class, id, data))
}

type HValue struct {
	Level int
	Data  string
	Class string
	Id    string
}

func HFunc(value *HValue) template.HTML {
	if value.Class != "" {
		value.Class = " class=\"" + value.Class + "\""
	}
	if value.Id != "" {
		value.Id = " id=\"" + value.Id + "\""
	}
	return template.HTML(fmt.Sprintf("<h%v %v%v>%v</h1>", value.Level, value.Class, value.Id, value.Data))
}

type Button struct {
	Type     string // 按钮类型
	Text     string // 按钮文本
	Outline  bool   // 是否带轮廓线
	Size     string // 按钮大小 lg sm
	Disabled bool   // 是否禁用
}

var ButtonsType = map[string]string{
	"primary":   `btn btn-primary`,
	"secondary": `btn btn-secondary`,
	"success":   `btn btn-success`,
	"danger":    `btn btn-danger`,
	"warning":   `btn btn-warning`,
	"info":      `btn btn-info`,
	"light":     `btn btn-light`,
	"dark":      `btn btn-dark`,
}

func ButtonFunc(btn *Button, args ...string) template.HTML {
	class := ""
	if value, has := ButtonsType[btn.Type]; has {
		class = value
	}

	id := ""

	for _, arg := range args {
		if strings.Contains(arg, "class=") {
			class += strings.ReplaceAll(arg, "class=", " ")
		}

		if strings.Contains(arg, "id=") {
			id = strings.ReplaceAll(arg, "id=", "")
		}
	}

	if btn.Outline {
		switch btn.Type {
		case "primary":
			class = strings.ReplaceAll(class, "btn-primary", "btn-outline-primary")
		case "secondary":
			class = strings.ReplaceAll(class, "btn-secondary", "btn-outline-secondary")
		case "success":
			class = strings.ReplaceAll(class, "btn-success", "btn-outline-success")
		case "danger":
			class = strings.ReplaceAll(class, "btn-danger", "btn-outline-danger")
		case "warning":
			class = strings.ReplaceAll(class, "btn-warning", "btn-outline-warning")
		case "info":
			class = strings.ReplaceAll(class, "btn-info", "btn-outline-info")
		case "light":
			class = strings.ReplaceAll(class, "btn-light", "btn-outline-light")
		case "dark":
			class = strings.ReplaceAll(class, "btn-dark", "btn-outline-dark")
		}
	}

	switch btn.Size {
	case "lg":
		class += " btn-lg"
	case "sm":
		class += " btn-sm"
	}

	if class != "" {
		class = " class=\"" + class + "\""
	}

	if id != "" {
		id = " id=\"" + id + "\""
	}

	disabled := ""
	if btn.Disabled {
		disabled = " disabled"
	}

	return template.HTML(fmt.Sprintf("<button type=\"button\" %v%v%v>%v</button>", class, id, disabled, btn.Text))
}

func HarderHtml() template.HTML {
	t := `<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="Mark Otto, Jacob Thornton, and Bootstrap contributors">
    <meta name="generator" content="Hugo 0.108.0">
    <title>Blog examples</title>
    <link href="/static/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/font/bootstrap-icons.css" rel="stylesheet">
    <link href="/static/css/blog.css" rel="stylesheet">
	<script src="/static/js/bootstrap.bundle.min.js"></script>
</head>`
	return template.HTML(t)
}

type FormControl struct {
	Id    string
	Label string
}

func FormControlFunc(data *FormControl) template.HTML {
	t := fmt.Sprintf(`<div class="form-floating">
  <input type="email" class="form-control" id="%v" placeholder="name@example.com">
  <label for="floatingInput">%v</label>
</div>`, data.Id, data.Label)
	return template.HTML(t)
}

type Modal struct {
	Title                   string
	ContentTemplateFileName string
}

func ModalFunc(id string, data *Modal) template.HTML {
	aa := template.Must(template.ParseFiles(fmt.Sprintf("%s%s", HtmlPath, data.ContentTemplateFileName)))
	b := []byte("")
	bb := bytes.NewBuffer(b)
	_ = aa.Execute(bb, "")
	t := fmt.Sprintf(`<div class="modal fade" id="%v" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <h1 class="modal-title fs-5" id="exampleModalLabel">%v</h1>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <div class="modal-body">
        %v
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        <button type="button" class="btn btn-primary">Save changes</button>
      </div>
    </div>
  </div>
</div>`, id, data.Title, bb.String())

	return template.HTML(t)
}
