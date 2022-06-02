package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
//    "time"
    "log"
    "os"
    "os/exec"
    "strings"
)

var conf_dir = ".cs-server/"
var works_dir = conf_dir + "works/"
var archiv_dir = conf_dir + "archiv/"
var run_dir = conf_dir + "run/"

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

func LoadWork(id string) (Work, map[string]interface{}) {
    var w Work
    fname := fmt.Sprintf("%s/%s", works_dir, id)
    content, err := ioutil.ReadFile(fname)
    if err != nil { /*log.Fatal("Error when opening file: ", err)*/ return w, nil }
    var data map[string]interface{}
    err = json.Unmarshal(content, &data)
    if err != nil { /*log.Fatal("Error during Unmarshal(): ", err)*/ return w, nil }
    w.Id = id
    w.Title = fmt.Sprintf("%s", data["Title"])
    w.Temp = fmt.Sprintf("%s", data["Temp"])
    w.Pressure = fmt.Sprintf("%s", data["Pressure"])
    w.Stat = fmt.Sprintf("%s", data["Stat"])
    w.Compounds = fmt.Sprintf("%s", data["Compounds"])
    w.Status = fmt.Sprintf("%s", data["Status"])
    w.Done = false
    return w, data
}

func msys3(comstr string) string {
//    arg := strings.Split(strconv.Unquote(comstr, '"'))
    arg := strings.Split(comstr, " ")
    out, err := exec.Command(arg[0], arg[1:]...).Output()
    if err != nil {
        log.Fatal(err)
    }
    return string(out)
}

func msys2(prg string, args ...string) string {
    out, err := exec.Command(prg, args...).Output()
    if err != nil {
        log.Fatal(err)
    }
    return string(out)
}

func msys(prg string, args ...string) string {
    var out = make([]byte, 0)
    var err error = nil
    switch len(args) {
	case 0 : { out, err = exec.Command(prg).Output() }
	case 1 : { out, err = exec.Command(prg, args[0]).Output() }
	case 2 : { out, err = exec.Command(prg, args[0], args[1]).Output() }
	case 3 : { out, err = exec.Command(prg, args[0], args[1], args[2]).Output() }
	default : { out, err = exec.Command(prg).Output() }
    }
    if err != nil {
        log.Fatal(err)
    }
    return string(out)
}


func sysrun(app string, args ...string) {
    out, err := exec.Command(app, args...).Output()
    if err != nil { /* log.Fatal(err) */fmt.Println(err) }
    fmt.Println(string(out))
}

func MakeWork(w Work, smi string) { // prepared
    dir := run_dir + w.Id + "/"
    _ = os.Mkdir(dir, 0755)
//    if err != nil { /* log.Fatal(err) */ return }
    _ = ioutil.WriteFile(dir + "mol.can", []byte(smi), 0644)
//    com := fmt.Sprintf("babel", "-ican", dir+"mol.can", "-oxyz", "--gen3D", dir+"v1.xyz", dir, dir)
//    fmt.Println(com)
    sysrun("./prepare-work.sh", run_dir, w.Id, w.Stat)
}

func CheckWorkList() {
    files, err := ioutil.ReadDir(works_dir)
    if err != nil { /*log.Fatal(err)*/ return }
    for _, file := range files {
	w, _ := LoadWork(file.Name())
	if w.Status=="wait" {
	    Compounds := strings.Split(w.Compounds, "\r\n")
	    smi := Compounds[0]
	    for i := range Compounds[1:] {
		smi = smi + "." + Compounds[i]
	    }
	    fmt.Println(w.Id, w.Title)
	    for i := range Compounds {
		Compound := Compounds[i]
		if len(Compound)<1 { continue }
		fmt.Println("    ", Compound)
	    }
	    MakeWork(w, smi)
	    log.Println("Prepare work:", w.Id)
	}
    }
    return
}

func main() {
    file, err := os.OpenFile(conf_dir + "cs-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    // if err != nil { log.Fatal(err) }
    if err == nil { log.SetOutput(file) }

    CheckWorkList()
}
