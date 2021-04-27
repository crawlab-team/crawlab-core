package grpc

type Address struct {
	Host string
	Port string
}

type AddressOptions struct {
	Host string
	Port string
}

func NewAddress(opts *AddressOptions) Address {
	if opts == nil {
		opts = &AddressOptions{}
	}
	if opts.Host == "" {
		opts.Host = "localhost"
	}
	if opts.Port == "" {
		opts.Port = "8765"
	}
	return Address{
		Host: opts.Host,
		Port: opts.Port,
	}
}
