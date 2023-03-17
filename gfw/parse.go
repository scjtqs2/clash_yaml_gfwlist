// Package gfw 解析gfw
package gfw

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

func LoadGfwList() ([]string, error) {
	err := InitGfwList()
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(Gfwlist)
	if err != nil {
		return nil, err
	}
	orgBuf, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	return readList(orgBuf), nil
}

func readList(data []byte) []string {
	var list []string
	tmpMap := make(map[string]bool)
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(dst, data)
	if err == nil {
		data = dst[:n]
	}

	if !bytes.HasPrefix(data, []byte("[AutoProxy ")) {
		log.Fatal("invalid auto proxy file")
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.ContainsRune(line, '.') {
			continue
		}
		domain := line
		switch line[0] {
		case '.':
			domain = "|h://" + line[1:]
			fallthrough
		case '|':
			if line[1] == '|' {
				domain = strings.TrimRight(domain[2:], "/")
				if strings.ContainsRune(domain, '/') {
					log.Printf("unsupported line: %s", line)
					continue
				}
			} else {
				u, err := url.Parse(strings.Replace(domain[1:], "*", "/", -1))
				if err != nil || !strings.ContainsRune(u.Host, '.') || strings.ContainsRune(u.Host, ':') {
					log.Printf("unsupported line: %s", line)
					continue
				}
				domain = u.Host
			}
			if net.ParseIP(domain) == nil {
				// list = append(list, domain)
				tmpMap[domain] = true
			}
		}
	}
	for s := range tmpMap {
		list = append(list, s)
	}
	return list
}
