package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var paramSNI, paramHost, paramUrl string
var paramHttps bool

func main() {
	flag.StringVar(&paramSNI, "sni", "", "Set SNI TLS header. If not set, will use URL.")
	flag.StringVar(&paramHost, "host", "", "Set Host HTTP header. If not set, will use URL.")
	flag.StringVar(&paramUrl, "url", "", "Set the URL. Required. Must not contain the protocol.")
	flag.BoolVar(&paramHttps, "https", true, "Enable or disable https usage.")
	flag.Parse()
	if paramUrl == "" {
		panic("No URL given. Try with --help.")
	}
	if strings.Contains(paramUrl, "://") {
		panic("No protocol please. example.com instead of http://example.com")
	}
	if paramSNI == "" {
		paramSNI = paramUrl
	}
	if paramHost == "" {
		paramHost = paramUrl
	}
	resp, err := frontCheck(paramUrl, paramSNI, paramHost, paramHttps)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Proto: %s\n", resp.Proto)
		content, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err.Error())
		}
		if len(content) > 1000 {
			content = content[:1000]
		}
		fmt.Printf("Up to first 1000 body bytes:\n%s\n", content)

	}
}

func frontCheck(url string, sniname string, hostheader string, https bool) (response http.Response, err error) {
	var client http.Client
	var proto string
	if https {
		proto = "https"
		client = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					ServerName: sniname,
				},
			},
		}
	} else {
		proto = "http"
		client = http.Client{}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s/", proto, url), nil)
	if err != nil {
		return http.Response{}, err
	}
	req.Host = hostheader
	resp, err := client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil
}
