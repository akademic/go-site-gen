package main

import (
    "os"
    "io"
    "path/filepath"
)

func genSite() {
    filepath.Walk(SiteDir, siteDirVisit) 
}

func siteDirVisit( path string, f os.FileInfo, err error) error {
    if !f.IsDir() {
        dir, name  := filepath.Split(path)
        dir, _ = filepath.Rel(SiteDir, dir)

        toPath := filepath.Join(PublicDir, dir)
        os.MkdirAll(toPath, 0700)

        toFile := filepath.Join(toPath, name)
        from, _ := os.Open(path)
        to, _ := os.Create(toFile)

        io.Copy(to, from)
    }

    return err
}
