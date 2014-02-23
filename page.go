package main

import (
    "time"
    "regexp"
    "io/ioutil"
    "errors"
    "fmt"
)

var (
    HeaderRE = regexp.MustCompile("^(?s)(@.+?)\n\n")
    AttrsRE = regexp.MustCompile("(?Um)^@([^\\:]+?)\\: (.+)$")
    TitleRE = regexp.MustCompile("#([^#\n]+)")
)

type Page struct {
    Title string
    Author string
    PubTime time.Time
    Layout string
    Content string
}

func loadPage(path string) ( *Page ) {
    fmt.Println(path)
    bytes, err := ioutil.ReadFile(path)

    if err != nil {
        die("Unable to read file [%s]", path)
    }

    page, err := newPage(bytes)
    fmt.Println(page)
    if err != nil {
        die("Unable to create new page [%s]: "+err.Error(), path)
    }

    return page
}

func newPage(data []byte) ( *Page, error ) {
    content := string(data)

    page, err := parsePage(content)
    if err != nil {
        return nil, err
    }

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
        page.Layout = h["layout"]
        page.Author = h["author"]

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
        
        page.Content = content

        if m := TitleRE.FindStringSubmatch(content); m != nil {
            page.Title = m[1]
            fmt.Println("Title found:", page.Title)
        } else {
            return nil, errors.New("Title not found")
        }
    }

    return page, nil
}
