package cmd

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type process struct {
	pid int
	mem int
	cpu float64
}

type processList []process

func (p processList) Len() int {
	return len(p)
}

func (p processList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type byMem struct{ processList }

func (p byMem) Less(i, j int) bool {
	return p.processList[i].mem > p.processList[j].mem
}

type byCPU struct{ processList }

func (p byCPU) Less(i, j int) bool {
	return p.processList[i].cpu > p.processList[j].cpu
}

// TopByMem 列出系统中最占内存的 n 个进程
func TopByMem(n int) ([]process, error) {
	pl, err := parseProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to parse processes: %v", err)
	}

	sort.Sort(byMem{pl})
	if len(pl) > n {
		pl = pl[:n]
	}

	return pl, nil
}

// TopByCPU 列出系统中最占 CPU 的 n 个进程
func TopByCPU(n int) ([]process, error) {
	pl, err := parseProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to parse processes: %v", err)
	}

	sort.Sort(byCPU{pl})
	if len(pl) > n {
		pl = pl[:n]
	}

	return pl, nil
}

// parseProcesses 解析进程信息
func parseProcesses() ([]process, error) {
	cmd := exec.Command("ps", "-eo", "pid,rss,%cpu")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	pl := make(processList, len(lines)-1)
	for i, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		mem, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}
		pl[i-1] = process{pid, mem, cpu}
	}

	return pl, nil
}
