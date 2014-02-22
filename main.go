package main

import (
    "os"
    "path/filepath"
    "log"
)

var (
    PostsDir string
    PublicDir string
    TemplatesDir string
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
