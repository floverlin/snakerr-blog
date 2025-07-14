package config

type Config struct {
	Server    ServerCf   `json:"server"`
	Database  DatabaseCf `json:"database"`
	Locales   LocalesCf  `json:"locales"`
	Templates TemplateCf `json:"templates"`
}

type ServerCf struct {
	Secret        string   `json:"secret"`
	Timeouts      int      `json:"timeouts"`
	WSPingRate    int      `json:"websocket_ping_rate"`
	WSPongTimeout int      `json:"websocket_pong_timeout"`
	FlashTimeout  int      `json:"flash_timeout"`
	StaticPath    string   `json:"static_path"`
	UploadsPath   string   `json:"uploads_path"`
	Visual        VisualCF `json:"visual"`
	Enviroment    string
}

type VisualCF struct {
	PostPerPage int `json:"post_per_page"`
}

type DatabaseCf struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Migrations string `json:"migrations"`
}

type LocalesCf struct {
	Path string `json:"path"`
}

type TemplateCf struct {
	Path string `json:"path"`
}
