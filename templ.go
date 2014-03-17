package main

import (
    "html/template"
    "path/filepath"
    "os"
    "time"
    "io/ioutil"
    "strings"
)

type templIndex struct {
    SiteData SiteData
    Recent []*Post
}

type templPost struct {
    SiteData SiteData
    Post *Post
}

type templRss struct {
    XmlHeader string
    SiteData SiteData
    PubTime time.Time
    Recent []*Post
}

type templStatic struct {
    SiteData SiteData
    Page *Page
}

var FuncsMap = template.FuncMap{
    "fmttime": func(t time.Time, f string) string {
        return t.Format(f)
    },
    "html": func(x string) interface{} {return template.HTML(x)},
    "include": func(file string) string {
        data, _ := ioutil.ReadFile(filepath.Join(IncludesDir, file))
        return string(data)
    },
}

func createPost(post *Post) {

    tpl := new(templPost)
    tpl.SiteData = SiteDataVar
    tpl.Post = post

    t := template.Must(template.New("default.html").Funcs(FuncsMap).ParseFiles(filepath.Join( TemplatesDir, "default.html"), filepath.Join( TemplatesDir, "post.html")) )

    f, err := os.Create( filepath.Join(post.SavePath, "index.html") )
    if err != nil {
        panic(err)
    }
    defer f.Close()

    t.Execute(f, tpl)
}

func createIndexPage() {
    tpl := new(templIndex)
    tpl.SiteData = SiteDataVar
    tpl.Recent = getPostRecent()

    t := template.Must(template.New("default.html").Funcs(FuncsMap).ParseFiles(filepath.Join( TemplatesDir, "default.html"), filepath.Join( TemplatesDir, "index.html")) )
    
    f, err := os.Create( filepath.Join(PublicDir, "index.html") )

    if err != nil {
        panic(err)
    }
    defer f.Close()

    t.Execute(f, tpl)
}

func createPostRss() {
    tpl := new(templRss)
    tpl.XmlHeader = `<?xml version="1.0" encoding="UTF-8"?>`
    tpl.SiteData = SiteDataVar
    tpl.Recent = getPostRecent()
    tpl.PubTime = tpl.Recent[0].PubTime

    t := template.Must(template.New("rss.xml").Funcs(FuncsMap).ParseFiles(filepath.Join( TemplatesDir, "rss.xml")) )
    
    f, err := os.Create( filepath.Join(PublicDir, "rss.xml") )
    
    if err != nil {
        panic(err)
    }
    defer f.Close()

    err = t.Execute(f, tpl)
    if err != nil {
        panic(err)
    }
}

func createStaticPage(page *Page, path string, name string) {
    tpl := new(templStatic)
    tpl.SiteData = SiteDataVar
    tpl.Page = page

    fname := strings.Split(name, ".")[0]

    save_path := filepath.Join(path, fname+".html")

    main_layout := "default.html"
    if page.Layout != "" {
        main_layout = page.Layout
    }

    t := template.Must(template.New(main_layout).Funcs(FuncsMap).ParseFiles(filepath.Join( TemplatesDir, main_layout), filepath.Join( TemplatesDir, "static.html")) )

    f, err := os.Create(save_path)

    if err != nil {
        panic(err)
    }
    defer f.Close()

    err = t.Execute(f, tpl)
    if err != nil {
        panic(err)
    }
}
