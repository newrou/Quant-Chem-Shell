package main

import (
    "fmt"
    "html/template"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "time"
    "crypto/md5"
    "log"
)

var works_dir = "/home/alex/.md-server/works"

type WorkDetails struct {
    Title        string
    Temp         string
    Pressure     string
    Stat         string
    Compounds    string
    Status       string
}

type Todo struct {
    Title   string
    Status  string
    Done    bool
}

type TodoPageData struct {
    PageTitle string
    Todos     []Todo
}

func GetWorkList(works []Todo) []Todo {
    files, err := ioutil.ReadDir(works_dir)
    if err != nil { log.Fatal(err) }
    for _, file := range files {
//        fmt.Println(file.Name(), file.IsDir())
	works = append(works, Todo{Title: file.Name(), Status: "done", Done: false})
    }
    return works
}

func main() {
    tmpl_list := template.Must(template.ParseFiles("form_list.html"))
    tmpl_add := template.Must(template.ParseFiles("form_add.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := TodoPageData{
            PageTitle: "List of works:",
            Todos: []Todo{},
        }
	data.Todos = GetWorkList(data.Todos)
        tmpl_list.Execute(w, data)
    })

    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl_add.Execute(w, nil)
            return
        }

        work := WorkDetails{
            Title:     r.FormValue("Title"),
            Temp:      r.FormValue("Temp"),
            Pressure:  r.FormValue("Pressure"),
            Stat:      r.FormValue("Stat"),
            Compounds: r.FormValue("Compounds"),
            Status:    r.FormValue("Status"),
        }

        dat, err := json.MarshalIndent(work, "", " ")
        if err != nil { fmt.Println(err) }
        utime := int32(time.Now().Unix())
	hmd5 := md5.Sum([]byte(work.Title))
	fname := fmt.Sprintf("%s/%d-%x", works_dir, utime, hmd5)
        _ = ioutil.WriteFile(fname, dat, 0644)
        fmt.Println(fname, string(dat))
        _ = work
        tmpl_add.Execute(w, struct{ Success bool }{true})
    })

    http.ListenAndServe(":8080", nil)
}
