package plugin

// RunService initializes and runs a service using the provided ServiceAPI
// implementation.
func RunService(api ServiceAPI) {
	s, err := NewService(api)
	if err != nil {
		Error(err)
		return
	}

	if err := s.Run(); err != nil {
		Error(err)
		return
	}
}

// RunFeature initializes a feature using the provided FeatureAPI and executes it.
func RunFeature(api FeatureAPI) {
	f, err := NewFeature(api)
	if err != nil {
		Error(err)
		return
	}

	if err := f.Run(); err != nil {
		Error(err)
		return
	}
}
