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
    "os"
//    "reflect"
)

var conf_dir = ".cs-server/"
var works_dir = conf_dir + "works"
var archiv_dir = conf_dir + "archiv"

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
    Message      string
    WorkList     []Work
}

func LoadWork(id string) Work {
    var w Work
    fname := fmt.Sprintf("%s/%s", works_dir, id)
    content, err := ioutil.ReadFile(fname)
    if err != nil { /*log.Fatal("Error when opening file: ", err)*/ return w }
    var data map[string]interface{}
    err = json.Unmarshal(content, &data)
    if err != nil { /*log.Fatal("Error during Unmarshal(): ", err)*/ return w }
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
    if err != nil { /*log.Fatal(err)*/ return works }
    for _, file := range files {
	w := LoadWork(file.Name())
	works = append(works, w)
    }
    return works
}

func main() {
    file, err := os.OpenFile(conf_dir + "cs-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    // if err != nil { log.Fatal(err) }
    if err == nil { log.SetOutput(file) }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "List of works:",
            WorkList: []Work{},
        }
	tmpl_list := template.Must(template.ParseFiles("form_list.html"))
	data.WorkList = GetWorkList(data.WorkList)
        tmpl_list.Execute(w, data)
	log.Println("View list of works")
    })


    http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
	tmpl_view := template.Must(template.ParseFiles("form_view.html"))
        if r.Method != http.MethodGet {
            tmpl_view.Execute(w, nil)
            return
        }
	Id := r.FormValue("id")
	data := LoadWork(Id)
        tmpl_view.Execute(w, data)
	log.Println("View work: ", Id)
    })


    http.HandleFunc("/set-status", func(w http.ResponseWriter, r *http.Request) {
	tmpl_view := template.Must(template.ParseFiles("form_view.html"))
        if r.Method != http.MethodGet {
            tmpl_view.Execute(w, nil)
            return
        }
	Id := r.FormValue("id")
	Status := r.FormValue("status")
	work := LoadWork(Id)
	work.Status = Status
        dat, err := json.MarshalIndent(work, "", " ")
        if err != nil { fmt.Println(err) }
	fname := fmt.Sprintf("%s/%s", works_dir, Id)
        _ = ioutil.WriteFile(fname, dat, 0644)
        tmpl_view.Execute(w, work)
	log.Println("Set status work", Id, "to", Status)
    })


    http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
	tmpl_auto := template.Must(template.ParseFiles("form_auto.html"))
        if r.Method != http.MethodGet {
            tmpl_auto.Execute(w, nil)
            return
        }
	Id := r.FormValue("id")
	p1 := fmt.Sprintf("%s/%s", works_dir, Id)
	p2 := fmt.Sprintf("%s/%s", archiv_dir, Id)
	e := os.Rename(p1, p2)
	if e != nil { /* log.Fatal(e) */ }
	msg := fmt.Sprintf("Remove work: %s", Id)
        data := PageData { Message: msg, }
        tmpl_auto.Execute(w, data)
	log.Println("Remove work: ", Id)
    })


    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
	tmpl_add := template.Must(template.ParseFiles("form_add.html"))
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
	hmd5 := md5.Sum([]byte(dat))
	work.Id = fmt.Sprintf("%d-%x", utime, hmd5)
        dat, err = json.MarshalIndent(work, "", " ")
        if err != nil { fmt.Println(err) }
	fname := fmt.Sprintf("%s/%s", works_dir, work.Id)
        _ = ioutil.WriteFile(fname, dat, 0644)
//        fmt.Println(fname, string(dat))
        _ = work
        tmpl_add.Execute(w, struct{ Success bool }{true})
	log.Println("Add new work: ", work.Id)
    })


    http.ListenAndServe(":8080", nil)
}

