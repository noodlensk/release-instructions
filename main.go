package main

import (
	"net/http"
	"os"

	jira "github.com/andygrunwald/go-jira"
	"github.com/gorilla/mux"
	bitbucket "github.com/noodlensk/go-bitbucket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "development"

type Repo struct {
	Name          string
	Branch        string `json:"-"`
	ExcludeBranch string `json:"-"`
	Weight        int
}
type Configuration struct {
	HTTPListen string
	Jira       struct {
		Host     string
		User     string
		Password string
		Project  string
	}
	Bitbucket struct {
		User     string
		Password string
		Owner    string
		Repos    []Repo
	}
}

var config Configuration
var cfgFile string

var rootCmd = &cobra.Command{
	Use:           "release-sintructions",
	Short:         "Run release instructions TBD",
	Long:          `Run release instructions TBD`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		tp := jira.BasicAuthTransport{
			Username: config.Jira.User,
			Password: config.Jira.Password,
		}
		jiraClient, _ := jira.NewClient(tp.Client(), config.Jira.Host)
		bitbucketClient := bitbucket.NewBasicAuth(config.Bitbucket.User, config.Bitbucket.Password)
		app := NewApp(jiraClient, bitbucketClient, config.Bitbucket.Repos)
		r := mux.NewRouter()
		r.Handle("/api/instructions", appHTTPHandler{app, releaseInstrucionsHandler}).Methods("GET")
		r.PathPrefix("/").Handler(appHTTPHandler{app, indexHandler}).Methods("GET")
		return http.ListenAndServe(config.HTTPListen, r)
	},
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file
func initConfig() {
	log.Infof("Running version %s", version)
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("config file required")
	}
	log.Infof("using config file: %s", viper.ConfigFileUsed())
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
}
