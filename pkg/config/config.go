package config

type Config struct {
	Server struct {
		Port      int    `toml:"port"`
		Algorithm string `toml:"algorithm"`
	} `toml:"server"`
	Nodes []struct {
		ID     string `toml:"id"`
		URL    string `toml:"url"`
		MaxBPM int    `toml:"max_bpm"`
		MaxRPM int    `toml:"max_rpm"`
	} `toml:"nodes"`
}

func LoadConfig(path string) (*Config, error) {

	config := &Config{
		Server: struct {
			Port      int    `toml:"port"`
			Algorithm string `toml:"algorithm"`
		}{
			Port:      8081,
			Algorithm: "round_robin",
		},
		Nodes: []struct {
			ID     string `toml:"id"`
			URL    string `toml:"url"`
			MaxBPM int    `toml:"max_bpm"`
			MaxRPM int    `toml:"max_rpm"`
		}{
			{
				ID:     "node1",
				URL:    "http://localhost:8082",
				MaxBPM: 1000,
				MaxRPM: 100,
			},
			{
				ID:     "node2",
				URL:    "http://localhost:8083",
				MaxBPM: 1000,
				MaxRPM: 100,
			},
			{
				ID:     "node3",
				URL:    "http://localhost:8084",
				MaxBPM: 1000,
				MaxRPM: 100,
			},
		},
		// Nodes: []struct {
		// 	ID     string `toml:"id"`
		// 	URL    string `toml:"url"`
		// 	MaxBPM int    `toml:"max_bpm"`
		// 	MaxRPM int    `toml:"max_rpm"`
		// }{
		// 	{
		// 		ID:     "api-server",
		// 		URL:    "http://localhost:8082",
		// 		MaxBPM: 30,
		// 		MaxRPM: 3,
		// 	},
		// },
	}

	return config, nil
}
