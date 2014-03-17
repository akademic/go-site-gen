package main

import (
    "os"
    "time"
    "regexp"
    "io/ioutil"
    "errors"
    "fmt"
    "github.com/russross/blackfriday"
)

var (
    HeaderRE = regexp.MustCompile("^(?s)(@.+?)\n\n")
    AttrsRE = regexp.MustCompile("(?Um)^@([^\\:]+?)\\: (.+)$")
    TitleRE = regexp.MustCompile("^# ([^#\n]+)")
    FirstParaRE = regexp.MustCompile(`^# [^#\n]+\n\n(([^\n]+\n)+)\n`)
)

type Page struct {
    Path string
    Title string
    Headers map[string]string
    PubTime time.Time
    UpdateTime time.Time
    Layout string
    Summary string
    Content string
}

func loadPage(path string) ( *Page ) {
    fmt.Println(path)

    page_fi, err := os.Stat(path)
    if err != nil {
        die("Unable to stat page file [%s]", path)
    }

    bytes, err := ioutil.ReadFile(path)

    if err != nil {
        die("Unable to read file [%s]", path)
    }

    page, err := newPage(bytes)
    if err != nil {
        die("Unable to create new page [%s]: "+err.Error(), path)
    }
    page.Path = path
    page.UpdateTime = page_fi.ModTime()
    fmt.Println(page)

    return page
}

func newPage(data []byte) ( *Page, error ) {
    content := string(data)

    page, err := parsePage(content)
    if err != nil {
        return nil, err
    }

    page.Content = string(blackfriday.MarkdownCommon([]byte(page.Content)))
    page.Summary = string(blackfriday.MarkdownCommon([]byte(page.Summary)))

    return page, nil
}

func parsePage(content string) (*Page, error) {
    page := &Page{}
    h := make(map[string]string)

    if len(content) > 0 && content[0] == '@' {
        //parsing headers
        fmt.Println("Start parsing headers")
        if m := HeaderRE.FindStringSubmatch(content); m != nil {
            header := m[0]
            if m := AttrsRE.FindAllStringSubmatch(header, -1); m != nil {
                for _, pair := range m {
                    name, value := pair[1], pair[2]
                    h[name] = value
                }
            } else {
                return nil, errors.New("Bad headers format")
            }
        } else {
            return nil, errors.New("Headers not found")
        }

        fmt.Println("Headers:", h)
        page.Headers = h
        page.Layout = h["layout"]

        //date header parsing
        var err error = nil
        var d time.Time

        if len(h["date"]) == 10 {
            d, err = time.Parse(DateFormat, h["date"])
        } else if len(h["date"]) == 16 {
            d, err = time.Parse(DateTimeFormat, h["date"])
        }

        if err != nil {
            return nil, fmt.Errorf("Unable to parse date [%s]", h["date"])
        }

        page.PubTime = d
    
        //get clean content
        content = HeaderRE.ReplaceAllLiteralString(content, "")

        if m := FirstParaRE.FindStringSubmatch(content); len(m) != 0 {
            page.Summary = m[1]
        } else {
            return nil, errors.New("Summary not found")
        }
        

        if m := TitleRE.FindStringSubmatch(content); m != nil {
            page.Title = m[1]
            fmt.Println("Title found:", page.Title)
        } else {
            return nil, errors.New("Title not found")
        }

        content = TitleRE.ReplaceAllLiteralString(content, "")
        
        page.Content = content
    }

    return page, nil
}
