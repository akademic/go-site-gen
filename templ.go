package main

import (
    "html/template"
    "path/filepath"
    "os"
)

func createPost(post *Page, save_path string) {

    postTemplData := map[string] interface {} {
        "Title": post.Title,
        "Body": template.HTML(post.Content),
        "PubTime": post.PubTime,
    }

    t, err := template.ParseFiles( filepath.Join( TemplatesDir, "default.html"),
                                    filepath.Join( TemplatesDir, "post.html"))
    if err != nil {
        panic(err)
    }

    f, err := os.Create( save_path )
    if err != nil {
        panic(err)
    }
    defer f.Close()

    t.Execute(f, postTemplData)
}
