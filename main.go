package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

const webContent = "Hello World!"

type Database struct {
	Name string `gorm:"column:Database"`
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/services/", serviceDiscoveryHandler)
	log.Println("start service")
	log.Fatal(http.ListenAndServe(":19008", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, webContent)
}

func serviceDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	c, err := ReadFromYaml("./config.yaml")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	ret := map[string]interface{}{
		"in-namespace":    nil,
		"cross-namespace": nil,
	}

	// Test in-namespace discovery
	{
		url := fmt.Sprintf("http://%s/", c.InNamespace.Service)
		var success bool
		var msg string
		resp, err := http.Get(url)
		if err != nil {
			success = false
			msg = err.Error()
		} else {
			if resp.StatusCode == http.StatusOK {
				success = true

				b, err2 := ioutil.ReadAll(resp.Body)
				if err2 != nil {
					msg = err2.Error()
				} else {
					msg = string(b)
				}
				_ = resp.Body.Close()

			} else {
				success = false
				msg = string(resp.Status)
			}
		}
		ret["in-namespace"] = map[string]interface{}{
			"service": c.InNamespace.Service,
			"success": success,
			"msg":     msg,
		}
	}

	// Test cross namespace discovery
	{
		var success bool
		var msg interface{}
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True",
			c.CrossNamespace.User,
			c.CrossNamespace.Password,
			c.CrossNamespace.Service,
			c.CrossNamespace.DbName,
		)

		db, err := gorm.Open("mysql", dsn)
		if err != nil {
			success = false
			msg = err.Error()
		} else {
			var dbs []Database
			db.Raw("SHOW DATABASES").Scan(&dbs)
			err2 := db.Error
			if err2 != nil {
				success = false
				msg = err2.Error()
			} else {
				success = true
				msg = map[string]interface{}{
					"show_databases": dbs,
				}
			}
			_ = db.Close()
		}

		ret["cross-namespace"] = map[string]interface{}{
			"service": c.CrossNamespace.Service,
			"success": success,
			"msg":     msg,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(ret)
	_, _ = w.Write(j)
}
