package main

import (
    "html/template"
    "path/filepath"
    "os"
    "time"
    "io/ioutil"
)

type templIndex struct {
    SiteData SiteData
    Recent []*Post
}

type templPost struct {
    SiteData SiteData
    Post *Post
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
