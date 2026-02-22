package scanner

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
)

// abstracts command execution
type CommandRunner interface {
	Run(name string, args ...string) (io.ReadCloser, error)
}

// executes real OS commands
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

// ARPScanner handles ARP table retrieval
type ARPScanner struct {
	runner CommandRunner
}

// NewARPScanner creates scanner with default runner
func NewARPScanner() *ARPScanner {
	return &ARPScanner{
		runner: DefaultRunner{},
	}
}

// GetARPTable fetches and parses ARP entries
func (a *ARPScanner) GetARPTable() (map[string]string, error) {
	stdout, err := a.runner.Run("arp", "-n")
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	return parseARP(stdout), nil
}

// parseARP parses ARP table from reader
func parseARP(r io.Reader) map[string]string {
	arpMap := make(map[string]string)
	scanner := bufio.NewScanner(r)

	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\s+.*\s+([0-9a-f:]{17})`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			arpMap[matches[1]] = matches[2]
		}
	}
	return arpMap
}
