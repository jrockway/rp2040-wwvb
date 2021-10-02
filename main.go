package main

import (
	"machine"
	"time"

	"github.com/jrockway/rp2040-wwvb/screen"
	"tinygo.org/x/drivers/ssd1306"
)

const (
	TWakeup = time.Millisecond
)

func main() {
	enable := machine.D5
	irq := machine.D6
	enable.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	irq.Configure(machine.PinConfig{
		Mode: machine.PinInput,
	})
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Low()

	machine.I2C1.Configure(machine.I2CConfig{
		SDA:       machine.I2C1_SDA_PIN,
		SCL:       machine.I2C1_SCL_PIN,
		Frequency: machine.TWI_FREQ_400KHZ,
	})
	disp := ssd1306.NewI2C(machine.I2C1)
	disp.Configure(ssd1306.Config{
		Width:    128,
		Height:   32,
		VccState: ssd1306.SWITCHCAPVCC,
		Address:  0x3c,
	})
	s := screen.Screen{
		Display: &disp,
	}
	s.Clear()

	start := func() {
		enable.High()
		time.Sleep(TWakeup)
		s.Printf(".")
		if err := machine.I2C1.WriteRegister(0x32, 0x00, []byte{0b00000001}); err != nil {
			time.Sleep(time.Second)
			s.Printf("control: %v", err)
		}
		s.Printf(".")
		out := make([]byte, 2)
		if err := machine.I2C1.ReadRegister(0x32, 0x00, out[0:1]); err != nil {
			s.Printf("read control: %v", err)
		}
		s.Printf("ctrl %x ", out[0])
		if err := machine.I2C1.ReadRegister(0x32, 0x02, out[0:1]); err != nil {
			s.Printf("read irq: %v", err)
		}
		s.Printf("irq %x ", out[0])
		if err := machine.I2C1.ReadRegister(0x32, 0x03, out[0:1]); err != nil {
			s.Printf("read status: %v", err)
		}
		s.Printf("st0 %x ", out[0])
		if err := machine.I2C1.ReadRegister(0x32, 0x0D, out[0:1]); err != nil {
			s.Printf("read deviceid: %v", err)
		}
		s.Printf("id %x ", out[0])
	}
	start()

	var cycles, nok int
	var lastTime time.Time
	out := make([]byte, 9)
	for i := 0; ; i++ {
		s.Printf("\nn %d nok %d cyc %d ", i, nok, cycles)
		if !lastTime.IsZero() {
			s.Printf("%s\n", lastTime.Format(time.RFC3339Nano))

		}
		var gotIRQ time.Time
		for j := 0; j < 500; j++ {
			time.Sleep(10 * time.Millisecond)
			if !irq.Get() {
				s.Printf("irq.")
				gotIRQ = time.Now()
				break
			}
		}
		if !gotIRQ.IsZero() {
			cycles++
			if err := machine.I2C1.ReadRegister(0x32, 0x02, out[0:1]); err != nil {
				s.Printf("read irq: %v", err)
				continue
			}
			if out[0] == 0x04 {
				s.Printf("cycle complete.")
			} else if out[0] == 0x01 {
				s.Printf("rx complete.")
			} else {
				s.Printf("unk irq %03b.", out[0])
			}
		}
		if err := machine.I2C1.ReadRegister(0x32, 0x03, out[0:1]); err != nil {
			s.Printf("read status: %v", err)
			continue
		}
		s.Printf("ant %x rxok %x ", out[0]&0b00000010>>1, out[0]&0b00000001)

		if out[0]&00000001 > 0 {
			s.Printf("TIME! ")
			machine.LED.High()
			if err := machine.I2C1.ReadRegister(0x32, 0x04, out); err != nil {
				s.Printf("read time: %v", err)
				continue
			}
			year := fromBCD(out[0])
			month := fromBCD(out[1])
			day := fromBCD(out[2])
			hour := fromBCD(out[3])
			minute := fromBCD(out[4])
			second := fromBCD(out[5])
			lastTime = time.Date(2000+year, time.Month(month), day, hour, minute, second, 0, time.UTC).Add(time.Since(gotIRQ).Round(time.Millisecond))

			cycles = 0
			nok++
			enable.Low()
			time.Sleep(time.Second)
			start()
		}
	}
}

func fromBCD(x byte) int {
	return int(10*(x>>4)) + int(0b1111&x)
}
