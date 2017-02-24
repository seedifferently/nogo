package main

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func RunLocalDNSServer(laddr string, echo bool) (*dns.Server, string, error) {
	pc, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return nil, "", err
	}

	server := &dns.Server{PacketConn: pc, ReadTimeout: time.Hour, WriteTimeout: time.Hour}
	if echo { // Act as a simple echo server
		server.Handler = dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
			m := new(dns.Msg)
			m.SetReply(r)
			w.WriteMsg(m)
		})
	}

	waitLock := sync.Mutex{}
	waitLock.Lock()
	server.NotifyStartedFunc = func() { waitLock.Unlock() }

	go func() {
		server.ActivateAndServe()
		pc.Close()
	}()

	waitLock.Lock()
	return server, pc.LocalAddr().String(), nil
}

func Test_dnsHandler(t *testing.T) {
	db.Reset()

	s, addrstr, err := RunLocalDNSServer("127.0.0.1:0", false)
	if err != nil {
		t.Fatalf("unable to run test server: %v", err)
	}
	defer s.Shutdown()
	es, eaddrstr, err := RunLocalDNSServer("127.0.0.1:0", true)
	if err != nil {
		t.Fatalf("unable to run echo test server: %v", err)
	}
	defer es.Shutdown()

	*dnsProxyTo = eaddrstr
	dns.HandleFunc(".", dnsHandler)
	defer dns.HandleRemove(".")

	if err := db.put("test.disallowed", &Record{}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	if err := db.put("test.allowed", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}

	m := new(dns.Msg)
	m.SetQuestion("test.disallowed.", dns.TypeA)

	r, err := dns.Exchange(m, addrstr)
	if err != nil {
		t.Errorf("failed to exchange: %+v", err)
	}
	testEqual(t, "Disallowed Rcode = %+v, want %+v", r.Rcode, dns.RcodeNameError)
	testEqual(t, "Disallowed Authoritative = %+v, want %+v", r.Authoritative, true)
	testEqual(t, "Disallowed RecursionAvailable = %+v, want %+v", r.RecursionAvailable, false)
	testEqual(t, "Disallowed Questions = %+v, want %+v", r.Question, m.Question)

	m = new(dns.Msg)
	m.SetQuestion("test.allowed.", dns.TypeA)
	r, err = dns.Exchange(m, addrstr)
	if err != nil {
		t.Errorf("failed to exchange: %+v", err)
	}
	testEqual(t, "Allowed Rcode = %+v, want %+v", r.Rcode, dns.RcodeSuccess)
	testEqual(t, "Allowed Questions = %+v, want %+v", r.Question, m.Question)

	// Multiple questions
	m = new(dns.Msg)
	m.SetQuestion("test.allowed.", dns.TypeA)
	m.Question = append(m.Question, dns.Question{Name: "test.disallowed.", Qtype: dns.TypeA, Qclass: dns.ClassINET})
	r, err = dns.Exchange(m, addrstr)
	if err != nil {
		t.Errorf("failed to exchange: %+v", err)
	}
	testEqual(t, "Multiple Rcode = %+v, want %+v", r.Rcode, dns.RcodeSuccess)
	testEqual(t, "Multiple Response len() = %+v, want %+v", len(r.Question), 1)
	testEqual(t, "Multiple Response Question = %+v, want %+v", r.Question[0].Name, "test.allowed.")
}

func Test_filterQuestions(t *testing.T) {
	db.Reset()

	if err := db.put("test.allowed", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	if err := db.put("test.disallowed", &Record{}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}

	qs := filterQuestions([]dns.Question{{Name: "Test.Allowed."}, {Name: "Test.disallowed"}})
	testEqual(t, "len(filterQuestions(...)) = %+v, want %+v", len(qs), 1)
	testEqual(t, "filterQuestions(...)[0].Name = %+v, want %+v", qs[0].Name, "Test.Allowed.")
}

func Test_isNameAllowed(t *testing.T) {
	db.Reset()

	db.put("Test.allowed", &Record{Paused: true})
	db.put("Test.disallowed", &Record{})
	testEqual(t, "isNameAllowed('test.allowed') = %+v, want %+v", isNameAllowed("Test.Allowed."), true)
	testEqual(t, "isNameAllowed('test.disallowed') = %+v, want %+v", isNameAllowed("Test.Disallowed"), false)
	testEqual(t, "isNameAllowed('not.in.db') = %+v, want %+v", isNameAllowed("not.in.db."), true)
}
