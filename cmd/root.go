package cmd

import (
	"context"
	"fmt"
	"os"

	"fiber-admin/internal/app"
	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/wire"
	"fiber-admin/pkg/jwt"
	"fiber-admin/pkg/mongo"
	"fiber-admin/pkg/redis"
	"fiber-admin/pkg/utils/check"
	logging "fiber-admin/pkg/zap"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	configFile  string // config file path
	port        string // port to listen on (default is 3000)
	host        string // host to listen on (default is localhost)
	logLevel    string // log level (default is info)
	tls         bool   // enable tls
	tlsCertFile string // tls cert file path
	tlsKeyFile  string // tls key file path (default is "")
	fiberApp    *app.App
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fiber-admin",
	Short: "Fiber Admin Server",
	Long:  `Fiber Admin Server is a web server that provides APIs for managing users, notices, and documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := context.Background()
		// Start the app
		fiberApp, err = wire.InitializeApp(ctx)

		if err != nil {
			panic(err)
		}

		// Start the app
		if fiberApp.Config.BaseConfig.EnableTls {
			fiberApp.Logger.Info(
				"Starting server with TLS enabled",
				zap.String("host", fiberApp.Config.BaseConfig.AppHost),
				zap.String("port", fiberApp.Config.BaseConfig.AppPort),
				zap.String("tls_cert_file", fiberApp.Config.BaseConfig.TlsCertFile),
				zap.String("tls_key_file", fiberApp.Config.BaseConfig.TlsKeyFile),
			)
			if err := fiberApp.App.ListenTLS(
				fmt.Sprintf("%s:%s", fiberApp.Config.BaseConfig.AppHost, fiberApp.Config.BaseConfig.AppPort),
				fiberApp.Config.BaseConfig.TlsCertFile,
				fiberApp.Config.BaseConfig.TlsKeyFile,
			); err != nil {
				panic(err)
			}
		} else {
			fiberApp.Logger.Info(
				"Starting server",
				zap.String("host", fiberApp.Config.BaseConfig.AppHost),
				zap.String("port", fiberApp.Config.BaseConfig.AppPort),
			)
			if err := fiberApp.App.Listen(
				fmt.Sprintf("%s:%s", fiberApp.Config.BaseConfig.AppHost, fiberApp.Config.BaseConfig.AppPort),
			); err != nil {
				panic(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"",
		"config file",
	)
	rootCmd.PersistentFlags().StringVarP(
		&port,
		"port",
		"p",
		"",
		"port to listen on (default is 8080)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&host,
		"host",
		"H",
		"",
		"host to listen on (default is localhost)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&logLevel,
		"log-level",
		"l",
		"info",
		"log level (default is info)",
	)
	rootCmd.PersistentFlags().BoolVar(
		&tls,
		"tls",
		false,
		"enable tls (default is false)",
	)
	rootCmd.PersistentFlags().StringVar(
		&tlsCertFile,
		"tls-cert-file",
		"",
		"tls cert file path (default is \"\")",
	)
	rootCmd.PersistentFlags().StringVar(
		&tlsKeyFile,
		"tls-key-file",
		"",
		"tls key file path (default is \"\")",
	)
}

// initConfig creates a new config object and initializes it with the values from the config file and
// priority: flags > config file > default values
func initConfig() {
	cfg := config.New()
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		var home string
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fiber-admin")
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
	err = viper.Unmarshal(cfg)
	cobra.CheckErr(err)
	if port != "" {
		if !check.IsValidAppPort(port) {
			cobra.CheckErr(fmt.Errorf("invalid port: %s", port))
		}
		cfg.BaseConfig.AppPort = port
	}
	if host != "" {
		if !check.IsValidAppHost(host) {
			cobra.CheckErr(fmt.Errorf("invalid host: %s", host))
		}
		cfg.BaseConfig.AppHost = host
	}
	if logLevel != "" {
		if !check.IsValidLogLevel(logLevel) {
			cobra.CheckErr(fmt.Errorf("invalid log level: %s", logLevel))
		}
		cfg.ZapConfig.Level = logLevel
	}
	if tls {
		cfg.BaseConfig.EnableTls = true
	}
	if tlsCertFile != "" {
		// Check if the tls cert file exists
		if _, err := os.Stat(tlsCertFile); os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("tls cert file does not exist: %s", tlsCertFile))
		}
		cfg.BaseConfig.TlsCertFile = tlsCertFile
	}
	if tlsKeyFile != "" {
		// Check if the tls key file exists
		if _, err := os.Stat(tlsKeyFile); os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("tls key file does not exist: %s", tlsKeyFile))
		}
		cfg.BaseConfig.TlsKeyFile = tlsKeyFile
	}
	config.Set(cfg)
	viper.WatchConfig()
	viper.OnConfigChange(
		func(in fsnotify.Event) {
			fmt.Printf("Reloading config file: %s ...\n", in.Name)

			if err := viper.Unmarshal(cfg); err != nil {
				fmt.Printf("error reading config: %s\n", err)
			}
			config.Set(cfg)

			// Reload singletons
			if err := redis.Update(context.Background(), cfg.CacheConfig.RedisConfig.GetRedisOptions()); err != nil {
				fmt.Printf("error setting redis: %s\n", err)
			}
			if err := logging.Update(cfg.ZapConfig.GetZapConfig()); err != nil {
				fmt.Printf("error setting zap: %s\n", err)
			}
			if err := mongo.Update(
				context.Background(), cfg.MongoConfig.GetQmgoConfig(), cfg.MongoConfig.PingTimeoutS,
				cfg.MongoConfig.Database,
			); err != nil {
				fmt.Printf("error setting mongo: %s\n", err)
			}
			if err := jwt.Update(
				nil, cfg.JWTConfig.TokenDuration, cfg.JWTConfig.RefreshDuration, cfg.JWTConfig.RefreshBuffer,
			); err != nil {
				fmt.Printf("error setting jwt: %s\n", err)
			}

			fmt.Printf("Config file: %s reloaded\n", in.Name)
		},
	)
}
