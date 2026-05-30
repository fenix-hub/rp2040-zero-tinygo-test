package main

import (
	"fmt"
	"image/color"
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"

	"github.com/pavelanni/tinygo-drivers/rotaryencoder"
)

const (
	CLK_PIN machine.Pin = machine.GPIO1
	DT_PIN  machine.Pin = machine.GPIO2
	SW_PIN  machine.Pin = machine.GPIO3

	SDA_PIN machine.Pin = machine.GPIO4
	SCL_PIN machine.Pin = machine.GPIO5
)

var (
	pos      int  = 0
	counter  int  = 0
	clicked  bool = false
	pressing bool = false
)

var (
	WHITE = color.RGBA{255, 255, 255, 255}
	BLACK = color.RGBA{A: 255}
)

var display *ssd1306.Device
var enc rotaryencoder.Device
var led ws2812.Device

func init() {
	// Serial
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})

	// Display
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       SDA_PIN,
		SCL:       SCL_PIN,
	})
	display = ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32, // or ssd1306.Address
		Width:   128,
		Height:  32, // or 64
	})
	display.ClearDisplay()

	// Rotary Encoder
	enc = rotaryencoder.New(CLK_PIN, DT_PIN, SW_PIN)
	enc.Configure()

	// LED
	pin := machine.GPIO16
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led = ws2812.NewWS2812(pin)

}

func main() {
	fmt.Fprintln(machine.Serial, "hello from RP2040")

	var x_pos int16 = 0

	for {
		display.ClearBuffer()
		if x_pos > 128 {
			x_pos = -5 * 5
		}
		x_pos += 1
		tinyfont.WriteLine(display, &tinyfont.Org01, x_pos, 5, "hello", WHITE)

		select {
		case dir := <-enc.Dir:
			pos += dir
		default:
			pressing = !SW_PIN.Get()
			if enc.SwitchWasClicked() {
				clicked = !clicked
				counter += 1
			}
		}

		tinyfont.WriteLine(display, &tinyfont.Org01, 0, 11, fmt.Sprintf("rotation: %v", strconv.Itoa(pos)), WHITE)
		tinyfont.WriteLine(display, &tinyfont.Org01, 0, 17, fmt.Sprintf("clicks: %v", strconv.Itoa(counter)), WHITE)
		tinyfont.WriteLine(display, &tinyfont.Org01, 0, 23, fmt.Sprintf("pressing: %v", strconv.FormatBool(pressing)), WHITE)

		tinydraw.Circle(display, int16(128-(128/6)), 16, int16(uint8(pos)%16), WHITE)

		if pressing {
			tinydraw.Rectangle(display, int16(pos%128), 26, 5, 5, WHITE)
			led.WriteColors([]color.RGBA{{255, 0, 0, 255}})
		} else {
			tinydraw.FilledRectangle(display, int16(pos%128), 26, 5, 5, WHITE)
			led.WriteColors([]color.RGBA{{255, 255, 0, 255}})
		}

		display.Display()

		line := fmt.Sprintf("rotation: %d counter: %d pressing: %t", pos, counter, pressing)
		fmt.Fprintf(machine.Serial, "\r%-60s", line) // pads/truncates to 60 chars

		time.Sleep(time.Millisecond * 50)
	}
}
