# RP2040 Zero TinyGo Test

A playground for testing device drivers and examples on RP2040 using TinyGo.

<video src="output.webp" controls autoplay loop muted></video>

## What's included

- **Rotary encoder** with push switch
- **SSD1306 OLED** display
- **WS2812 RGB LED**

## Setup

### Prerequisites

- [TinyGo](https://tinygo.org/getting-started/)
- RP2040 board (Raspberry Pi Pico or compatible)
- USB cable

### Flash the project

```bash
tinygo flash -target=<board> -port=/dev/ttyACM0 ./
```

Replace `<board>` with your target (e.g., `pico`).

### Read the serial output

**Install picocom** (Linux):

```bash
sudo apt update && sudo apt install picocom
```

**Connect to console**:

```bash
picocom -b 115200 /dev/ttyACM0
```

Exit: `Ctrl-A Ctrl-X`

## Serial output

`main.go` writes a single updating status line to the serial console at 115200 baud:

```go
fmt.Fprintf(machine.Serial, "\r%-60s", line)
```

This overwrites the same line in-place. Use `fmt.Fprintln()` for new lines instead:

```go
fmt.Fprintln(machine.Serial, line)
```

To clear the line before printing:

```go
fmt.Fprintf(machine.Serial, "\r\033[K%s", line)
```

## Example: main.go

### Rotary encoder

- Reads via the `enc.Dir` channel (emits `+1` or `-1` per click)
- Maintains a position counter by accumulating direction values

```go
pos := 0
for {
    select {
    case dir := <-enc.Dir:
        pos += dir
    }
}
```

### Switch / button

- Polls the pin for press state (active-low with pull-up)
- Uses `enc.SwitchWasClicked()` to detect release (debounced)

### OLED display

- Renders text with `tinyfont` and shapes with `tinydraw`
- Updates each frame

### Colors

Color constants must be `var`, not `const` (composite literals):

```go
var WHITE = color.RGBA{255, 255, 255, 255}
```

## Files

- `main.go` — example application
- `go.mod` / `go.sum` — dependencies

## License

None specified.

This repository is a small test project for experimenting with devices on an RP2040 board using TinyGo and Go.

It is intended as a playground to validate drivers and example code for common peripherals such as:

- Rotary encoders (hardware with push switch)
- Small OLED displays (SSD1306 / similar)

Goals
- Demonstrate interrupt-driven and polled input for a rotary encoder
- Show how to display simple state on an OLED
- Provide short, copy-pasteable examples for rapid testing

Quick start
1. Install TinyGo (https://tinygo.org/getting-started/).
2. Connect your RP2040 board (Raspberry Pi Pico or compatible) via USB.
3. Build and flash an example:

   tinygo flash -target=<your-rp2040-board> -port=<your-serial-port> ./

   Replace <your-rp2040-board> and <your-serial-port> with your board and OS-specific port (example: /dev/ttyACM0 or COM3).

Rotary encoder notes
- The included rotary encoder driver exposes a Dir channel which emits +1 or -1 on each full detent (click). Use this channel to maintain a stable position counter:

```go
pos := 0
for {
    select {
    case d := <-enc.Dir:
        pos += d
        println("position:", pos)
    }
}
```

- The driver's Value() method is not a reliable position read because the interrupt handler resets its internal counter after each detent. Prefer the Dir channel.

Switch (button) behavior
- The driver currently sends on enc.Switch only when the button is released (active-low with pull-up).
- To detect press (falling edge) you can:
  - Poll the pin yourself with debouncing (simple and reliable), or
  - Modify the driver to emit a Press channel on the falling edge of the switch interrupt.

Example: display "true"/"false" on OLED when button toggles
- Keep a package-level boolean state and toggle it on each click; redraw the display with fmt.Sprintf("%t", state) or write a short label.

Code style / constants
- Use package-level vars for composite literals such as colors (e.g. `var white = color.RGBA{255,255,255,255}`) — composite literals cannot be declared as const.

Files of interest
- main.go — example application entry
- go.mod / go.sum — module info and dependencies
- drivers/ or vendor/ — driver code (rotary encoder, oled) if present

Contributing
- This project is a small test repo; open issues or PRs with improvements to examples or driver behavior.

License
- No license specified. Add one if you intend to share this project publicly.
