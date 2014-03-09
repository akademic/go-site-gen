package main

import (
    "html/template"
    "path/filepath"
    "strings"
    "os"
    "strconv"
)

func createPost(post *Page) {

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

    save_path := getPostSavePath(post)

    f, err := os.Create( save_path )
    if err != nil {
        panic(err)
    }
    defer f.Close()

    t.Execute(f, postTemplData)
}

func getPostSavePath(post *Page) (string) {
    parts := strings.Split(post.Path, "/")
    name := parts[len(parts)-2 : len(parts)-1][0]
    year := strconv.Itoa(post.PubTime.Year())

    save_dir := filepath.Join(BlogDir, year, name)
    os.MkdirAll(save_dir, 0700)
    save_path := filepath.Join(save_dir, "index.html")

    return save_path
}
