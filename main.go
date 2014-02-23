package main

import (
    "os"
    "path/filepath"
    "log"
    "fmt"
)

var (
    PostsDir string
    PublicDir string
    TemplatesDir string

)

const (
    DateFormat     = "2006-01-02"
    DateTimeFormat = DateFormat + " 15:04"
)

func _init() {

    // Initialize directories
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL ", err)
	}

    PublicDir = filepath.Join(pwd, "_public")
	PostsDir = filepath.Join(pwd, "_posts")
	TemplatesDir = filepath.Join(pwd, "_layouts")

}

func main() {
    _init()
    loadPage(PostsDir+"/test_post/index.md")
}

func die(format string, v ...interface{}) {
    os.Stderr.WriteString(fmt.Sprintf(format+"\n", v...))
    os.Exit(1)
}
