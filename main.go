package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

type ScriptName struct {
	Name     string `yaml:"name"`
	Dir      string `yaml:"dir"`
	Executor string `yaml:"executor"`
}

type Yml struct {
	ScriptList []ScriptName `yaml:"scripts"`
}

var (
	scriptMetric = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "script_success",
			Help: "Return script exit code (0 is success, 1 is failed)",
		},
		[]string{"name"})

	configDir = os.Args[1]
)

func main() {
	go metricScrape(configDir)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2307", nil)
}

func metricScrape(configDir string) {
	yml := Yml{}
	yamlFile, err := ioutil.ReadFile(configDir)
	handlError("Read yml file error ->", err)

	e := yaml.Unmarshal(yamlFile, &yml)
	handlError("Unmarshal error ->", e)
	for {
		for _, name := range yml.ScriptList {
			metricHandler(name.Name, name.Dir, name.Executor)
		}
		time.Sleep(5 * time.Second)
	}

}
func metricHandler(scriptName, scriptDir, scriptExecutor string) {
	_, err := exec.Command(scriptExecutor, scriptDir).Output()
	if err != nil {
		log.Printf("Script %v return %v", scriptName, err.(*exec.ExitError))
		scriptMetric.WithLabelValues(scriptName).Set(1)
	} else {
		log.Printf("Script %v success return exit code 0", scriptName)
		scriptMetric.WithLabelValues(scriptName).Set(0)

	}
}

func handlError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
