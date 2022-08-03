package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Change value : echo 150 | sudo tee /sys/class/backlight/intel_backlight/brightness
// Get max : cat /sys/class/backlight/intel_backlight/max_brightness

// Example in go : https://github.com/psimoesSsimoes/i-brightness/blob/master/main.go

// Sys fs brightness paths
const (
	DefaultSysFSBrightnessPath    = "/sys/class/backlight/intel_backlight/brightness"
	DefaultSysFSMaxBrightnessPath = "/sys/class/backlight/intel_backlight/max_brightness"
	DefaultSysFSLuminanceInput    = "/sys/bus/iio/devices/iio:device0/in_illuminance_raw"
)

func main() {
	// Loag cli args
	set := flag.Int("set", 0, "Set brightness")
	get := flag.Bool("get", true, "Print current brightness to stdout (in percent unless -raw)")
	raw := flag.Bool("raw", false, "Set & get brightness in raw format")
	lum := flag.Bool("lum", false, "Get input luminance in raw format")
	brightnessSysFSPath := flag.String("brpath", "", "Brightness sys fs path (default to /sys/class/backlight/intel_backlight/brightness)")
	maxBrightnessSysFsPath := flag.String("maxbrpath", "", "Max brightness sys fs path (default to /sys/class/backlight/intel_backlight/max_brightness)")
	lumInputBrightnessSysFsPath := flag.String("luminpath", "", "Input luminance sys fs path (default to /sys/bus/iio/devices/iio:device0/in_illuminance_raw)")
	flag.Parse()

	// Set correct sys fs brightness paths
	if len(*brightnessSysFSPath) == 0 {
		*brightnessSysFSPath = DefaultSysFSBrightnessPath
	}

	if len(*maxBrightnessSysFsPath) == 0 {
		*maxBrightnessSysFsPath = DefaultSysFSMaxBrightnessPath
	}

	if len(*lumInputBrightnessSysFsPath) == 0 {
		*lumInputBrightnessSysFsPath = DefaultSysFSLuminanceInput
	}

	// Checks
	if *lum {
		*get = false
	}

	if *lum && (*set != 0 || *get) {
		panic("-lum and -set or -get are mutually exclusive")
	}

	// Get input luminance
	if *lum {
		rawInputLuminance, err := readFileValue(*lumInputBrightnessSysFsPath)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(os.Stdout, "%d\n", rawInputLuminance)
		return
	}

	// Get current brightness
	if *get {
		rawBrightness, err := readFileValue(*brightnessSysFSPath)
		brightness := rawBrightness
		if err != nil {
			panic(err)
		}

		if !*raw {
			maxRawBrightness, err := readFileValue(*maxBrightnessSysFsPath)
			if err != nil {
				panic(err)
			}

			brightness = rawBrightness * 100 / maxRawBrightness
		}

		fmt.Fprintf(os.Stdout, "%d\n", brightness)
	}

	// Set current brightness
	if *set > 0 {
		brightness := *set

		if !*raw {
			maxRawBrightness, err := readFileValue(*maxBrightnessSysFsPath)
			if err != nil {
				panic(err)
			}

			brightness = (brightness * maxRawBrightness / 100)
		}

		if err := writeFileValue(*brightnessSysFSPath, brightness); err != nil {
			panic(err)
		}
	}
}

func readFileValue(filename string) (int, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(strings.TrimSpace(string(content)))
}

func writeFileValue(filename string, value int) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0o444)
	if err != nil {
		return err
	}

	n, err := f.Write([]byte(strconv.Itoa(value)))
	if err == nil && n < len(([]byte(string(value)))) {
		err = io.ErrShortWrite
	}

	if err1 := f.Close(); err == nil {
		err = err1
	}

	return err
}
