package model


//v1 - only piepline per repository
type Pipeline struct {
	Version    int            `yaml:"version" `
	Pipeline   map[string]Job `yaml:"pipeline"`
	Repository string
}


//one pipeline can have many jobs
type Job struct {
	Name        string       `yaml:"name,omitempty"`
	Language    string       `yaml:"language,omitempty"`
	Environment string       `yaml:"environment,omitempty"`
	Needs       []string     `yaml:"needs,omitempty"`
	SupplyChain *SupplyChain `yaml:"supply_chain,omitempty"`
}

type JobResult struct {
	Name string
	Err  string
}

type SupplyChain struct {
	Enabled  bool   `yaml:"enabled"`
	TokenRef string `yaml:"token_ref,omitempty"`
}

type Policies struct {
	AllowNetwork bool `yaml:"allow_network,omitempty"`
	AllowWriteFS bool `yaml:"allow_write_fs,omitempty"`
}

type LogEvent struct {
	Stream string // stdout/stderr
	Line   string
}

/*
type Pipeline struct {
    Version int             `yaml:"version"`
    Pipeline map[string]Job `yaml:"pipeline"`
}

type Job struct {
    Name        string            `yaml:"name,omitempty"`
    Intent      string            `yaml:"intent,omitempty"`
    Runtime     Runtime           `yaml:"runtime,omitempty"`
    Environment string            `yaml:"environment,omitempty"`
    DependsOn   []string          `yaml:"depends_on,omitempty"`
    Variables   map[string]string `yaml:"variables,omitempty"`
    Source      Source            `yaml:"source,omitempty"`
    Targets     Targets           `yaml:"targets,omitempty"`
    Permissions map[string]string `yaml:"permissions,omitempty"`
    Artifacts   Artifacts         `yaml:"artifacts,omitempty"`
    Cache       Cache             `yaml:"cache,omitempty"`
    Triggers    []string          `yaml:"triggers,omitempty"`
    Metadata    map[string]string `yaml:"metadata,omitempty"`
    Config      map[string]any    `yaml:"config,omitempty"`
}

type Runtime struct {
    Language string `yaml:"language,omitempty"`
    Version  string `yaml:"version,omitempty"`
    Image    string `yaml:"image,omitempty"`
}

type Source struct {
    Checkout  bool   `yaml:"checkout,omitempty"`
    Repository string `yaml:"repository,omitempty"`
    Ref        string `yaml:"ref,omitempty"`
}

type Targets struct {
    Platform     []string              `yaml:"platform,omitempty"`
    Architecture []string              `yaml:"architecture,omitempty"`
    Runtime      map[string][]string   `yaml:"runtime,omitempty"`
    Custom       map[string][]string   `yaml:"custom,omitempty"`
}

type Artifacts struct {
    Inputs  []string `yaml:"inputs,omitempty"`
    Outputs []string `yaml:"outputs,omitempty"`
}

type Cache struct {
    Enabled bool     `yaml:"enabled,omitempty"`
    Paths   []string `yaml:"paths,omitempty"`
}
*/
