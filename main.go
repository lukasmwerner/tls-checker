package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/lukasmwerner/tls-checker/scan"
)

func main() {

	hostname := ""

	if len(os.Args) == 2 {
		hostname  = os.Args[1]
	} else {
		err := huh.NewInput().Title("Hostname please").Value(&hostname).Run()
		if err != nil {
			return
		}
	}

	p := tea.NewProgram(initalModel(hostname))

	go func() {
		link, _ := url.Parse(hostname)
		if link.Hostname() == "" {
			link.Host = hostname
		}

		c := scan.TestAllTLSVersions(link)
		for v := range c {
			if version, ok := strings.CutPrefix(v, "ok: "); ok {
				p.Send(successMsg(version))
			} else {
				p.Send(failureMsg(v))
			}
		}
		p.Send(tea.QuitMsg{})

	}() 


	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
