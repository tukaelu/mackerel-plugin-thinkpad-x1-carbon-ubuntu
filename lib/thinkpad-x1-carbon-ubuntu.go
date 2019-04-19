package mpthinkpad

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	pathBattery = "/sys/class/power_supply/BAT0/uevent"
)

var graphdef = map[string]mp.Graphs{}

// ThinkpadX1CarbonLinuxPlugin plugin struct
type ThinkpadX1CarbonLinuxPlugin struct {
}

// GraphDefinition impl mackerel plugin interface
func (p *ThinkpadX1CarbonLinuxPlugin) GraphDefinition() map[string](mp.Graphs) {
	graphdef["battery.BAT0.capacity"] = mp.Graphs{
		Label: "Battery Capacity",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "capacity", Label: "Capacity", Diff: false, Stacked: false},
		},
	}

	graphdef["battery.BAT0.energy"] = mp.Graphs{
		Label: "Battery Energy",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "energy_now", Label: "Now", Diff: false, Stacked: false},
			{Name: "energy_full", Label: "FCC", Diff: false, Stacked: false},
			{Name: "energy_design", Label: "Design", Diff: false, Stacked: false},
		},
	}

	graphdef["battery.BAT0.cycle"] = mp.Graphs{
		Label: "Battery Cycle Count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "cycle_count", Label: "Cycle", Diff: false, Stacked: false},
		},
	}

	graphdef["cpu.temp"] = mp.Graphs{
		Label: "CPU temperature",
		Unit:  "integer",
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
func (p *ThinkpadX1CarbonLinuxPlugin) FetchMetrics() (map[string]interface{}, error) {
	var err error

	m := make(map[string]interface{})

	if err = collectBattery(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func collectBattery(m *map[string]interface{}) error {
	file, err := os.Open(pathBattery)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		slice := strings.Split(scanner.Text(), "=")

		switch slice[0:1][0] {
		case "POWER_SUPPLY_CAPACITY":
			v, err := atoi(string(slice[1:2][0]))
			if err != nil {
				return err
			}
			(*m)["capacity"] = v
		case "POWER_SUPPLY_CYCLE_COUNT":
			v, err := atoi(string(slice[1:2][0]))
			if err != nil {
				return err
			}
			(*m)["cycle_count"] = v
		case "POWER_SUPPLY_ENERGY_NOW":
			v, err := atoi(string(slice[1:2][0]))
			if err != nil {
				return err
			}
			(*m)["energy_now"] = v
		case "POWER_SUPPLY_ENERGY_FULL":
			v, err := atoi(string(slice[1:2][0]))
			if err != nil {
				return err
			}
			(*m)["energy_full"] = v
		case "POWER_SUPPLY_ENERGY_FULL_DESIGN":
			v, err := atoi(string(slice[1:2][0]))
			if err != nil {
				return err
			}
			(*m)["energy_design"] = v
		}
	}
	return nil
}

func atoi(s string) (int64, error) {
	return strconv.ParseInt(strings.Trim(s, " "), 10, 64)
}

func Do() {

}
