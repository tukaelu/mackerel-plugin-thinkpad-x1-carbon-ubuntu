package mpthinkpad

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

const (
	pathBattery = "/sys/class/power_supply/BAT0/uevent"

	pathCPU   = "/sys/devices/platform/coretemp.0/hwmon/hwmon3/temp1_input"
	pathCore0 = "/sys/devices/platform/coretemp.0/hwmon/hwmon3/temp2_input"
	pathCore1 = "/sys/devices/platform/coretemp.0/hwmon/hwmon3/temp3_input"
	pathCore2 = "/sys/devices/platform/coretemp.0/hwmon/hwmon3/temp4_input"
	pathCore3 = "/sys/devices/platform/coretemp.0/hwmon/hwmon3/temp5_input"
)

var graphdef = map[string]mp.Graphs{}

// ThinkpadX1CarbonPlugin plugin struct
type ThinkpadX1CarbonPlugin struct {
	Prefix string
}

// GraphDefinition impl mackerel plugin interface
func (p *ThinkpadX1CarbonPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef["battery.capacity"] = mp.Graphs{
		Label: "Battery Capacity",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "capacity", Label: "Capacity", Diff: false, Stacked: false},
		},
	}

	graphdef["battery.energy"] = mp.Graphs{
		Label: "Battery Energy",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "energy_now", Label: "Now", Diff: false, Stacked: false},
			{Name: "energy_full", Label: "FCC", Diff: false, Stacked: false},
			{Name: "energy_design", Label: "Design", Diff: false, Stacked: false},
		},
	}

	graphdef["battery.cycle"] = mp.Graphs{
		Label: "Battery Cycle Count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "cycle_count", Label: "Cycle", Diff: false, Stacked: false},
		},
	}

	graphdef["cpu.temp"] = mp.Graphs{
		Label: "CPU Temperature",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "cpu", Label: "CPU", Diff: false, Stacked: false},
			{Name: "core0", Label: "Core 0", Diff: false, Stacked: false},
			{Name: "core1", Label: "Core 1", Diff: false, Stacked: false},
			{Name: "core2", Label: "Core 2", Diff: false, Stacked: false},
			{Name: "core3", Label: "Core 3", Diff: false, Stacked: false},
		},
	}

	return graphdef
}

// FetchMetrics impl mackerel plugin interface
func (p *ThinkpadX1CarbonPlugin) FetchMetrics() (map[string]float64, error) {

	m := make(map[string]float64)

	if err := collectBattery(&m); err != nil {
		return nil, err
	}

	if err := collectCPUTemp(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// MetricKeyPrefix impl mackerel plugin interface
func (p *ThinkpadX1CarbonPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "thinkpad"
	}
	return p.Prefix
}

func collectBattery(m *map[string]float64) error {
	file, err := os.Open(pathBattery)
	if err != nil {
		return err
	}
	defer file.Close()

	var key string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key = ""

		slice := strings.Split(scanner.Text(), "=")

		switch slice[0:1][0] {
		case "POWER_SUPPLY_CAPACITY":
			key = "capacity"
		case "POWER_SUPPLY_CYCLE_COUNT":
			key = "cycle_count"
		case "POWER_SUPPLY_ENERGY_NOW":
			key = "energy_now"
		case "POWER_SUPPLY_ENERGY_FULL":
			key = "energy_full"
		case "POWER_SUPPLY_ENERGY_FULL_DESIGN":
			key = "energy_design"
		}

		if key == "" {
			continue
		}

		value, err := atof(string(slice[1:2][0]))
		if err != nil {
			return err
		}
		(*m)[key] = value
	}
	return nil
}

func collectCPUTemp(m *map[string]float64) error {
	var v float64
	var err error

	v, err = parseCPUTemp(pathCPU)
	if err != nil {
		return err
	}
	(*m)["cpu"] = v

	v, err = parseCPUTemp(pathCore0)
	if err != nil {
		return err
	}
	(*m)["core0"] = v

	v, err = parseCPUTemp(pathCore1)
	if err != nil {
		return err
	}
	(*m)["core1"] = v

	v, err = parseCPUTemp(pathCore2)
	if err != nil {
		return err
	}
	(*m)["core2"] = v

	v, err = parseCPUTemp(pathCore3)
	if err != nil {
		return err
	}
	(*m)["core3"] = v

	return nil
}

func parseCPUTemp(statFile string) (float64, error) {
	str, err := parseStatFile(statFile)
	if err != nil {
		return 0, err
	}

	v, err := atof(str)
	if err != nil {
		return 0, err
	}
	return v / 1000, nil
}

func parseStatFile(statFile string) (string, error) {
	str, err := ioutil.ReadFile(statFile)
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.Trim(string(str), "\n"), " "), nil
}

func atoi(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func atof(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Do the plugin
func Do() {
	optPrefix := flag.String("metric-key-prefix", "thinkpad", "Metric key prefix")
	flag.Parse()

	p := mp.NewMackerelPlugin(&ThinkpadX1CarbonPlugin{
		Prefix: *optPrefix,
	})
	p.Run()
}
