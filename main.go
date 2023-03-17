// Package 入口
package main

import (
	conf2 "clash_yaml_gfwlist/conf"
	"clash_yaml_gfwlist/gfw"
	"flag"
	"fmt"
	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	rules "github.com/Dreamacro/clash/rule"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	c string
	h bool
	// debug   bool
	d       bool
	Version string
	o       string
)

func init() {

	flag.StringVar(&c, "c", "config.yaml", "待转换的clash配置文件")
	flag.StringVar(&o, "o", "config_gfw.yaml", "转换完成后的clash新配置文件")
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&d, "d", true, "是否每次运行都重新下载gfwlist.txt")
	flag.StringVar(&gfw.GfwlistUrl, "gfw", "https://pagure.io/gfwlist/raw/master/f/gfwlist.txt", "gfwlist download url for https://github.com/gfwlist/gfwlist")
	flag.Parse()
}

// Help cli命令行-h的帮助提示
func Help() {
	fmt.Printf(`clash 规则添加gfwlist命令工具
version: %s
Usage:
server [OPTIONS]
Options:
`, Version)
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	if h {
		Help()
	}
	if d {
		_ = os.RemoveAll(gfw.Gfwlist)
	}
	confBuf, err := readConfig(c)
	if err != nil {
		panic(err)
	}
	var rawConf config.RawConfig
	err = yaml.Unmarshal(confBuf, &rawConf)
	if err != nil {
		panic(err)
	}
	conf, err := config.Parse(confBuf)
	if err != nil {
		panic(err)
	}
	domainList, err := gfw.LoadGfwList()
	if err != nil {
		panic(err)
	}
	// 判断是否有 🈲 GFW
	isGfw := false
	for _, provider := range conf.Providers {
		if provider.Name() == "🈲 GFW" {
			isGfw = true
			break
		}
	}
	if !isGfw {
		gfwGroup := make(map[string]any)
		gfwGroup["name"] = "🈲 GFW"
		gfwGroup["type"] = "select"
		gfwGroup["proxies"] = []string{"🚀 节点选择", "♻️ 自动选择", "🎯 全球直连"}
		rawConf.ProxyGroup = append(rawConf.ProxyGroup, gfwGroup)
	}
	tmpRules := make([]C.Rule, 0)
	for _, s := range domainList {
		tmpRules = append(tmpRules, rules.NewDomainSuffix(s, "🈲 GFW"))
	}
	// 去除重复的rules
	for _, rule := range conf.Rules {
		switch rule.RuleType() {
		case C.DomainSuffix:
			if !slices.Contains(domainList, rule.Payload()) {
				tmpRules = append(tmpRules, rule)
			}
		default:
			tmpRules = append(tmpRules, rule)
		}
	}
	rawConf.Rule = conf2.TransRule(tmpRules)
	buf, _ := yaml.Marshal(rawConf)
	err = writefile(o, buf)
	if err != nil {
		panic(err)
	}
}

// readConfig 读去文件
func readConfig(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file %s is empty", path)
	}

	return data, err
}

// writefile 写入文件
func writefile(file string, buf []byte) error {
	return os.WriteFile(file, buf, 0o644)
}
