# rp2040-wwvb

It's been on my todo list for years to add a
[WWVB](https://www.nist.gov/pml/time-and-frequency-division/time-distribution/radio-station-wwvb)
recevier to my [GPS clock](https://github.com/jrockway/beaglebone-gps-clock). I got an all-in-one
module to do that which did exactly what I wanted -- receive the radio signal and then output the
time code on a GPIO pin. Unfortunately, it is pretty much impossible to receive the AM signal
reliably on the east coast, so I never got it to work. I knew that WWVB adopted a phase-modulation
signal that was designed to solve that problem, but when I looked, there weren't any modules
available that received the signal. I'm also not sure how to build a 60KHz phase modulation receiver
myself. ("Hook up a random SDR to an adequate antenna and see what happens" didn't yield good
results.) Anyway, I knew there was a wall clock on the market that received the signal, so I bought
it and took it apart. Inside was a blob of epoxy with antenna wires leading into it and two
suspicious-looking traces with pull-up resisitors on them leading out. Plugged in my oscilloscope
and it was obviously I²C, so I went online looking for datasheets. There was
[one](http://www.leapsecond.com/pages/es100/ES100DataSheetver0p97.pdf), and better,
[modules available without a clock](https://www.universal-solder.ca/product/everset-es100-mod-wwvb-bpsk-phase-modulation-receiver-module/)
attached. So I got one of those, connected it to a random microcontroller (I had an RP2040 Feather
laying around), and wrote the code in this repository. It works. I receive the WWVB time signal in
"full minute mode" multiple times per night. (The clock I took apart and wrote my own code to
drive... still works! I ran it on 4V instead of 3, and had to scratch off some solder mask to invade
its bus, and it didn't mind at all.)

I haven't written the code to use the module's "tracking" mode (where it listens to a part of the
signal to find the start of the minute, if you know the time within ±4s), but will add it soon. (I
kind of like the "interactive" development experience, and you can only receive the radio station at
night here, so hacking on this depends on how badly I need to get up in the morning.) That might be
useful to hook into my GPS clock, but I really wanted a 60kHz signal and UTC second pulses to track
the frequency relative to GPS (is it a fixed propagation delay, or does it change?) This module
won't do that, I don't think, but I guess we'll see.

In the meantime, enjoy programming a $1 microcontroller with functions like
`t.Format(time.RFC3339Nano)` which feels very good.

## Building

Build a binary:

    tinygo build -target feather-rp2040 -o out.hex

I flash with a soldered-on SWD header + a JLink EDU + JFlash Lite, but the board I guess supports
UF2 if you want to screw around with USB for 30 minutes every time you change one line of code.
OpenOCD is also an option, which I suppose is great if you don't develop on a Linux VM inside
Windows.

There are some support libraries for driving the screen that you can test with normal go:

    go test ./screen

## Final thoughts

Good module. The RP2040 is also great. I have always been an STM32 fanboy, but for $1 my aliegence
might be changing. Plus they exist during the pandemic because no car manufacturer has them on their
BOMs yet.

Tinygo is also great. Never writing a C program again.
