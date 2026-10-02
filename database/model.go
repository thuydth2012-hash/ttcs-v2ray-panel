package database

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type Setting struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"uniqueIndex;not null"`
	Value string
}

type Inbound struct {
	ID uint `gorm:"primaryKey"`

	Remark   string `gorm:"not null"`
	Protocol string `gorm:"not null"`

	MonitorIP string
	Port      int

	TotalTraffic uint64
	Upload       uint64
	Download     uint64

	ExpireDate string

	// VMess
	UUID                      string
	AlterID                   int
	DisableInsecureEncryption bool

	// Transmission
	Transmission string

	// TCP
	TCPAcceptProxyProtocol bool

	// HTTP transport
	HTTPVersion string
	HTTPMethod  string
	HTTPPath    string
	HTTPHeaders string

	// TLS
	TLS            bool
	TLSDomain      string
	TLSCertPath    string
	TLSCertContent string

	Enabled bool
}
