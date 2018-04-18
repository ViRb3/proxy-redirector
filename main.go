package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"
)

var settingsFile string
var lineRegex = regexp.MustCompile(`^([0-9.*]+:[0-9*]+)[\s\t]+([0-9.]+:[0-9]+)$`)

func main() {
	var port int
	var verbose bool
	var help bool
	flag.IntVar(&port, "port", 8868, "Port to listen on")
	flag.StringVar(&settingsFile, "settings", "settings.txt", "Settings file with routes")
	flag.BoolVar(&verbose, "verbose", true, "Verbose proxy output")
	flag.BoolVar(&help, "help", false, "Help screen")
	flag.Parse()

	if help {
		helpScreen()
		return
	}

	var routes, err = readSettings()
	if err != nil {
		return
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	for src, dst := range routes {
		setupRedirect(proxy, src, dst)
	}

	fmt.Println()
	fmt.Printf("HTTP/S proxy up on port %d!", port)
	fmt.Println()
	fmt.Println()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), proxy))
}

func setupRedirect(proxy *goproxy.ProxyHttpServer, src string, dst string) {
	var pattern = getHostMatchRegex(src)
	proxy.OnRequest(goproxy.UrlMatches(pattern)).HandleConnectFunc(
		func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			return goproxy.OkConnect, dst
		})
	fmt.Printf("Redirecting %s -> %s", src, dst)
	fmt.Println()
}

func getHostMatchRegex(src string) *regexp.Regexp {
	var srcSplit = strings.Split(src, ":")
	if srcSplit[0] == "*" && srcSplit[1] == "*" {
		return regexp.MustCompile(".+")
	} else if srcSplit[0] == "*" {
		return regexp.MustCompile(fmt.Sprintf("^.+:%s$", srcSplit[1]))
	} else if srcSplit[1] == "*" {
		return regexp.MustCompile(fmt.Sprintf("^%s:.+$", srcSplit[0]))
	} else {
		return regexp.MustCompile(fmt.Sprintf("^%s:%s$", srcSplit[0], srcSplit[1]))
	}
}

func readSettings() (map[string]string, error) {
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		err := errors.New("settings file doesn't exist")
		printError(err)
		return nil, err
	}
	bytes, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return nil, err
	}
	parsed, err := parseSettings(string(bytes))
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func parseSettings(content string) (map[string]string, error) {
	var result = make(map[string]string)
	for _, line := range splitLines(content) {
		if strings.TrimSpace(line) == "" {
			continue
		}
		matches := lineRegex.FindStringSubmatch(line)
		if len(matches) < 2 {
			err := errors.New("bad settings format")
			printError(err)
			return nil, err
		}
		result[matches[1]] = matches[2]
	}
	return result, nil
}

func helpScreen() {
	fmt.Println(`A HTTP/S proxy that redirects connections.
Designed to be used as system proxy or forced for specific programs via software like Proxifier.`)

	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()

	fmt.Println(`Settings file consists of lines defining redirection routes. 

A redirection route has the following format:
{ip}:{port} {ip}:{port}

Multiple whitespaces/tabs are permitted as a separator.

This program will redirect the first (source) ip&port to the second (destination) ip&port.
You can use a wildcard '*' to match the ip, port, or both, for the source.`)
}
