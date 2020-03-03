// Copyright (C) liasica. 2020-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2020-03-03
// Based on proxy by liasica, magicrolan@qq.com.

package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "os/exec"
    "runtime"
    "strings"
)

var (
    h bool
    u string
    p string
)

func init() {
    flag.BoolVar(&h, "h", false, "this help")
    flag.StringVar(&u, "u", "", "请输入你的订阅地址")
    flag.StringVar(&p, "p", "8123", "请输入代理端口，默认为 8123")
}

func openbrowser(url string) {
    var err error

    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = fmt.Errorf("unsupported platform")
    }
    if err != nil {
        log.Fatal(err)
    }

}

func readInput() string {
    reader := bufio.NewReader(os.Stdin)
    text, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalln(err)
    }
    return strings.ReplaceAll(text, "\n", "")
}

func getIP() string {
    var ip net.IP
    ifaces, _ := net.Interfaces()
    // handle err
    for _, i := range ifaces {
        addrs, _ := i.Addrs()
        // handle err
        for _, addr := range addrs {
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            // process IP address
        }
    }
    if ip == nil {
        return ""
    }
    return ip.String()
}

func main() {
    // flag.Parse()
    //
    // if h || u == "" || p == "" {
    //     flag.Usage()
    //     return
    // }

    fmt.Print("请输入需代理的url: ")
    u = readInput()

    fmt.Print("请输入代理端口(默认是8123): ")
    p = readInput()
    if p == "" {
        p = "8123"
    }
    remote, err := url.Parse(u)
    if err != nil {
        log.Fatalln(err)
    }
    target, _ := url.Parse(remote.Scheme + "://" + remote.Host)

    http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {

        proxy := httputil.NewSingleHostReverseProxy(target)

        req.URL.Host = remote.Host
        req.URL.Scheme = remote.Scheme
        req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
        req.Host = remote.Host

        proxy.ServeHTTP(res, req)
    })

    // http.HandleFunc("/body", func(res http.ResponseWriter, req *http.Request) {
    //
    // })
    ip := getIP()
    if ip == "" {
        ip = "你的本地IP"
    }

    q := remote.RawQuery
    if q != "" {
        q = "?" + q
    }

    fmt.Printf("路由器订阅地址填写 http://%s:%s%s%s\n订阅完成后可关闭，下次订阅或更新前请重新运行", ip, p, remote.Path, q)
    log.Fatalln(http.ListenAndServe(":"+p, nil))
}
