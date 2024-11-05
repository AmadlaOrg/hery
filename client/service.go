package client

// NewClientService to set up the entity Client service
func NewClientService() IClient {
	return &SClient{}
}
