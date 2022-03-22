package requests

import (
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

const MaxDuration = 1<<31 - 1

var (
	localTCPAddrList []*net.TCPAddr
	ProxyAddr        string
)

func SetLocalTCPAddrList(ips ...string) {
	list := make([]*net.TCPAddr, 0, len(ips))
	for k := range ips {
		p := net.ParseIP(ips[k])
		if p == nil {
			continue
		}
		list = append(list, &net.TCPAddr{IP: p})
	}
	localTCPAddrList = list
}

func proxyFunc(req *http.Request) (*url.URL, error) {
	u, err := checkProxyAddr(ProxyAddr)
	if err != nil {
		return http.ProxyFromEnvironment(req)
	}
	return u, err
}

func getLocalTCPAddr() *net.TCPAddr {
	if len(localTCPAddrList) == 0 {
		return nil
	}
	i := rand.Intn(len(localTCPAddrList))
	return localTCPAddrList[i]
}

func getDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		LocalAddr: getLocalTCPAddr(),
	}
}

func checkProxyAddr(proxyAddr string) (u *url.URL, err error) {
	if proxyAddr == "" {
		return nil, NewErrProxyAddrEmpty()
	}
	host, port, err := net.SplitHostPort(proxyAddr)
	if err == nil {
		u = &url.URL{Host: net.JoinHostPort(host, port)}
		return
	}
	u, err = url.Parse(proxyAddr)
	if err == nil {
		return
	}
	return
}

func SetGlobalProxy(proxyAddr string) {
	ProxyAddr = proxyAddr
}

// SetTCPHostBind 设置host绑定ip
func SetTCPHostBind(host, ip string) {
	//tcpCache.Store(host, expires.NewDataExpires(net.ParseIP(ip), MaxDuration))
	//return
}
