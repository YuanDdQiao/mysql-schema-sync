package main

import (
	"flag"
	"fmt"
	"log"
	"mysql-schema-sync/internal"
	"os"
	"strings"
)

var configPath = flag.String("conf", "./config.json", "json config file path")
var sync = flag.Bool("sync", false, "sync shcema change to dest db")
var drop = flag.Bool("drop", false, "drop fields,index,foreign key")

var source = flag.String("from", "", "mysql dsn source,eg: username[:password]@(ip:port)/dbname\n\twhen it is not empty ignore [-conf] param")
var dest = flag.String("to", "", "mysql dsn dest,eg username[:password]@(ip:port)/dbname")
var tables = flag.String("tables", "", "table names to check\n\teg : product_base,order_*")
var mailTo = flag.String("mail_to", "", "overwrite config's email.to")

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate)
	df := flag.Usage
	flag.Usage = func() {
		df()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "mysql schema sync tools "+internal.Version)
		fmt.Fprintln(os.Stderr, internal.AppURL+"\n")
	}
}

var cfg *internal.Config

func main() {
	flag.Parse()
	if *source == "" {
		cfg = internal.LoadConfig(*configPath)
	} else {
		cfg = new(internal.Config)
		cfg.SourceDSN = *source
		cfg.DestDSN = *dest
	}
	cfg.Sync = *sync
	cfg.Drop = *drop

	if *mailTo != "" && cfg.Email != nil {
		cfg.Email.To = *mailTo
	}

	if cfg.Tables == nil {
		cfg.Tables = []string{}
	}
	if *tables != "" {
		_ts := strings.Split(*tables, ",")
		for _, _name := range _ts {
			_name = strings.TrimSpace(_name)
			if _name != "" {
				cfg.Tables = append(cfg.Tables, _name)
			}
		}
	}
	defer (func() {
		if err := recover(); err != nil {
			log.Println(err)
			cfg.SendMailFail(fmt.Sprintf("%s", err))
			log.Fatalln("exit")
		}
	})()

	cfg.Check()
	internal.CheckSchemaDiff(cfg)
}
