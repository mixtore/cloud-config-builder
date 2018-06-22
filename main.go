package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

type (
	Kubernetes struct {
		EnvVars EnvVars

		Namespace string
		Name      string
	}

	AppEngine struct {
		EnvVars EnvVars

		ScalingField AppEngineScalingField
		Scaling      bool

		ResourcesField AppEngineResourcesField
		Resources      bool

		DisableHealthCheck bool

		Name    string
		Runtime string
		Env     string
		Command string
	}

	AppEngineScalingField struct {
		MinNumInstances string
		MaxNumInstances string
		CPUUtilization  string
	}

	AppEngineResourcesField struct {
		CPUCount string
		MemoryGB string
	}

	EnvVars map[string]string
)

var (
	kubernetesConfigMap = template.Must(template.New("kubernetesConfigMap").Parse(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: {{ .Namespace }}
  name: {{ .Name }}
data:
{{- range $key, $value := .EnvVars }}
  {{ $key }}: "{{ $value }}"
{{- end }}
`))

	appEngineWeb = template.Must(template.New("appEngine").Parse(`service: {{ .Name }}
runtime: {{ .Runtime }}
env: {{ .Env }}
entrypoint: {{ .Command }}

{{ if .DisableHealthCheck -}}
health_check:
  enable_health_check: false
{{- end }}

{{ if .Scaling -}}
automatic_scaling:
  min_num_instances: {{ .ScalingField.MinNumInstances }}
  max_num_instances: {{ .ScalingField.MaxNumInstances }}
  cpu_utilization:
    target_utilization: {{ .ScalingField.CPUUtilization }}
{{- end }}

{{ if .Resources -}}
resources:
  memory_gb: {{ .ResourcesField.MemoryGB }}
  cpu: {{ .ResourcesField.CPUCount }}
{{- end }}

env_variables:
{{- range $key, $value := .EnvVars }}
  {{ $key }}: "{{ $value }}"
{{- end }}
`))
)

func main() {
	flagSettingsNamespacePtr := flag.String("namespace", "", "set namespace")
	flagSettingsNamePtr := flag.String("name", "", "set name")
	flagSettingsRuntimePtr := flag.String("runtime", "", "set runtime")
	flagSettingsEnvPtr := flag.String("env", "", "set env")
	flagSettingsCommandPtr := flag.String("command", "", "set command")
	flagSettingsDisableHealthCheckPtr := flag.Bool("disable-healthcheck", false, "disable healthcheck")

	flagSettingsScalingPtr := flag.Bool("scaling", true, "enable automatic scaling")
	flagSettingsScalingMinPtr := flag.String("scaling-min", "", "set scaling min num instances")
	flagSettingsScalingMaxPtr := flag.String("scaling-max", "", "set scaling max num instances")
	flagSettingsScalingCPUPtr := flag.String("scaling-cpu", "", "set scaling cpu utilization target")

	flagSettingsResourcesPtr := flag.Bool("resources", true, "enable resources requisition")
	flagSettingsResourcesMemoryPtr := flag.String("resources-memory", "", "set resource memory ( in GB )")
	flagSettingsResourcesCPUPtr := flag.String("resources-cpu-count", "", "set resource cpu count")

	flagEnvFilePtr := flag.String("env-file", "", "envfile for the config file")
	flagTypePtr := flag.String("type", "", "select the config type")
	flagOutputFilePtr := flag.String("output-file", "", "output file")

	flag.Parse()

	switch *flagTypePtr {
	case "kubernetes-configmap":
		settings := Kubernetes{
			EnvVars:   LoadVars(*flagEnvFilePtr),
			Namespace: *flagSettingsNamespacePtr,
			Name:      *flagSettingsNamePtr,
		}
		settings.Validate()
		WriteToFile(kubernetesConfigMap, settings, *flagOutputFilePtr)
	case "appengine":
		settings := AppEngine{
			EnvVars:            LoadVars(*flagEnvFilePtr),
			Name:               *flagSettingsNamePtr,
			Runtime:            *flagSettingsRuntimePtr,
			Env:                *flagSettingsEnvPtr,
			Command:            *flagSettingsCommandPtr,
			DisableHealthCheck: *flagSettingsDisableHealthCheckPtr,
			ScalingField: AppEngineScalingField{
				MinNumInstances: *flagSettingsScalingMinPtr,
				MaxNumInstances: *flagSettingsScalingMaxPtr,
				CPUUtilization:  *flagSettingsScalingCPUPtr,
			},
			Scaling: *flagSettingsScalingPtr,
			ResourcesField: AppEngineResourcesField{
				MemoryGB: *flagSettingsResourcesMemoryPtr,
				CPUCount: *flagSettingsResourcesCPUPtr,
			},
			Resources: *flagSettingsResourcesPtr,
		}
		settings.Validate()
		WriteToFile(appEngineWeb, settings, *flagOutputFilePtr)
	default:
		log.Fatal("type not found")
	}

}

func LoadVars(filename string) EnvVars {
	if len(filename) < 1 {
		log.Fatal("missing -env-file flag")
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	envvars, err := godotenv.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	return envvars
}

func WriteToFile(tpl *template.Template, s interface{}, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = tpl.Execute(f, s)
	if err != nil {
		log.Fatal(err)
	}
}

func (k *Kubernetes) Validate() {
	if len(k.Namespace) == 0 {
		log.Fatal("missing -namespace flag")
	}
	if len(k.Name) == 0 {
		log.Fatal("missing -name flag")
	}
}

func (a *AppEngine) Validate() {
	if len(a.Name) == 0 {
		log.Fatal("missing -name flag")
	}
	if len(a.Runtime) == 0 {
		log.Fatal("missing -runtime flag")
	}
	if len(a.Env) == 0 {
		log.Fatal("missing -env flag")
	}
	if len(a.Command) == 0 {
		log.Fatal("missing -command flag")
	}

	if a.Scaling {
		if len(a.ScalingField.MinNumInstances) == 0 {
			log.Fatal("missing -scaling-min flag")
		}
		if len(a.ScalingField.MaxNumInstances) == 0 {
			log.Fatal("missing -scaling-max flag")
		}
		if len(a.ScalingField.CPUUtilization) == 0 {
			log.Fatal("missing -scaling-cpu flag")
		}

	}

	if a.Resources {
		if len(a.ResourcesField.MemoryGB) == 0 {
			log.Fatal("missing -resources-memory")
		}
		if len(a.ResourcesField.CPUCount) == 0 {
			log.Fatal("missing -resources-cpu-count")
		}
	}
}
