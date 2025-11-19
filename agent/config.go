package agent

// AgentConfig 代理配置
type AgentConfig struct {
	UseLLM       bool
	AutoExecute  bool
	WorkspaceDir string
	ReportDir    string
	ADBPath      string
	RemoteDir    string
}

func (cfg *AgentConfig) normalize() {
	if cfg.WorkspaceDir == "" {
		cfg.WorkspaceDir = "./workspace"
	}
	if cfg.ReportDir == "" {
		cfg.ReportDir = "./workspace/reports"
	}
}

