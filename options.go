package conform

import "fmt"

type Option func(*loadOptions) error

type source interface {
	load(*schema) error
}

type loadOptions struct {
	sources []source
}

func FromEnv() Option {
	return func(o *loadOptions) error {
		if o == nil {
			return fmt.Errorf("%w: nil load options", ErrInvalidSchema)
		}

		o.sources = append(o.sources, envSource{})
		return nil
	}
}

func FromYAMLFile(path string) Option {
	return func(o *loadOptions) error {
		if o == nil {
			return fmt.Errorf("%w: nil load options", ErrInvalidSchema)
		}

		if path == "" {
			return fmt.Errorf("%w: yaml path must not be empty", ErrInvalidSchema)
		}

		o.sources = append(o.sources, yamlFileSource{
			path: path,
		})
		return nil
	}
}
