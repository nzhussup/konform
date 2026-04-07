package conform

func Load(target any, opts ...Option) error {
	loadOpts := loadOptions{}

	for _, opt := range opts {
		if err := opt(&loadOpts); err != nil {
			return err
		}
	}

	sc, err := buildSchema(target)
	if err != nil {
		return err
	}

	if err := applyDefaults(sc); err != nil {
		return err
	}

	for _, src := range loadOpts.sources {
		if err := src.load(sc); err != nil {
			return err
		}
	}

	if err := validateRequired(sc); err != nil {
		return err
	}

	return nil
}
