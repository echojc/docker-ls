package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var maybeConfigPath string
	if len(os.Args) > 1 {
		maybeConfigPath = os.Args[1]
	}

	config, err := LoadConfig(maybeConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		containers, err := Containers()
		if err != nil {
			fmt.Fprintf(w, "Error retrieving containers")
		} else {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<h2>Containers</h2>")
			for _, c := range containers {
				image := strings.Split(c.Image, ":")[0]
				if _, contains := config.Blacklist[image]; contains {
					continue
				}

				if len(c.PublicPorts()) == 0 {
					continue
				}

				fmt.Fprintf(w, "<ul>")
				fmt.Fprintf(w, "<li>%s</li>", FormatContainer(c, config))
				fmt.Fprintf(w, "</ul>")
			}
			fmt.Fprintf(w, "<h2>Config</h2>")
			fmt.Fprintf(w, "<pre>")
			fmt.Fprintf(w, config.String())
			fmt.Fprintf(w, "</pre>")
		}
	})

	listenAddr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func FormatContainer(c *Container, config *Config) string {
	name := c.Id[0:12]
	for _, n := range c.Names {
		if strings.Index(n[1:], "/") == -1 {
			name = n[1:]
			break
		}
	}

	pubPorts := c.PublicPorts()
	urls := make([]string, len(pubPorts))
	image := strings.Split(c.Image, ":")[0]

	for i, port := range pubPorts {
		urls[i] = fmt.Sprintf("%s:%d", config.Host, port)

		if protocols, ok := config.Protocols[image]; ok {
			if protocol, ok := protocols[port]; ok {
				url := fmt.Sprintf("%s://%s:%d", protocol, config.Host, port)
				switch protocol {
				case "http", "https":
					urls[i] = fmt.Sprintf("<a href=\"%s\">%s</a>", url, url)
				default:
					urls[i] = url
				}
			}
		}

		urls[i] = fmt.Sprintf("<li>%s</li>", urls[i])
	}

	return fmt.Sprintf("<b>%s</b> (%s)<ul>%s</ul>", name, c.Image, strings.Join(urls, ""))
}
