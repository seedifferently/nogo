package main

import (
	"log"
	"strings"

	"github.com/miekg/dns"
)

func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	if !isDisabled {
		// Make a copy of the questions (in case we need them for the error response)
		qs := make([]dns.Question, len(r.Question))
		copy(qs, r.Question)

		// If none of the questions are allowed, respond with an error message
		if r.Question = filterQuestions(r.Question); len(r.Question) == 0 {
			m := &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Id:                 r.Id,
					Response:           true,
					Opcode:             dns.OpcodeQuery,
					Authoritative:      true,
					RecursionDesired:   r.RecursionDesired,
					RecursionAvailable: false,
					Rcode:              dns.RcodeNameError, // NXDOMAIN
				},
				Question: qs,
			}

			w.WriteMsg(m)
			return
		}
	}

	// Proxy allowed questions upstream
	for _, addr := range strings.Split(*dnsProxyTo, ",") {
		in, _, err := dnsClient.Exchange(r, addr)
		if err != nil {
			log.Printf("Exchange(%s) Error: %s\n", addr, err)
			continue
		}

		w.WriteMsg(in)
		return
	}

	dns.HandleFailed(w, r)
}

func filterQuestions(qs []dns.Question) []dns.Question {
	var keep []dns.Question

	for _, q := range qs {
		if isNameAllowed(q.Name) {
			keep = append(keep, q)
		}
	}

	return keep
}

func isNameAllowed(n string) bool {
	n = strings.TrimSuffix(n, ".")

	r, err := db.get(n)
	if err != nil {
		if err == errRecordNotFound {
			// If no record by that name was found, assume it is allowed
			return true
		}

		// For other errors, assume the name is now allowed
		log.Printf("db.get(%s) Error: %s\n", n, err)
		return false
	}

	return r.isAllowed()
}
