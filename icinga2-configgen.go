package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/cheggaaa/pb"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"math/rand"
	"os"
	"path"
	"sync"
	"text/template"
)

type Host struct {
	Name string
}

type Service struct {
	Name string
	Host string
}

var (
	numHosts    = kingpin.Flag("hosts", "Number of hosts").Short('H').Required().Int()
	numServices = kingpin.Flag("services", "Number of services per Host").Short('s').Required().Int()
	confDir     = kingpin.Flag("confdir", "Output directory for configs").Short('c').Required().String()
	templateDir = kingpin.Flag("tmpldir", "Template directory (defaults to /etc/icinga2-configgen/templates/ or templates/").PlaceHolder("TEMPLATE-DIR").Default("/etc/icinga2-configgen/").String()
)

var templates = template.New("foo")


func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	// check if config dir exists
	finfo, err := os.Stat(*confDir)
	if err != nil {
		log.Fatalf("Config directory %s does not exist\n", *confDir)
	}
	if !finfo.IsDir() {
		log.Fatalf("Config directory %s is not a directory\n", *confDir)
	}

    finfo, err = os.Stat(*templateDir)
    if err != nil {
        log.Printf("Config directory %s does not exist, try to fallback to templates/\n", *templateDir)
        *templateDir = "templates"
        finfo, err = os.Stat(*templateDir)
        if err != nil {
            log.Fatalf("Failed, no template directory found\n")
        }
    }
    templates = template.Must(template.ParseGlob(*templateDir + "/*.tmpl"))

	log.Printf("Create %d hosts with %d services each", *numHosts, *numServices)
	bar := pb.StartNew(*numHosts)
	// create a syncgroup to make sure all templates are written before exiting
	var wg sync.WaitGroup
	for host := 0; host < *numHosts; host++ {
		// add us to the syncgroup
		wg.Add(1)
		go func() {
			// remove from syncgroup if done
			defer wg.Done()
			genHost()
			bar.Increment()
		}()
	}
	// exit if all subroutines are done
	wg.Wait()
	bar.FinishPrint("The End!")
}

func genHost() {
	hostname := getName()
	filename := path.Join(*confDir, hostname+".conf")

	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not open %s: %s", filename, err)
	}
	// close if done
	defer f.Close()

	err = templates.ExecuteTemplate(f, "host", Host{
		Name: hostname,
	})
	if err != nil {
		log.Fatalf("Could not execute host template: %s", err)
	}
	for service := 0; service < *numServices; service++ {
		name := getName()
		err = templates.ExecuteTemplate(f, "service", Service{
			Name: name,
			Host: hostname,
		})
		if err != nil {
			log.Fatalf("Could not execute service template: %s", err)
		}
	}
}

func getName() string {
	name := randomdata.SillyName() + fmt.Sprintf("%4d", rand.Int())
	return name
}
