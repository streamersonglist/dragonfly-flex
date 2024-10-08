package fly

type ErrorResponse struct {
	Error string `json:"error"`
}

type Machine struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Region    string        `json:"region"`
	State     string        `json:"state"`
	Config    MachineConfig `json:"config"`
	PrivateIP string        `json:"private_ip"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

type MachineConfig struct {
	Env struct {
		PrimaryRegion string `json:"PRIMARY_REGION"`
	} `json:"env"`
	Init     map[string]interface{} `json:"init"`
	Guest    Guest                  `json:"guest"`
	Metadata map[string]string      `json:"metadata"`
	Mounts   []Mount                `json:"mounts"`
	Services []Service              `json:"services"`
	Metrics  Metrics                `json:"metrics"`
	Checks   map[string]HTTPCheck   `json:"checks"`
	Image    string                 `json:"image"`
	Restart  RestartPolicy          `json:"restart"`
}

type Guest struct {
	CpuKind  string `json:"cpu_kind"`
	Cpus     int    `json:"cpus"`
	MemoryMb int    `json:"memory_mb"`
}

type Mount struct {
	Encrypted bool   `json:"encrypted"`
	Path      string `json:"path"`
	SizeGb    int    `json:"size_gb"`
	Volume    string `json:"volume"`
	Name      string `json:"name"`
}

type Service struct {
	Protocol         string      `json:"protocol"`
	InternalPort     int         `json:"internal_port"`
	Autostart        bool        `json:"autostart"`
	Ports            []Port      `json:"ports"`
	Concurrency      Concurrency `json:"concurrency"`
	ForceInstanceKey interface{} `json:"force_instance_key"`
}

type Port struct {
	Port     int      `json:"port"`
	Handlers []string `json:"handlers"`
}

type Concurrency struct {
	Type      string `json:"type"`
	HardLimit int    `json:"hard_limit"`
	SoftLimit int    `json:"soft_limit"`
}

type Metrics struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}

type HTTPCheck struct {
	Port     int    `json:"port"`
	Type     string `json:"type"`
	Interval string `json:"interval"`
	Timeout  string `json:"timeout"`
	Path     string `json:"path"`
}

type RestartPolicy struct {
	Policy string `json:"policy"`
}
type CreateMachineRequest struct {
	Config                  MachineConfig `json:"config"`
	LeaseTTL                int           `json:"lease_ttl"`
	Lsvd                    bool          `json:"lsvd"`
	Name                    string        `json:"name,omitempty"`
	Region                  string        `json:"region,omitempty"`
	SkipLaunch              bool          `json:"skip_launch"`
	SkipServiceRegistration bool          `json:"skip_service_registration"`
}

type UpdateMachineRequest struct {
	Config                  MachineConfig `json:"config"`
	CurrentVersion          string        `json:"current_version"`
	LeaseTTL                int           `json:"lease_ttl"`
	Lsvd                    bool          `json:"lsvd"`
	Name                    string        `json:"name,omitempty"`
	Region                  string        `json:"region,omitempty"`
	SkipLaunch              bool          `json:"skip_launch"`
	SkipServiceRegistration bool          `json:"skip_service_registration"`
}

type MachineEvent struct {
	ID        string      `json:"id"`
	Request   interface{} `json:"request"`
	Source    string      `json:"source"`
	Status    string      `json:"status"`
	Timestamp int         `json:"timestamp"`
	Type      string      `json:"type"`
}

type Lease struct {
	Description string `json:"description"`
	ExpiresAt   int    `json:"expires_at"`
	Nonce       string `json:"nonce"`
	Owner       string `json:"owner"`
	Version     string `json:"version"`
}

type CreateLeaseRequest struct {
	Description string `json:"description"`
	TTL         int    `json:"ttl"` // `ttl` is the time in seconds for which the lease will be valid
}

type MachineExecRequest struct {
	Command []string `json:"command"`
	Timeout int      `json:"timeout"`
}

type SignalRequest struct {
	Signal string `json:"signal"`
}

type StopRequest struct {
	Signal  string   `json:"signal"`
	Timeout Duration `json:"timeout"`
}

type Duration struct {
	// Define Duration fields
}

type ListenSocket struct {
	Address string `json:"address"`
	Proto   string `json:"proto"`
}

type ProcessStat struct {
	Command       string         `json:"command"`
	CPU           int            `json:"cpu"`
	Directory     string         `json:"directory"`
	ListenSockets []ListenSocket `json:"listen_sockets"`
	PID           int            `json:"pid"`
	RSS           int            `json:"rss"`
	Rtime         int            `json:"rtime"`
	Stime         int            `json:"stime"`
}

type ExecResponse struct {
	ExitCode   int    `json:"exit_code"`
	ExitSignal int    `json:"exit_signal"`
	Stderr     string `json:"stderr"`
	Stdout     string `json:"stdout"`
}
