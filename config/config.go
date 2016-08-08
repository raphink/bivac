package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/jessevdk/go-flags"
)

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version             bool     `short:"V" long:"version" description:"Display version."`
	Loglevel            string   `short:"l" long:"loglevel" description:"Set loglevel ('debug', 'info', 'warn', 'error', 'fatal', 'panic')." env:"CONPLICITY_LOG_LEVEL" default:"info"`
	VolumesBlacklist    []string `short:"b" long:"blacklist" description:"Volumes to blacklist in backups." env:"CONPLICITY_VOLUMES_BLACKLIST" env-delim:","`
	Manpage             bool     `short:"m" long:"manpage" description:"Output manpage."`
	NoVerify            bool     `long:"no-verify" description:"Do not verify backup." env:"CONPLICITY_NO_VERIFY"`
	JSON                bool     `short:"j" long:"json" description:"Log as JSON (to stderr)." env:"CONPLICITY_JSON_OUTPUT"`
	Engine              string   `short:"E" long:"engine" description:"Backup engine to use." env:"CONPLICITY_ENGINE" default:"duplicity" ini-name:""`
	HostnameFromRancher bool     `short:"H" long:"hostname-from-rancher" description:"Retrieve hostname from Rancher metadata." env:"CONPLICITY_HOSTNAME_FROM_RANCHER"`

	Duplicity struct {
		Image           string `long:"duplicity-image" description:"The duplicity docker image." env:"DUPLICITY_DOCKER_IMAGE" default:"camptocamp/duplicity:latest"`
		TargetURL       string `long:"duplicity-url" description:"The duplicity target URL to push to." env:"DUPLICITY_TARGET_URL"`
		FullIfOlderThan string `long:"full-if-older-than" description:"The number of days after which a full backup must be performed." env:"CONPLICITY_FULL_IF_OLDER_THAN" default:"15D"`
		RemoveOlderThan string `long:"remove-older-than" description:"The number days after which backups must be removed." env:"CONPLICITY_REMOVE_OLDER_THAN" default:"30D"`
	} `group:"Duplicity Options"`

	RClone struct {
		Image     string `long:"rclone-image" description:"The rclone docker image." env:"RCLONE_DOCKER_IMAGE" default:"camptocamp/rclone:latest"`
		TargetURL string `long:"rclone-url" description:"The RClone target URL to push to." env:"RCLONE_TARGET_URL"`
	} `group:"RClone Options"`

	Metrics struct {
		PushgatewayURL string `short:"g" long:"gateway-url" description:"The prometheus push gateway URL to use." env:"PUSHGATEWAY_URL"`
	} `group:"Metrics Options"`

	AWS struct {
		AccessKeyID     string `long:"aws-access-key-id" description:"The AWS access key ID." env:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `long:"aws-secret-key-id" description:"The AWS secret access key." env:"AWS_SECRET_ACCESS_KEY"`
	} `group:"AWS Options"`

	Swift struct {
		Username   string `long:"swift-username" description:"The Swift user name." env:"SWIFT_USERNAME"`
		Password   string `long:"swift-password" description:"The Swift password." env:"SWIFT_PASSWORD"`
		AuthURL    string `long:"swift-auth_url" description:"The Swift auth URL." env:"SWIFT_AUTHURL"`
		TenantName string `long:"swift-tenant-name" description:"The Swift tenant name." env:"SWIFT_TENANTNAME"`
		RegionName string `long:"swift-region-name" description:"The Swift region name." env:"SWIFT_REGIONNAME"`
	} `group:"Swift Options"`

	Docker struct {
		Endpoint string `short:"e" long:"docker-endpoint" description:"The Docker endpoint." env:"DOCKER_ENDPOINT" default:"unix:///var/run/docker.sock"`
	} `group:"Docker Options"`
}

// LoadConfig loads the config from flags & environment
func LoadConfig(version string) *Config {
	var c Config
	parser := flags.NewParser(&c, flags.Default)

	conffile := fmt.Sprintf("%s/.conplicity.conf", os.Getenv("HOME"))
	log.Println(conffile)

	parser.Parse()
	iniParser := flags.NewIniParser(parser)
	if err := iniParser.ParseFile(conffile); err != nil {
		os.Exit(1)
	}

	if err := iniParser.WriteFile(conffile, flags.IniNone); err != nil {
		log.Fatal(err)
	}

	if c.Version {
		fmt.Printf("Conplicity v%v\n", version)
		os.Exit(0)
	}

	if c.Manpage {
		var buf bytes.Buffer
		parser.WriteManPage(&buf)
		fmt.Printf(buf.String())
		os.Exit(0)
	}

	sort.Strings(c.VolumesBlacklist)
	return &c
}
