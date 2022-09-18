package main

import (
    "fmt"
//    "io/ioutil"
//    "time"
    "log"
//    "os"
    "os/exec"
    "strings"
    "regexp"
)

var conf_dir = "~/QCS/"
var works_dir = conf_dir + "works/"
var archiv_dir = conf_dir + "results/"
var run_dir = conf_dir + "run/"

func msys3(comstr string) string {
//    arg := strings.Split(strconv.Unquote(comstr, '"'))
    arg := strings.Split(comstr, " ")
    out, err := exec.Command(arg[0], arg[1:]...).Output()
    if err != nil {
        log.Fatal(err)
    }
    return string(out)
}

func main() {
    r := regexp.MustCompile(`[ns]`)
    pids := msys3("pgrep qcs-run-cpu01").split()
    fmt.Println(pids)
}
