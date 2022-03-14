package handler

var dnsTypes = map[uint16]string{
	1:   "A",
	2:   "NS",
	5:   "CNAME",
	6:   "SOA",
	13:  "HINFO",
	15:  "MX",
	16:  "TXT",
	28:  "AAAA",
	46:  "RRSIG",
	257: "CAA",
}
