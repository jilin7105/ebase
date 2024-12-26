package ebase

//
//import (
//	esv7 "github.com/elastic/go-elasticsearch/v7"
//	esv8 "github.com/elastic/go-elasticsearch/v8"
//	"github.com/jilin7105/ebase/config"
//	"github.com/jilin7105/ebase/logger"
//)
//
//type EsEbase struct {
//	esv7    *esv7.Client
//	esv8    *esv8.Client
//	Version string
//}
//
//func newEs(config config.EsConfig) (EsEbase, error) {
//	var es = EsEbase{
//		Version: config.Version,
//	}
//	cfgv8 := esv8.Config{}
//	cfgv7 := esv7.Config{}
//
//	if config.Type == "host" {
//
//		cfgv7 = esv7.Config{
//			Addresses: config.Hosts,
//		}
//		cfgv8 = esv8.Config{
//			Addresses: config.Hosts,
//		}
//		if config.User != "" {
//			cfgv8.Username = config.User
//			cfgv7.Username = config.User
//
//		}
//
//		if config.Pass != "" {
//			cfgv8.Password = config.Pass
//			cfgv7.Password = config.Pass
//		}
//	}
//
//	if config.Type == "cloud" {
//		cfgv8.CloudID = config.CloudId
//		cfgv7.CloudID = config.CloudId
//		cfgv8.APIKey = config.ApiKey
//		cfgv7.APIKey = config.ApiKey
//	}
//
//	var err error
//	if config.Version == "v8" {
//		es.esv8, err = esv8.NewClient(cfgv8)
//
//	}
//
//	if config.Version == "v7" {
//		es.esv7, err = esv7.NewClient(cfgv7)
//	}
//
//	return es, err
//
//}
//
//func (e *Eb) inites() {
//
//	for _, config := range e.Config.Es {
//		esclient, err := newEs(config)
//		if err != nil {
//			logger.Info("es  初始化失败 %s %s", config.Name, err.Error())
//		}
//		e.ES[config.Name] = esclient
//	}
//
//}
