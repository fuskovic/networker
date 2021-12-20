package test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TlsSuite struct {
	RawPrivateKey string
	RawCert       string
	Cert          tls.Certificate
	CACert        []byte
	ClientConfig  *tls.Config
	ServerConfig  *tls.Config
}

type HitCounter interface {
	HitCount() int
}

type MockTlsServer interface {
	HitCounter
	SetHandler(http.Handler)
	StartTLS()
	URL() string
	Close()
}

type mockTlsServer struct {
	*httptest.Server
	hitCount int
}

// ServeHTTP implements http.Handler
func (s *mockTlsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.hitCount++
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello from %s", r.Host)
}

// HitCount returns the number of requests the server has processed and implements HitCounter.
func (s *mockTlsServer) HitCount() int {
	return s.hitCount
}

// URL returns the embedded testservers url and implements MockTlsServer.
func (s *mockTlsServer) URL() string {
	return s.Server.URL
}

// Close shuts down the server and implements MockTlsServer.
func (s *mockTlsServer) Close() {
	s.Server.Close()
}

// Config returns the underlying test server configuration and implements MockTlsServer
func (s *mockTlsServer) SetHandler(h http.Handler) {
	s.Server.Config.Handler = h
}

// NewUnstartedMockTlsServer returns a MockTlsServer implementation
// that the caller is responsible for starting and tearing down.
func (s *TlsSuite) NewUnstartedMockTlsServer(t *testing.T) MockTlsServer {
	t.Helper()
	tlsServer := new(mockTlsServer)
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			tlsServer.hitCount++
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "hello from %s", r.Host)
		}),
	)
	testServer.TLS = &tls.Config{
		Certificates: []tls.Certificate{s.Cert},
	}
	tlsServer.Server = testServer
	return tlsServer
}

// NewStartedMockTlsServers returns n amount of MockTlsServer implementations that have already been started.
// The caller should defer the returned cleanup func on initialization to tear down the servers after the test.
func (s *TlsSuite) NewStartedMockTlsServers(t *testing.T, n int) ([]MockTlsServer, func()) {
	var tlsServers []MockTlsServer
	var cleanupFuncs []func()
	for len(tlsServers) < n {
		tlsServer := s.NewUnstartedMockTlsServer(t)
		tlsServer.StartTLS()
		tlsServers = append(tlsServers, tlsServer)
		cleanupFuncs = append(cleanupFuncs, func() { tlsServer.Close() })
	}
	return tlsServers, func() {
		for _, f := range cleanupFuncs {
			f()
		}
	}
}

func newTlsTestSuite(t *testing.T) *TlsSuite {
	t.Helper()

	subjectName := pkix.Name{
		Organization:  []string{"Networker"},
		Country:       []string{"US"},
		Province:      []string{""},
		Locality:      []string{"Austin"},
		StreetAddress: []string{"The Whitley"},
		PostalCode:    []string{"78701"},
	}

	now := time.Now()

	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               subjectName,
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// generate keypair
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	require.NoError(t, err)

	// pem encode
	caPEM := new(bytes.Buffer)
	require.NoError(t, pem.Encode(caPEM,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		},
	))

	caPrivKeyPEM := new(bytes.Buffer)
	require.NoError(t, pem.Encode(caPrivKeyPEM,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
		},
	))

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject:      subjectName,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    now,
		NotAfter:     now.AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	require.NoError(t, err)

	certPEM := new(bytes.Buffer)
	require.NoError(t, pem.Encode(certPEM,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		},
	))

	certPrivKeyPEM := new(bytes.Buffer)
	require.NoError(t, pem.Encode(certPrivKeyPEM,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
		},
	))

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	require.NoError(t, err)

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPEM.Bytes())

	// initialize the test suite
	return &TlsSuite{
		RawPrivateKey: certPrivKeyPEM.String(),
		RawCert:       certPEM.String(),
		Cert:          serverCert,
		CACert:        caPEM.Bytes(),
		ClientConfig:  &tls.Config{RootCAs: certpool},
		ServerConfig: &tls.Config{
			Certificates: []tls.Certificate{serverCert},
		},
	}
}
