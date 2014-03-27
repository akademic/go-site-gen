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
    Posts []*Post
)

type Post struct {
    *Page
    SavePath string
    RelSavePath string
    Next *Post
    Prev *Post
}

type sortablePosts []*Post
func (s sortablePosts) Len() int           { return len(s) }
func (s sortablePosts) Less(i, j int) bool { return s[i].PubTime.Before(s[j].PubTime) }
func (s sortablePosts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func genPosts() {

    Posts = getPosts()
    for i, post := range Posts {

        if i > 0 {
            if _, ok := Posts[i-1].Headers["draft"]; !ok {
                post.Next = Posts[i-1]
            }
        }

        if i < len(Posts) - 1 {
            if _, ok := Posts[i+1].Headers["draft"]; !ok {
                post.Prev = Posts[i+1]
            }
        }

        createPost(post)
        movePostFiles(filepath.Dir(post.Path), post.SavePath)
    }
}

func getPostRecent() []*Post {

    var recent_count int
    if len(Posts) < int(RecentPostsCount) {
        recent_count = len(Posts)
    } else {
        recent_count = int(RecentPostsCount)
    }

    recent := make( []*Post, 0 )
    for i := 0; i < recent_count; i++ {
        if _, ok := Posts[i].Headers["draft"]; ok {
            if recent_count < len(Posts) {
                recent_count++
            }
        } else {
            recent = append(recent, Posts[i])
        }
    }
    
    return recent
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

func getPostSavePath(post *Page) (string, string) {
    parts := strings.Split(post.Path, "/")
    name := parts[len(parts)-2]
    year := strconv.Itoa(post.PubTime.Year())

    save_dir := filepath.Join(BlogDir, year, name)
    os.MkdirAll(save_dir, 0700)

    rel_save_dir, _ := filepath.Rel(PublicDir, save_dir)

    return rel_save_dir, save_dir
}

func getPosts() ([]*Post) {

    list, err := getPostsList()
    if err != nil {
        log.Fatal("FATAL ", err)
    }

    posts := make([]*Post, 0, len(list))

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

func loadPost( fi os.FileInfo ) ( *Post ) {
    post_path := filepath.Join(PostsDir, fi.Name(), "index.md")
    
    page := loadPage(post_path)
    
    rel_save_dir, save_dir := getPostSavePath(page)

    post := &Post{
        page,
        save_dir,
        rel_save_dir,
        nil,
        nil,
    }

    return post
}

