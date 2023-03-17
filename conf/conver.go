// Package conf 配置文件转换
package conf

import (
	"fmt"
	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"net"
	"strings"
)

// CoverConfigToRawConfig 未完成
func CoverConfigToRawConfig(cfg *config.Config) (*config.RawConfig, error) {
	conf := &config.RawConfig{}
	// General
	conf.AllowLan = cfg.General.AllowLan
	conf.BindAddress = cfg.General.BindAddress
	conf.Port = cfg.General.Port
	conf.SocksPort = cfg.General.SocksPort
	conf.RedirPort = cfg.General.RedirPort
	conf.MixedPort = cfg.General.MixedPort
	conf.TProxyPort = cfg.General.TProxyPort
	conf.ExternalController = cfg.General.ExternalController
	conf.Secret = cfg.General.Secret
	conf.LogLevel = cfg.General.LogLevel
	conf.IPv6 = cfg.General.IPv6
	conf.ExternalUI = cfg.General.ExternalUI
	conf.Mode = cfg.General.Mode
	conf.Interface = cfg.General.Interface
	conf.RoutingMark = cfg.General.RoutingMark
	// Proxies
	for _, proxy := range cfg.Proxies {

		switch proxy.Type() {
		case C.Direct:
		case C.Reject:
		case C.Shadowsocks:
			tp := make(map[string]any)
			conf.Proxy = append(conf.Proxy, tp)
		case C.ShadowsocksR:
			tp := make(map[string]any)
			buf, _ := proxy.MarshalJSON()
			tp[proxy.Name()] = string(buf)
			conf.Proxy = append(conf.Proxy, tp)
		case C.Snell:
			tp := make(map[string]any)
			buf, _ := proxy.MarshalJSON()
			tp[proxy.Name()] = string(buf)
			conf.Proxy = append(conf.Proxy, tp)
		case C.Socks5:
			tp := make(map[string]any)
			buf, _ := proxy.MarshalJSON()
			tp[proxy.Name()] = string(buf)
			conf.Proxy = append(conf.Proxy, tp)
		case C.Http:
			tp := make(map[string]any)
			buf, _ := proxy.MarshalJSON()
			tp[proxy.Name()] = string(buf)
			conf.Proxy = append(conf.Proxy, tp)
		case C.Vmess:
			tp := make(map[string]any)
			buf, _ := proxy.MarshalJSON()
			tp[proxy.Name()] = string(buf)
			conf.Proxy = append(conf.Proxy, tp)
		case C.Trojan:
			tp := make(map[string]any)
			buf, _ := proxy.MarshalJSON()
			tp[proxy.Name()] = string(buf)
			conf.Proxy = append(conf.Proxy, tp)
		case C.Selector:
		case C.Fallback:
		case C.URLTest:
		case C.LoadBalance:
		}
	}
	return conf, nil
}

// TransRule 内部的rules转换成配置文件yaml的 rules
func TransRule(cfg []C.Rule) []string {
	rules := make([]string, 0)
	for _, rule := range cfg {
		var str string
		switch rule.RuleType() {
		case C.Domain:
			str = "DOMAIN"
		case C.DomainSuffix:
			str = "DOMAIN-SUFFIX"
		case C.DomainKeyword:
			str = "DOMAIN-KEYWORD"
		case C.GEOIP:
			str = "GEOIP"
		case C.IPCIDR:
			// str = "IP-CIDR"  或者 IP-CIDR6
			if checkIpv6(rule.Payload()) {
				str = "IP-CIDR6"
			} else {
				str = "IP-CIDR"
			}
		case C.SrcIPCIDR:
			str = "SRC-IP-CIDR"
		case C.SrcPort:
			str = "SRC-PORT"
		case C.DstPort:
			str = "DST-PORT"
		case C.Process:
			str = "PROCESS-NAME"
		case C.ProcessPath:
			str = "PROCESS-PATH"
		case C.MATCH:
			str = "MATCH"
		}
		if str != "" {
			rules = append(rules, fmt.Sprintf("%s,%s,%s", str, rule.Payload(), rule.Adapter()))
		}
	}
	return rules
}

func getRuleTargetByRuleType(rt C.RuleType) string {
	var str string
	switch rt {
	case C.Domain:
		str = "DOMAIN"
	case C.DomainSuffix:
		str = "DOMAIN-SUFFIX"
	case C.DomainKeyword:
		str = "DOMAIN-KEYWORD"
	case C.GEOIP:
		str = "GEOIP"
	case C.IPCIDR:
		str = "IP-CIDR"
	case C.SrcIPCIDR:
		str = "SRC-IP-CIDR"
	case C.SrcPort:
		str = "SRC-PORT"
	case C.DstPort:
		str = "DST-PORT"
	case C.Process:
		str = "PROCESS-NAME"
	case C.ProcessPath:
	case C.MATCH:
		str = "MATCH"
	}
	return str
}

// 0: invalid ip
// 4: IPv4
// 6: IPv6
func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}

func checkIpv6(ip string) bool {
	return strings.Contains(ip, "::")
}
