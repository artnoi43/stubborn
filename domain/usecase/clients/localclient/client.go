package localclient

type jsonClient struct{}

func New(tableFile string) *jsonClient {
	InitTable(tableFile)
	return &jsonClient{}
}

// Implement handler.dnsClient
func (k *jsonClient) Shutdown() {}
