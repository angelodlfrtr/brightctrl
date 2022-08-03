# Get / set screen brightness on linux systems (via sysfs)

## Installation

```sh
go install github.com/angelodlfrtr/brightctrl@latest
```

Non-root users cannot by default write to /sys/class, to allow other users :

```sh
sudo chown root $GOPATH/bin/brightctrl
sudo chmod u+s $GOPATH/bin/brightctrl
```

## Usage

```sh
# Print help
brightctrl --help

# Get current brightness
brightctrl

# Get current brightness in raw
brightctrl -raw

# Set current brightness in percent
brightctrl -set 50

# Set current brightness in raw
brightctrl -set 488 -raw
```
