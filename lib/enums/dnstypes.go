package enums

// DnsTypes references https://en.wikipedia.org/wiki/List_of_DNS_record_types
var DnsTypes = map[uint16]string{
	1:   "A",
	2:   "NS",
	5:   "CNAME",
	6:   "SOA",
	12:  "PTR",
	13:  "HINFO",
	15:  "MX",
	16:  "TXT",
	18:  "AFSDB",
	28:  "AAAA",
	46:  "RRSIG",
	257: "CAA",
}
