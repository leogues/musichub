package musichub

type Provider struct{}

// ProviderService represents a service that provides platforms.
type ProviderService interface {

	// Platforms returns a list of platforms.
	Platforms() []*Platform
}

type Platform struct {
	Id           string `json:"id"`
	TypePlatform string `json:"type"`
	Name         string `json:"name"`
	Label        string `json:"label"`
}
