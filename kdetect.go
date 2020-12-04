package kdetect

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/pilebones/go-udev/crawler"
)

func charsToString(ca [65]int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}

func uname() string {
	utsname := syscall.Utsname{}
	syscall.Uname(&utsname)
	return charsToString(utsname.Release)
}

type moduleAlias map[*regexp.Regexp]string

func (a moduleAlias) findModule(alias string) string {
	for k, v := range a {
		if k.MatchString(alias) {
			return v
		}
	}
	return ""
}

func readAlias(alias string) (moduleAlias, error) {
	res := moduleAlias{}

	file, err := os.Open(alias)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")
		// modalias contains '?' which aren't valid golang regexes
		moduleAlias := strings.ReplaceAll(data[1], "?", ".*")
		module := data[2]
		r, err := regexp.Compile(moduleAlias)
		if err != nil {
			log.Println("Failed compiling", moduleAlias)
			continue
		}
		res[r] = module
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func unique(intSlice []string) []string {
	keys := make(map[string]interface{})
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = nil
			list = append(list, entry)
		}
	}
	return list
}

func ProbeKernelModules(aliasFile string) ([]string, error) {
	queue := make(chan crawler.Device)
	errors := make(chan error)

	var kernelModules []string

	crawler.ExistingDevices(queue, errors, nil)

	if aliasFile == "" {
		aliasFile = filepath.Join(string(filepath.Separator), "lib", "modules", uname(), "modules.alias")
	}

	// Parse modules.alias, and pre-compile regexes for each driver
	mods, err := readAlias(aliasFile)
	if err != nil {
		return kernelModules, err
	}

	// Consume the events, and find modules to load as needed
	for {
		select {
		case device, more := <-queue:
			if !more {
				kernelModules = unique(kernelModules)
				return kernelModules, nil
			}
			if alias, ok := device.Env["MODALIAS"]; ok {
				driver := mods.findModule(alias)
				if driver != "" {
					kernelModules = append(kernelModules, driver)
				}
			}
			if driver, ok := device.Env["DRIVER"]; ok {
				kernelModules = append(kernelModules, driver)
			}
		case err := <-errors:
			log.Println("ERROR:", err)
		}
	}
}
