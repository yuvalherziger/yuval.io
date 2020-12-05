package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "strings"
)

type Payload struct {
    Input  string `json:"input"`
    Output string `json:"output"`
    Width  int    `json:"width"`
}

func executeCommand(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodPost {
        http.NotFound(w, r)
    }
    var p Payload
    err := json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var stdout, stderr bytes.Buffer

    cmdArgs := strings.Split(p.Input, " ")
    if cmdArgs[0] != "bio" {
        p.Output = fmt.Sprintf("Unsupported command: %s", cmdArgs[0])
        pJson, _ := json.Marshal(p)
        _, _ = w.Write(pJson)
        return
    }
    cmdArgs = append(cmdArgs, fmt.Sprintf("--view-width=%d", p.Width))
    cmd := exec.Command("/usr/bin/bio", cmdArgs[1:]...)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err = cmd.Run()


    if err != nil {
        log.Printf(err.Error())
        p.Output = stderr.String()
    } else {
        p.Output = stdout.String()
    }
    pJson, _ := json.Marshal(p)
    _, err = w.Write(pJson)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    fs := http.FileServer(http.Dir("/opt/bio/app/static/"))
    http.Handle("/", fs)
    http.HandleFunc("/api/v1beta1/cmd", executeCommand)
    err := http.ListenAndServe(":8090", nil)
    if err != nil {
        log.Fatal(err)
    }
}
