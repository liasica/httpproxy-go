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
    "flag"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

var (
    h    bool
    u  string
    p string
)

func init() {
    flag.BoolVar(&h, "h", false, "this help")
    flag.StringVar(&u, "u", "", "请输入你的订阅地址")
    flag.StringVar(&p, "p", "8123", "请输入代理端口，默认为 8123")
}

func main() {
    flag.Parse()

    if h || u == "" || p == "" {
        flag.Usage()
        return
    }

    log.Println("订阅地址是:", u)
    log.Println("代理端口是:", p)

    http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
        remote, err := url.Parse(u)
        if err != nil {
            log.Fatalln(err)
        }
        proxy := httputil.NewSingleHostReverseProxy(remote)

        req.URL.Host = remote.Host
        req.URL.Scheme = remote.Scheme
        req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
        req.Host = remote.Host

        proxy.ServeHTTP(res, req)
    })

    log.Println("代理启动")
    log.Fatalln(http.ListenAndServe(":" + p, nil))
}
