package main

import (
    "os"
    "math/rand"
    "time"
    //"fmt"
)

func main() {
    filename := os.Args[1]

    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()


    lines := []string{
        "127.0.0.1 - name1 [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 1234\n",
        "127.0.0.2 - name2 [09/May/2018:16:00:39 +0000] \"POST /api/users HTTP/1.0\" 304 1234\n",
        "127.0.0.3 - name3 [09/May/2018:16:00:39 +0000] \"PUT /api/metric HTTP/1.0\" 500 1234\n",
        "127.0.0.4 - name4 [09/May/2018:16:00:39 +0000] \"GET / HTTP/1.0\" 404 1234\n",
    }

    s := rand.NewSource(time.Now().Unix())
    r := rand.New(s)

    for i := 0; i < 100000; i++ {
        idx := r.Intn(len(lines))

        if _, err = f.WriteString(lines[idx]); err != nil {
            panic(err)
        }
    }
}