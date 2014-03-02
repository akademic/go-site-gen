package main

/*
Posts are in _posts directory
Post -> Prepare custom tags -> MarkdownToHtml -> html/template into ready file
*/

import (
    "os"
    "io/ioutil"
    "log"
    "path/filepath"
    "sort"
)

type sortablePosts []*Page
func (s sortablePosts) Len() int           { return len(s) }
func (s sortablePosts) Less(i, j int) bool { return s[i].PubTime.Before(s[j].PubTime) }
func (s sortablePosts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func getPosts() ([]*Page) {

    list, err := getPostsList()
    if err != nil {
        log.Fatal("FATAL ", err)
    }

    posts := make([]*Page, 0, len(list))

    for _, fi := range list {
        post := loadPost(fi)
        posts = append(posts, post)
    }

    sort.Sort(sort.Reverse(sortablePosts(posts)))

    return posts
}

func getPostsList() ([]os.FileInfo, error) {
    fis, err := ioutil.ReadDir(PostsDir)
	if err != nil {
		return nil, err
	}

    //filter only directories
    for i := 0; i < len(fis); {
        if !fis[i].IsDir() {
            fis[i], fis = fis[len(fis)-1], fis[:len(fis)-1]
        } else {
            i++
        }
    }

    return fis, nil
}

func loadPost( fi os.FileInfo ) ( *Page ) {
    post_path := filepath.Join(PostsDir, fi.Name(), "index.md")
    
    post := loadPage(post_path)

    return post
}

