package main

/*
Posts are in _posts directory
Post -> Prepare custom tags -> MarkdownToHtml -> html/template into ready file
*/

import (
    "os"
    "time"
    "io/ioutil"
    "log"
    "path/filepath"
)

type Post struct {
    Slug string
    Title string
    PubTime time.Time
    Content []byte
}

func getPosts() ([]*Post) {

    list, err := getPostsList()
    if err != nil {
        log.Fatal("FATAL ", err)
    }

    posts := make([]*Post, 0, len(list))

    for _, fi := range list {
        post, err := loadPost(fi)
    }
}

func getPostsList() ([]os.FileInfo, error) {
    fis, err := ioutil.ReadDir(PostsDir)
	if err != nil {
		return nil, err
	}

    for i := 0; i < len(fis); {
        if !fis[i].IsDir() {
            fis[i], fis = fis[len(fis)-1], fis[:len(fis)-1]
        } else {
            i++
        }
    }

    return fis, nil
}

func loadPost( fi os.FileInfo ) ( *Post, error ) {
    f, err := os.Open(filepath.Join(PostsDir, fi.Name(), "index.md"))
}
