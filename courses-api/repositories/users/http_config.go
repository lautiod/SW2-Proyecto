package users

// HTTPConfig contiene la configuración para el cliente HTTP
type HTTPConfig struct {
	Host    string
	Port    string
	Timeout int    // en segundos
	APIKey  string // si es necesario
}
