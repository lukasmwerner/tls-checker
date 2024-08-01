package scan

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/url"
	"sync"
)


func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS1.0"
	case tls.VersionTLS11:
		return "TLS1.1"
	case tls.VersionTLS12:
		return "TLS1.2"
	case tls.VersionTLS13:
		return "TLS1.3"
	case tls.VersionSSL30:
		return "TLS3.0"
	}
	return ""
}

type MinTLSVersionResult struct {
	Hostname    string
	MinVersion  string
	CipherSuite uint16
	Cert        *x509.Certificate
	DomainNames []string
}

func FindMinTLSVersion(link *url.URL) (MinTLSVersionResult, error) {
	r := MinTLSVersionResult{}
	host := link.Hostname()
	if p := link.Port(); p == "" {
		host += ":443"
	} else {
		host += ":" + p
	}
	r.Hostname = host
	for _, version := range []uint16{tls.VersionTLS10, tls.VersionTLS11, tls.VersionTLS12, tls.VersionTLS13} {
		conn, err := tls.Dial("tcp", host, &tls.Config{
			MinVersion:         version,
			MaxVersion:         version,
			InsecureSkipVerify: true,
		})
		if err != nil {
			fmt.Println("Failed", host, version)
			continue
		}
		r.MinVersion = tlsVersionString(conn.ConnectionState().Version)

		r.CipherSuite = conn.ConnectionState().CipherSuite
		r.Cert = conn.ConnectionState().PeerCertificates[0]
		r.DomainNames = r.Cert.DNSNames
		log.Println(host, "passed loop", r.MinVersion)
		break
	}
	return r, nil
}

func TestAllTLSVersions(link *url.URL) chan string {
	c := make(chan string)
	var wg sync.WaitGroup

	host := link.Hostname()
	if p := link.Port(); p == "" {
		host += ":443"
	} else {
		host += ":" + p
	}
	for _, version := range []uint16{tls.VersionTLS10, tls.VersionTLS11, tls.VersionTLS12, tls.VersionTLS13} {
		wg.Add(1)
		go func (version uint16) {	
			defer wg.Done()
			conn, err := tls.Dial("tcp", host, &tls.Config{
				MinVersion:         version,
				MaxVersion:         version,
				InsecureSkipVerify: true,
			})

			if err != nil {
				c <- fmt.Sprintf("err: failed %s %s", err.Error(), tlsVersionString(version))
				return	
			}
			defer conn.Close()

			c <- fmt.Sprintf("ok: %s", tlsVersionString(conn.ConnectionState().Version))
		}(version)
	}

	go func () {
		wg.Wait()
		close(c)
	}()

	return c
}
