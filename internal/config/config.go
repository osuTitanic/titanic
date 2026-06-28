package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/osuTitanic/titanic-go/internal/storage"
)

// Config holds all application configurations
type Config struct {
	// Database configuration
	PostgresHost             string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort             int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser             string `env:"POSTGRES_USER" envDefault:"bancho"`
	PostgresDatabase         string `env:"POSTGRES_DATABASE"`
	PostgresPassword         string `env:"POSTGRES_PASSWORD,required"`
	PostgresPoolEnabled      bool   `env:"POSTGRES_POOL_ENABLED" envDefault:"true"`
	PostgresPoolSize         int    `env:"POSTGRES_POOL_SIZE" envDefault:"10"`
	PostgresPoolSizeOverflow int    `env:"POSTGRES_POOL_SIZE_OVERFLOW" envDefault:"30"`
	PostgresPoolPrePing      bool   `env:"POSTGRES_POOL_PRE_PING" envDefault:"true"`
	PostgresPoolRecycle      int    `env:"POSTGRES_POOL_RECYCLE" envDefault:"900"`
	PostgresPoolTimeout      int    `env:"POSTGRES_POOL_TIMEOUT" envDefault:"15"`

	// Redis configuration
	RedisHost string  `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort int     `env:"REDIS_PORT" envDefault:"6379"`
	RedisPass *string `env:"REDIS_PASS"`

	// Path to store application data locally
	DataPath string `env:"DATA_PATH" envDefault:".data"`

	// S3 configuration (optional)
	S3Enabled   bool   `env:"S3_ENABLED" envDefault:"false"`
	S3BaseUrl   string `env:"S3_BASEURL"`
	S3AccessKey string `env:"S3_ACCESS_KEY"`
	S3SecretKey string `env:"S3_SECRET_KEY"`
	S3Bucket    string `env:"S3_BUCKET" envDefault:"osutitanic"`
	S3Region    string `env:"S3_REGION" envDefault:""`

	// Cloudflare cache purge configuration (optional)
	CloudflarePurgeEnabled bool        `env:"CLOUDFLARE_PURGE_ENABLED" envDefault:"false"`
	CloudflareZoneId       string      `env:"CLOUDFLARE_ZONE_ID"`
	CloudflareApiToken     string      `env:"CLOUDFLARE_API_TOKEN"`
	CloudflarePurgeOszUrls StringSlice `env:"CLOUDFLARE_PURGE_OSZ_URLS" envDefault:""`

	// Menu icon configuration (optional)
	MenuIconImage string `env:"MENUICON_IMAGE"`
	MenuIconUrl   string `env:"MENUICON_URL"`

	// Seasonal backgrounds
	SeasonalBackgrounds StringSlice `env:"SEASONAL_BACKGROUNDS"`

	// Discord webhook configuration (optional)
	OfficerWebhookUrl        string `env:"OFFICER_WEBHOOK_URL"`
	AnnounceEventsWebhookUrl string `env:"ANNOUNCE_EVENTS_WEBHOOK_URL"`
	ForumEventsWebhookUrl    string `env:"FORUM_EVENTS_WEBHOOK_URL"`
	BeatmapEventsWebhookUrl  string `env:"BEATMAP_EVENTS_WEBHOOK_URL"`

	// Image proxy baseurl for bbcode (optional)
	ImageProxyBaseUrl string `env:"IMAGE_PROXY_BASEURL"`

	// Email configuration (optional)
	EmailProvider string `env:"EMAIL_PROVIDER" envDefault:"noop"`
	EmailSender   string `env:"EMAIL_SENDER" envDefault:"support@titanic.sh"`

	// SMTP configuration
	SmtpHost          string `env:"SMTP_HOST"`
	SmtpPort          int    `env:"SMTP_PORT" envDefault:"587"`
	SmtpUser          string `env:"SMTP_USER"`
	SmtpPassword      string `env:"SMTP_PASSWORD"`
	SmtpTls           bool   `env:"SMTP_TLS" envDefault:"true"`
	SmtpSkipTlsVerify bool   `env:"SMTP_SKIP_TLS_VERIFY" envDefault:"false"`

	// Score server configuration
	WebHost                    string `env:"WEB_HOST" envDefault:"localhost"`
	WebPort                    int    `env:"WEB_PORT" envDefault:"80"`
	WebWorkers                 int    `env:"WEB_WORKERS" envDefault:"5"`
	ScoreResponseLimit         int    `env:"SCORE_RESPONSE_LIMIT" envDefault:"50"`
	BeatmapFavoritesLimit      int    `env:"BEATMAP_FAVOURITES_LIMIT" envDefault:"100"`
	AllowRelax                 bool   `env:"ALLOW_RELAX" envDefault:"false"`
	AllowUnauthenticatedDirect bool   `env:"ALLOW_UNAUTHENTICATED_DIRECT" envDefault:"true"`
	ApprovedMapRewards         bool   `env:"APPROVED_MAP_REWARDS" envDefault:"false"`
	BeatmapSubmissionEnabled   bool   `env:"BEATMAP_SUBMISSION_ENABLED" envDefault:"false"`
	FrozenRankUpdates          bool   `env:"FROZEN_RANK_UPDATES" envDefault:"false"`
	FrozenPPv1Updates          bool   `env:"FROZEN_PPV1_UPDATES" envDefault:"false"`

	// Bancho configuration
	BanchoTcpPorts            IntSlice    `env:"BANCHO_TCP_PORTS" envDefault:"13380,13381,13382,13383"`
	BanchoHttpPort            int         `env:"BANCHO_HTTP_PORT" envDefault:"5000"`
	BanchoWsPort              int         `env:"BANCHO_WS_PORT" envDefault:"5001"`
	BanchoIrcPort             int         `env:"BANCHO_IRC_PORT" envDefault:"6667"`
	BanchoIrcPortSsl          int         `env:"BANCHO_IRC_PORT_SSL" envDefault:"6697"`
	BanchoWorkers             int         `env:"BANCHO_WORKERS" envDefault:"16"`
	BanchoRestartTime         int         `env:"BANCHO_RESTART_TIME" envDefault:"10"`
	IrcEnabled                bool        `env:"IRC_ENABLED" envDefault:"true"`
	OsuIrcEnabled             bool        `env:"OSU_IRC_ENABLED" envDefault:"true"`
	BanchoSslKeyfile          string      `env:"BANCHO_SSL_KEYFILE"`
	BanchoSslCertfile         string      `env:"BANCHO_SSL_CERTFILE"`
	BanchoSslVerifyFile       string      `env:"BANCHO_SSL_VERIFY_FILE"`
	BanchoMaintenance         bool        `env:"BANCHO_MAINTENANCE" envDefault:"false"`
	AllowMultiaccounting      bool        `env:"ALLOW_MULTIACCOUNTING" envDefault:"false"`
	AutojoinChannels          StringSlice `env:"AUTOJOIN_CHANNELS" envDefault:"#osu,#announce"`
	BanchoIp                  string      `env:"BANCHO_IP"`
	DisableClientVerification bool        `env:"DISABLE_CLIENT_VERIFICATION" envDefault:"true"`
	BanchoClientCutoff        int         `env:"BANCHO_CLIENT_CUTOFF"`
	MultiplayerMaxSlots       int         `env:"MULTIPLAYER_MAX_SLOTS" envDefault:"8"`

	// Website configuration
	FrontendHost          string      `env:"FRONTEND_HOST" envDefault:"localhost"`
	FrontendPort          int         `env:"FRONTEND_PORT" envDefault:"8080"`
	FrontendWorkers       int         `env:"FRONTEND_WORKERS" envDefault:"4"`
	FrontendSecretKey     string      `env:"FRONTEND_SECRET_KEY" envDefault:"somethingrandom"`
	FrontendTokenExpiry   int         `env:"FRONTEND_TOKEN_EXPIRY" envDefault:"3600"`
	FrontendRefreshExpiry int         `env:"FRONTEND_REFRESH_EXPIRY" envDefault:"2592000"`
	EnableSsl             bool        `env:"ENABLE_SSL" envDefault:"false"`
	AllowInsecureCookies  *bool       `env:"ALLOW_INSECURE_COOKIES"`
	RecaptchaSecretKey    string      `env:"RECAPTCHA_SECRET_KEY"`
	RecaptchaSiteKey      string      `env:"RECAPTCHA_SITE_KEY"`
	SuperFriendlyUsers    IntSlice    `env:"SUPER_FRIENDLY_USERS" envDefault:"[1]"`
	BeginningEndedAt      DynamicTime `env:"BEGINNING_ENDED_AT" envDefault:"2023-12-31T06:00:00Z"`

	// Wiki configuration
	WikiRepositoryOwner  string `env:"WIKI_REPOSITORY_OWNER" envDefault:"osuTitanic"`
	WikiRepositoryName   string `env:"WIKI_REPOSITORY_NAME" envDefault:"wiki"`
	WikiRepositoryBranch string `env:"WIKI_REPOSITORY_BRANCH" envDefault:"main"`
	WikiRepositoryPath   string `env:"WIKI_REPOSITORY_PATH" envDefault:"wiki"`
	WikiDefaultLanguage  string `env:"WIKI_DEFAULT_LANGUAGE" envDefault:"en"`
	RemoveScoresOnRanked bool   `env:"REMOVE_SCORES_ON_RANKED" envDefault:"true"`

	// API configuration
	ApiHost                   string `env:"API_HOST" envDefault:"localhost"`
	ApiPort                   int    `env:"API_PORT" envDefault:"8000"`
	ApiWorkers                int    `env:"API_WORKERS" envDefault:"4"`
	ApiRateLimitEnabled       bool   `env:"API_RATELIMIT_ENABLED" envDefault:"true"`
	ApiRateLimitWindow        int    `env:"API_RATELIMIT_WINDOW" envDefault:"60"`
	ApiRateLimitRegular       int    `env:"API_RATELIMIT_REGULAR" envDefault:"400"`
	ApiRateLimitAuthenticated int    `env:"API_RATELIMIT_AUTHENTICATED" envDefault:"800"`

	// Ko-Fi token for donation callbacks
	KofiVerificationToken string `env:"KOFI_VERIFICATION_TOKEN"`

	// Discord bot configuration (optional)
	EnableDiscordBot   bool   `env:"ENABLE_DISCORD_BOT" envDefault:"false"`
	DiscordBotPrefix   string `env:"DISCORD_BOT_PREFIX" envDefault:"!"`
	DiscordBotToken    string `env:"DISCORD_BOT_TOKEN"`
	DiscordStaffRoleId int64  `env:"DISCORD_STAFF_ROLE_ID"`
	DiscordBatRoleId   int64  `env:"DISCORD_BAT_ROLE_ID"`

	// osu! API credentials (optional)
	OsuClientId     string `env:"OSU_CLIENT_ID"`
	OsuClientSecret string `env:"OSU_CLIENT_SECRET"`

	// Chat webhook configuration (optional)
	ChatWebhookUrl      string      `env:"CHAT_WEBHOOK_URL"`
	ChatChannelId       int64       `env:"CHAT_CHANNEL_ID"`
	ChatWebhookChannels StringSlice `env:"CHAT_WEBHOOK_CHANNELS" envDefault:"#osu"`

	// Release stream updates (optional)
	ReleaseUpdatesEnabled      bool   `env:"RELEASE_UPDATES_ENABLED" envDefault:"false"`
	ReleaseUpdateLocation      string `env:"RELEASE_UPDATE_LOCATION" envDefault:"release"`
	ReleaseUpdateNotifyWebhook string `env:"RELEASE_UPDATE_NOTIFY_WEBHOOK"`

	// Debugging options
	Debug  bool `env:"DEBUG" envDefault:"false"`
	Reload bool `env:"RELOAD" envDefault:"false"`

	// Domain configuration
	DomainName string `env:"DOMAIN_NAME" envDefault:"localhost"`

	// Custom URL overrides (optional)
	OsuBaseUrlOverride      string `env:"OSU_BASEURL"`
	ApiBaseUrlOverride      string `env:"API_BASEURL"`
	StaticBaseUrlOverride   string `env:"STATIC_BASEURL"`
	EventsWebsocketOverride string `env:"EVENTS_WEBSOCKET"`
	LoungeBackendOverride   string `env:"LOUNGE_BACKEND"`

	// Image services that bypass the image proxy
	ValidImageServicesOverride StringSlice `env:"VALID_IMAGE_SERVICES"`
}

// LoadConfig loads configuration from environment variables & the .env file
func LoadConfig(envFiles ...string) (*Config, error) {
	if len(envFiles) == 0 {
		envFiles = []string{".env"}
	}
	for _, file := range envFiles {
		_ = godotenv.Load(file)
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) PostgresDSN() string {
	database := c.PostgresDatabase
	if database == "" {
		database = c.PostgresUser
	}
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, database,
	)
}

func (c *Config) S3Config() *storage.S3Config {
	baseUrl, err := url.Parse(c.S3BaseUrl)
	if err != nil {
		return nil
	}

	return &storage.S3Config{
		Endpoint:        baseUrl.Host,
		UseSSL:          baseUrl.Scheme == "https",
		Region:          c.S3Region,
		BucketName:      c.S3Bucket,
		AccessKeyID:     c.S3AccessKey,
		SecretAccessKey: c.S3SecretKey,
	}
}

func (c *Config) EmailsEnabled() bool {
	return c.EmailProvider != "" && c.EmailSender != ""
}

func (c *Config) EmailDomain() *string {
	if !strings.Contains(c.EmailSender, "@") {
		return nil
	}
	parts := strings.SplitN(c.EmailSender, "@", 2)
	return &parts[1]
}

func (c *Config) OsuBaseUrl() string {
	if c.OsuBaseUrlOverride != "" {
		return c.OsuBaseUrlOverride
	}
	return c.DefaultOsuBaseUrl()
}

func (c *Config) ApiBaseUrl() string {
	if c.ApiBaseUrlOverride != "" {
		return c.ApiBaseUrlOverride
	}
	return c.DefaultApiBaseUrl()
}

func (c *Config) StaticBaseUrl() string {
	if c.StaticBaseUrlOverride != "" {
		return c.StaticBaseUrlOverride
	}
	return c.DefaultStaticBaseUrl()
}

func (c *Config) EventsWebsocket() string {
	if c.EventsWebsocketOverride != "" {
		return c.EventsWebsocketOverride
	}
	return c.DefaultEventsWebsocket()
}

func (c *Config) LoungeBackend() string {
	if c.LoungeBackendOverride != "" {
		return c.LoungeBackendOverride
	}
	return c.DefaultLoungeBackend()
}

func (c *Config) DefaultOsuBaseUrl() string {
	scheme := "http"
	if c.EnableSsl {
		scheme = "https"
	}
	return fmt.Sprintf("%s://osu.%s", scheme, c.DomainName)
}

func (c *Config) DefaultApiBaseUrl() string {
	scheme := "http"
	if c.EnableSsl {
		scheme = "https"
	}
	return fmt.Sprintf("%s://api.%s", scheme, c.DomainName)
}

func (c *Config) DefaultStaticBaseUrl() string {
	scheme := "http"
	if c.EnableSsl {
		scheme = "https"
	}
	return fmt.Sprintf("%s://s.%s", scheme, c.DomainName)
}

func (c *Config) DefaultEventsWebsocket() string {
	scheme := "ws"
	if c.EnableSsl {
		scheme = "wss"
	}
	return fmt.Sprintf("%s://api.%s/events/ws", scheme, c.DomainName)
}

func (c *Config) DefaultLoungeBackend() string {
	scheme := "http"
	if c.EnableSsl {
		scheme = "https"
	}
	return fmt.Sprintf("%s://lounge.%s", scheme, c.DomainName)
}

var defaultValidImageServices = []string{
	"ibb.co",
	"i.ibb.co",
	"i.ppy.sh",
	"i.imgur.com",
	"cdn.discordapp.com",
	"media.discordapp.net",
}

func (c *Config) ValidImageServices() []string {
	services := make(map[string]struct{})
	for _, s := range defaultValidImageServices {
		services[s] = struct{}{}
	}
	for _, s := range c.ValidImageServicesOverride {
		services[s] = struct{}{}
	}
	services[fmt.Sprintf("i.%s", c.DomainName)] = struct{}{}
	services[fmt.Sprintf("osu.%s", c.DomainName)] = struct{}{}

	result := make([]string, 0, len(services))
	for s := range services {
		result = append(result, s)
	}
	return result
}

func (c *Config) GetAllowInsecureCookies() bool {
	if c.AllowInsecureCookies != nil {
		return *c.AllowInsecureCookies
	}
	return !c.EnableSsl || c.Debug
}

func (c *Config) SitemapEnabled() bool {
	return c.DomainName == "titanic.sh" || c.DomainName == "localhost"
}

func (c *Config) ChristmasMode() bool {
	now := time.Now()
	return now.Month() == time.December && now.Day() >= 15
}
