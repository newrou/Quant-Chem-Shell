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
//    "reflect"
)

var works_dir = "/home/alex/.md-server/works"

type Work struct {
    Id           string
    Title        string
    Temp         string
    Pressure     string
    Stat         string
    Compounds    string
    Status       string
    Done    bool
}

type PageData struct {
    PageTitle    string
    WorkList     []Work
}

func LoadWork(fname, id string) Work {
    var w Work
    content, err := ioutil.ReadFile(fname)
    if err != nil { log.Fatal("Error when opening file: ", err) }
    var data map[string]interface{}
    err = json.Unmarshal(content, &data)
    if err != nil { log.Fatal("Error during Unmarshal(): ", err) }
//    fmt.Println(data)
//    fmt.Println(reflect.TypeOf(data["Title"]), data["Title"])
    w.Id = id
    w.Title = fmt.Sprintf("%s", data["Title"])
    w.Temp = fmt.Sprintf("%s", data["Temp"])
    w.Pressure = fmt.Sprintf("%s", data["Pressure"])
    w.Stat = fmt.Sprintf("%s", data["Stat"])
    w.Compounds = fmt.Sprintf("%s", data["Compounds"])
    w.Status = fmt.Sprintf("%s", data["Status"])
    w.Done = false
    return w
}

func GetWorkList(works []Work) []Work {
    files, err := ioutil.ReadDir(works_dir)
    if err != nil { log.Fatal(err) }
    for _, file := range files {
//        fmt.Println(file.Name(), file.IsDir())
	fname := fmt.Sprintf("%s/%s", works_dir, file.Name())
//        fmt.Println(fname)
	w := LoadWork(fname, file.Name())
//        fmt.Println(w, w.Status)
	works = append(works, w)
//	works = append(works, Work{Title: fname, Status: "done", Done: false})
    }
    return works
}

func main() {
    tmpl_add := template.Must(template.ParseFiles("form_add.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "List of works:",
            WorkList: []Work{},
        }
	tmpl_list := template.Must(template.ParseFiles("form_list.html"))
	data.WorkList = GetWorkList(data.WorkList)
        tmpl_list.Execute(w, data)
    })

    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl_add.Execute(w, nil)
            return
        }

        work := Work{
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
