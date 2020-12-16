package config

import (
	"flag"
	"net/http"
	"os"

	"log"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"golang.org/x/sync/syncmap"
)

var mode string
var confFileLoc string
var env string
var cfg ini.File
var verboses bool
var exporterPath string
var configs syncmap.Map

func init() {
	configs = syncmap.Map{}
	flag.StringVar(&env, "e", "prod", "Environment Parameters (production,staging,local,etc) ")
	flag.StringVar(&confFileLoc, "c", "config.conf", "Config File Path -- if using remote mode, see example in http://10.222.12.46:7000/alfamikro-service-config/prod")
	flag.StringVar(&mode, "m", "local", "Config Mode -- local|remote, with local mode as default")

	flag.Parse()

	err := initConfig(mode)
	if FancyHandleError(err) {
		log.Fatal(err.Error())
	}

	/*//TODO -- ini buat remote
	loadConfig(env)

	ticker := time.NewTicker(time.Minute * 5)

	go func() {
		for _ = range ticker.C {
			loadConfig(env)
		}
	}()

	//TODO -- ini buat lokal
	config, err := ini.Load(confFileLoc)
	//log.Println(config)
	if err != nil {
		log.Println("Can't open config file /n creating new Config files")
		f, err := os.Create(confFileLoc)
		f.WriteString("")
		f.Close()
		config, err := ini.Load(confFileLoc)
		if err != nil {
			panic(err)
		}
		cfg = *config
	} else {
		cfg = *config
	}*/

}
func GetCurrentEnvValueFromConfig(key string) string {
	return GetValueFromConfig(env, key)
}
func GetValueFromConfig(environment string, key string) string {

	switch mode {
	case "remote":
		//TODO -- ini buat remote
		if val, ok := configs.Load(key); ok {
			if stringVal, ok := val.(string); ok {
				return stringVal
			}
		}
		break
	default:
		//TODO -- ini buat lokal
		if cfg.Section(environment).HasKey(key) {
			return cfg.Section(environment).Key(key).String()
		} else {
			return cfg.Section("").Key(key).String()
		}
	}

	return ""
}

func loadConfig(env string) error {
	if env == "" {
		env = "prod"
	}
	var netClient = &http.Client{}

	resp, err := netClient.Get(confFileLoc + env)
	if FancyHandleError(err) {
		return err
	}
	defer resp.Body.Close()

	body, err := jason.NewObjectFromReader(resp.Body)
	if FancyHandleError(err) {
		return err
	}

	if propSources, err := body.GetObjectArray("propertySources"); !FancyHandleError(err) {
		for _, obj := range propSources {
			if source, err := obj.GetObject("source"); !FancyHandleError(err) {
				for key, val := range source.Map() {
					valString, _ := val.String()
					configs.Store(key, valString)
				}
				break
			}
		}
	} else {
		return err
	}

	log.Println("global config loaded. timestamp ", time.Now().Local().String())

	return nil
}

func initConfig(mode string) error {
	switch mode {
	case "remote":
		err := loadConfig(env)
		if FancyHandleError(err) {
			return err
		}

		ticker := time.NewTicker(time.Minute * 5)

		go func() {
			for _ = range ticker.C {
				loadConfig(env)
			}
		}()
		break
	default:
		config, err := ini.Load(confFileLoc)
		//log.Println(config)
		if err != nil {
			log.Println("Can't open config file /n creating new Config files")
			f, err := os.Create(confFileLoc)
			f.WriteString("")
			f.Close()
			config, err := ini.Load(confFileLoc)
			if err != nil {
				panic(err)
			}
			cfg = *config
		} else {
			cfg = *config
		}
		break
	}
	return nil
}
