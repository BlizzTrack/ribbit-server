/*
 * Copyright (c) 2020. BlizzTrack
 */

package network

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"time"
)

type BlizzardClient struct {
	timeout time.Duration
	proxy   string
	region  string
	dialer  net.Dialer
}

func NewBlizzardClient(region string, proxy string) *BlizzardClient {
	BlizzardClient := new(BlizzardClient)
	BlizzardClient.proxy = proxy
	BlizzardClient.region = region
	BlizzardClient.dialer = net.Dialer{Timeout: 5 * time.Second}
	BlizzardClient.timeout = 5 * time.Second

	return BlizzardClient
}

func (client *BlizzardClient) Call(method Command) (string, string, error) {
	data, err := client.call(method.String())
	if err != nil {
		return "", "", err
	}

	return data, getSeqn(data), nil
}

func (client *BlizzardClient) call(method string) (string, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(client.timeout))
	defer cancel()

	server := fmt.Sprintf("%s.version.battle.net:1119", client.region)
	if client.proxy != "" {
		server = client.proxy
	}

	ribbitClient, err := client.dialer.DialContext(ctx, "tcp", server)
	if err != nil {
		return "", err
	}
	defer ribbitClient.Close()

	err = ribbitClient.SetDeadline(time.Now().Add(client.timeout))
	if err != nil {
		return "", err
	}
	err = ribbitClient.SetReadDeadline(time.Now().Add(client.timeout))
	if err != nil {
		return "", err
	}
	err = ribbitClient.SetWriteDeadline(time.Now().Add(client.timeout))
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprintf(ribbitClient, fmt.Sprintf("%s\r\n", method))
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(ribbitClient)
	if err != nil {
		return "", err
	}

	content := string(data)

	return content, nil
}

func getSeqn(file string) string {
	f := getRegexParams(`\#\#\s?seqn\s?=\s?(?P<seqn>[0-9]*)`, file)

	if seqn, ok := f["seqn"]; ok {
		return seqn
	}
	return ""
}

func getRegexParams(regEx, url string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}