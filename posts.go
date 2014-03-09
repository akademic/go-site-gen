package main

/*
Posts are in _posts directory
Post -> Prepare custom tags -> MarkdownToHtml -> html/template into ready file
*/

import (
    "os"
    "io"
    "io/ioutil"
    "log"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

var (
    filesToMove []string
)

type sortablePosts []*Page
func (s sortablePosts) Len() int           { return len(s) }
func (s sortablePosts) Less(i, j int) bool { return s[i].PubTime.Before(s[j].PubTime) }
func (s sortablePosts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func genPosts() {

    posts := getPosts()
    for _, post := range posts {
        save_dir := getPostSavePath(post)
        save_path := filepath.Join(save_dir, "index.html")

        createPost(post, save_path)
        movePostFiles(filepath.Dir(post.Path), save_dir)
    }
}

func movePostFiles(fromDir, toDir string) {
    filesToMove = []string{}
    filepath.Walk(fromDir, postDirVisit)

    for _, file := range filesToMove {
        dir, name := filepath.Split(file)
        dir, _ = filepath.Rel(fromDir, dir)

        toPath := filepath.Join(toDir, dir)
        os.MkdirAll(toPath, 0700)

        toFile := filepath.Join(toPath, name)

        from, _ := os.Open(file)
        to, _ := os.Create(toFile)

        io.Copy(to, from)

    }
}

//FIXME: Not tread safe
func postDirVisit(path string, f os.FileInfo, err error) error {
    if !f.IsDir() {
        if f.Name() != "index.md" {
            filesToMove = append(filesToMove, path)
        }
    }
    return err
}

func getPostSavePath(post *Page) (string) {
    parts := strings.Split(post.Path, "/")
    name := parts[len(parts)-2 : len(parts)-1][0]
    year := strconv.Itoa(post.PubTime.Year())

    save_dir := filepath.Join(BlogDir, year, name)
    os.MkdirAll(save_dir, 0700)

    return save_dir
}

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

