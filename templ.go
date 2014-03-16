package main

import (
    "html/template"
    "path/filepath"
    "os"
    "time"
)

var FuncsMap = template.FuncMap{
    "fmttime": func(t time.Time, f string) string {
        return t.Format(f)
    },
    "html": func(x string) interface{} {return template.HTML(x)},
}

func createPost(post *Post) {

    postTemplData := map[string] interface {} {
        "SiteData": SiteDataVar,
        "Title": post.Title,
        "Body": template.HTML(post.Content),
        "PubTime": post.PubTime,
        "Next": post.Next,
        "Prev": post.Prev,
    }

    t := template.Must(template.New("default.html").Funcs(FuncsMap).ParseFiles(filepath.Join( TemplatesDir, "default.html"), filepath.Join( TemplatesDir, "post.html")) )

    f, err := os.Create( filepath.Join(post.SavePath, "index.html") )
    if err != nil {
        panic(err)
    }
    defer f.Close()

    t.Execute(f, postTemplData)
}

type templIndex struct {
    SiteData SiteData
    Recent []*Post
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
