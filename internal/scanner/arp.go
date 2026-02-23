package scanner

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
)

/******** Command Execution ********/

type CommandRunner interface {
	Run(name string, args ...string) (io.ReadCloser, error)
}

type DefaultRunner struct{}

func (d DefaultRunner) Run(name string, args ...string) (io.ReadCloser, error) {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return stdout, nil
}

/******** ARP Scanner ********/

type ARPScanner struct {
	runner CommandRunner
}

func NewARPScanner() *ARPScanner {
	return &ARPScanner{
		runner: DefaultRunner{},
	}
}

func (a *ARPScanner) GetARPTable() (map[string]string, error) {
	stdout, err := a.runner.Run("arp", "-n")
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	return parseARP(stdout), nil
}

/******** Package-Level Helper ********/

func GetARPTable() (map[string]string, error) {
	return NewARPScanner().GetARPTable()
}

/******** Parsing ********/

func parseARP(r io.Reader) map[string]string {
	arpMap := make(map[string]string)
	scanner := bufio.NewScanner(r)

	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\s+.*\s+([0-9a-f:]{17})`)

	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) == 3 {
			arpMap[matches[1]] = matches[2]
		}
	}
	return arpMap
}
