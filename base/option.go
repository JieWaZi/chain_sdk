package base

// ClientOption describes a functional parameter for the client.
type ClientOption func(*Options)

// WithUser injects the user ID to the SDK.
func WithUser(username string) ClientOption {
	return func(opts *Options) {
		opts.Username = username
	}
}

// WithOrg injects the organization to the SDK.
func WithOrg(org string) ClientOption {
	return func(opts *Options) {
		opts.Organization = org
	}
}

// WithOrdererEndpoint sets the orderer endpoint.
func WithOrdererEndpoint(orderer string) ClientOption {
	return func(opts *Options) {
		opts.OrdererEndpoint = orderer
	}
}

// WithTargetEndpoint sets the orderer endpoint.
func WithTargetEndpoint(targetEndpoint []string) ClientOption {
	return func(opts *Options) {
		opts.TargetEndPoint = targetEndpoint
	}
}
