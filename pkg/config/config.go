package config

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config groups all of the global configuration parameters into
// a single struct
type Config struct {
	AuthEndpoint string        `mapstructure:"authEndpoint"`
	AuthPath     string        `mapstructure:"authPath"`
	AuthTimeout  time.Duration `mapstructure:"authTimeout"`
	TokenPath    string        `mapstructure:"tokenPath"`
	RedirectURL  string        `mapstructure:"redirectURL"`

	Verbosity string `mapstructure:"verbosity"`

	MultipathUploadURI string        `mapstructure:"multipathUploadURI"`
	GetOrgURI          string        `mapstructure:"getOrgURI"`
	LogsURI            string        `mapstructure:"logsURI"`
	StatusURI          string        `mapstructure:"statusURI"`
	TestsURI           string        `mapstructure:"testsURI"`
	SchemaURL          string        `mapstructure:"schemaURL"`
	WBHost             string        `mapstructure:"wbHost"`
	APITimeout         time.Duration `mapstructure:"apiTimeout"`
	BiomeDNSZone       string        `mapstructure:"biomeDNSZone"`
	VersionLocation    string        `mapstructure:"versionLocation"`
	CLIURL             string        `mapstructure:"cliURL"`
	StopTestURI        string        `mapstructure:"stopTestURI"`
	StopDefURI         string        `mapstructure:"stopDefURI"`
	RunTestURI         string        `mapstructure:"runTestURI"`
	CreateTestURI      string        `mapstructure:"createTestURI"`

	OrgID string `mapstructure:"orgID"`

	TokenFile          string `mapstructure:"tokenFile"`
	OrgFile            string `mapstructure:"orgFile"`
	SchemaFile         string `mapstructure:"schemaFile"`
	GenesisCredentials string `mapstructure:"genesisCredentials"`
	GenesisBanner      string `mapstructure:"genesisBanner"`

	Dir     configdir.ConfigDir `mapstructure:"-"`
	UserDir *configdir.Config   `mapstructure:"-"`
}

func (c Config) HTTPClient() *http.Client {
	return &http.Client{Timeout: c.APITimeout}
}

func (c Config) APIEndpoint() string {
	if !strings.HasPrefix(c.WBHost, "http") {
		return "https://" + c.WBHost
	}
	return c.WBHost
}

// GetLogger gets a logger according to the config
func (c Config) setupLogrus() {
	lvl, err := logrus.ParseLevel(c.Verbosity)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(lvl)
	}

	return
}

func setViperEnvBindings() {
	viper.BindEnv("authEndpoint", "AUTH_ENDPOINT")
	viper.BindEnv("authPath", "AUTH_PATH")
	viper.BindEnv("tokenPath", "TOKEN_PATH")
	viper.BindEnv("verbosity", "VERBOSITY")
	viper.BindEnv("redirectURL", "REDIRECT_URL")
	viper.BindEnv("authTimeout", "AUTH_TIMEOUT")
	viper.BindEnv("MultipathUploadURI", "MULTIPART_UPLOAD_URI")
	viper.BindEnv("GetOrgURI", "GET_ORG_URI")
	viper.BindEnv("LogsURI", "LOGS_URI")
	viper.BindEnv("wbHost", "WB_HOST")
	viper.BindEnv("orgID", "ORG_ID")
	viper.BindEnv("schemaURI", "SCHEMA_URL")
	viper.BindEnv("testsURI", "TESTS_URI")

	viper.BindEnv("schemaFile", "SCHEMA_FILE")
	viper.BindEnv("tokenFile", "TOKEN_FILE")
	viper.BindEnv("orgFile", "ORG_FILE")
	viper.BindEnv("apiTimeout", "API_TIMEOUT")
	viper.BindEnv("versionLocation", "VERSION_LOCATION")
	viper.BindEnv("biomeDNSZone", "BIOME_DNS_ZONE")
	viper.BindEnv("statusURI", "STATUS_URI")
	viper.BindEnv("genesisCredentials", "GENESIS_CREDENTIALS")
	viper.BindEnv("banner", "GENESIS_BANNER")

	viper.BindEnv("stopTestURI", "STOP_TEST_URI")
	viper.BindEnv("stopDefURI", "STOP_DEF_URI")
	viper.BindEnv("runTestURI", "RUN_TEST_URI")
	viper.BindEnv("createTestURI", "CREATE_DEF_URI")
}

