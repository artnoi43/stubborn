package dohclient

type provider string

const (
	cloudflare provider = "CLOUDFLARE"
	quad9      provider = "QUAD9"
	google     provider = "GOOGLE"
	dnspod     provider = "DNSPOD"
)
