package main

import (
    "html/template"
    "path/filepath"
    "os"
)

func createPost(post *Post) {

    postTemplData := map[string] interface {} {
        "SiteData": SiteDataVar,
        "Title": post.Title,
        "Body": template.HTML(post.Content),
        "PubTime": post.PubTime,
        "Next": post.Next,
        "Prev": post.Prev,
    }

    t, err := template.ParseFiles( filepath.Join( TemplatesDir, "default.html"),
                                    filepath.Join( TemplatesDir, "post.html"))
    if err != nil {
        panic(err)
    }

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

func html (x string) interface{} { return template.HTML(x) }

func createIndexPage() {
    tpl := new(templIndex)
    tpl.SiteData = SiteDataVar
    tpl.Recent = getPostRecent()

    t, err := template.ParseFiles( filepath.Join( TemplatesDir, "default.html"),
                                    filepath.Join( TemplatesDir, "index.html"))
    if err != nil {
        panic(err)
    }

    t = t.Funcs(template.FuncMap{"html": html})
    
    f, err := os.Create( filepath.Join(PublicDir, "index.html") )

    if err != nil {
        panic(err)
    }
    defer f.Close()

    t.Execute(f, tpl)
}