func setViperDefaults() {
	isDev := false
	for i := range os.Args {
		if os.Args[i] == "--dev" {
			isDev = true
			break
		}
	}

	if isDev {
		viper.SetDefault("authEndpoint", "auth.infra.whiteblock.io")
		viper.SetDefault("authPath", "/auth/realms/wb/protocol/openid-connect/auth")
		viper.SetDefault("tokenPath", "/auth/realms/wb/protocol/openid-connect/token")
		viper.SetDefault("multipathUploadURI", "/api/v1/files/organizations/%s/definitions")
		viper.SetDefault("wbHost", "www.infra.whiteblock.io")
		viper.SetDefault("schemaURL", "https://assets.whiteblock.io/schema/schema.json")
		viper.SetDefault("schemaFile", ".dev-test-definition-format-schema")
		viper.SetDefault("tokenFile", ".dev-auth-token")
		viper.SetDefault("orgFile", ".dev-org-name")
		viper.SetDefault("biomeDNSZone", "infra.whiteblock.io")
		viper.SetDefault("genesisBanner", "")
	} else {
		viper.SetDefault("authEndpoint", "auth.genesis.whiteblock.io")
		viper.SetDefault("authPath", "/auth/realms/wb/protocol/openid-connect/auth")
		viper.SetDefault("tokenPath", "/auth/realms/wb/protocol/openid-connect/token")
		viper.SetDefault("multipathUploadURI", "/api/v1/files/organizations/%s/definitions")
		viper.SetDefault("wbHost", "genesis.whiteblock.io")
		viper.SetDefault("schemaURL", "https://assets.whiteblock.io/schema/schema.json")
		viper.SetDefault("schemaFile", ".test-definition-format-schema")
		viper.SetDefault("tokenFile", ".auth-token")
		viper.SetDefault("orgFile", ".org-name")
		viper.SetDefault("biomeDNSZone", "biomes.whiteblock.io")
		viper.SetDefault("genesisBanner", "")
	}

	viper.SetDefault("statusURI", "/api/v1/testexecution/status/%s")
	viper.SetDefault("testsURI", "/api/v1/testexecution/organizations/%s/tests")
	viper.SetDefault("runTestURI", "/api/v1/testexecution/run/%s/%s") //org def
	viper.SetDefault("createTestURI", "/api/v1/testexecution/run/%s") //org
	viper.SetDefault("redirectURL", "localhost:56666")
	viper.SetDefault("authTimeout", 120*time.Second)
	viper.SetDefault("verbosity", "INFO")
	viper.SetDefault("getOrgURI", "/api/v1/registrar/organization/%s")
	viper.SetDefault("logsURI", "/api/v1/logs/data")
	viper.SetDefault("orgID", "")
	viper.SetDefault("apiTimeout", 5*time.Second)
	viper.SetDefault("versionLocation", "https://assets.whiteblock.io/cli/latest")
	viper.SetDefault("cliURL", "https://assets.whiteblock.io/cli/bin/genesis/%s/%s/genesis")
	viper.SetDefault("stopTestURI", "/api/v1/testexecution/stop/test/%s")
	viper.SetDefault("stopDefURI", "/api/v1/testexecution/stop/definition/%s")
}

func init() {
	setViperDefaults()
	setViperEnvBindings()
}

var (
	conf Config
	once = &sync.Once{}
)

// NewConfig creates a new config object from the global config
func NewConfig() Config {
	once.Do(func() {
		_ = viper.ReadInConfig()
		err := viper.Unmarshal(&conf)
		if err != nil {
			panic(err)
		}
		conf.setupLogrus()
		conf.Dir = configdir.New("whiteblock", "genesis-cli")
		conf.UserDir = conf.Dir.QueryFolders(configdir.Global)[0]
		conf.UserDir.MkdirAll()
	})
	return conf
}
