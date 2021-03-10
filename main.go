package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy/cacher"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

var cfgFile = pflag.StringP("config", "c", "config.yaml", "config file")

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetEnvPrefix("GU")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	spew.Dump(viper.AllSettings())
}

func main() {
	g := goproxy.New()
	g.GoBinEnv = append(
		os.Environ(),
	)
	g.Cacher = &cacher.Disk{
		Root: viper.GetString("cache.disk.root"),
	}
	g.ProxiedSUMDBs = viper.GetStringSlice("goproxy.proxied_sumdbs")

	log.Println("Server Listen on: ", viper.GetString("http.listen"))
	http.ListenAndServe(viper.GetString("http.listen"), g)
}
