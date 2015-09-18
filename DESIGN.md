
# Autopilot Design

<a href="mailto:sebastien.soudan@gmail.com">Sebastien Soudan</a>

## Purpose

Heading hold for now.

## Requirements

- emergency disconnect
- powerful enough actuator
- ability to tune/calibrate the system on-board

## Pitfalls?

- slackness in the steering wheel
- slackness in the rudder
- non-linearity of the compass
- tilt compensation(?)
- extreme positions?
- salty env.
- power stability

## Resources

### GPS or compas?

[Sparkfun Forum](https://forum.sparkfun.com/viewtopic.php?f=14&t=31443)
We don't know the declination so the GPS heading and the compass heading will be slightly different.

The compass is sensitive to angular variation though the gyroscope can help to compensate for that.
The compass is also sensitive to the alignement with the movement direction.
The GPS does not provide a meangingful heading when the speed is not enough.


### HMC5883L

- [Arduino drivers for HMC5883L](https://github.com/jarzebski/Arduino-HMC5883L)
- [HMC5883L with a RaspberryPi](http://www.instructables.com/id/Interfacing-Digital-Compass-HMC5883L-with-Raspberr/#step1)
- See `HMC5883L_compensation_MPU6050.ino` for a tilt-compensated compass

### PID tuning
"Re-tune your PID.  Sounds like one of the parameters is way off.

My understanding is that you start with all parameters 0 and the crank up the P parameter until you get a fast response with minimal overshoot.  Then increase D to eliminate the overshoot without slowing the response too much.  Then increase I to eliminate offset (where it settles close to but not on the desired setpoint).

There is also a PID auto-tune library that might help."
From [Arduino Forum](http://forum.arduino.cc/index.php?topic=232450.0)

### Extra

- [AIS on Pi](http://publiclab.org/notes/ajawitz/06-11-2015/raspberry-pi-as-marine-traffic-radar)
- [AIS soft](http://hackaday.com/2013/05/06/tracking-ships-using-software-defined-radio-sdr/)

## Design

    [Autopilot]                    
                        -------
                        | gps |
                        -------
                           | position
                           v                                   
     destination point  -------   heading (sp)  -------------------     actual
                ---->   | cpu | --------------->| Heading control |---> heading
                        -------                 ------------------- 
                                                 
    [Heading control]
                                    rudder angle             actual     
      heading (sp)   /-----\  e -------   -----------------  heading 
    ---------------> | +/- | ---| PID |---| rudder | boat |--+---->
                     \-----/    -------   -----------------  |
                        ^                                    |
                        |                                    |
                        |        ----------                  |
                        \-------| compass |<-----------------/
                                 ----------

    [Rudder control]
            rudder angle (sp)  /-----\    -------    ------------------------      rudder angle
            -----------------> | +/- |--->| PID |--->| motor/steering wheel |----+------>
                               \-----/    -------    ------------------------    |
                                  ^                                              |
                                  |                 --------------------------   |
                                  \-----------------| rudder position sensor |<--/
                                                    --------------------------


### Autopilot

For the first iteration, the heading set point will be defined as the current heading when a button is pressed (heading hold mode).

### Heading control

TODO(ssoudan)

Seems that compass calibration might be required here to prevent non-linear behaviors.
Would be nice to be able to disable the feedback to be able to experimentally identify the rudder-boat system and tune the PID controller from that. Which means being able to export data (serial interface?).

### Rudder control
Here we probably don't need a PID. A simple P=1 is enough unless we have inertia in the control of the position of the rudder.
But we need to remove slackness in the chain made of the motor, the rudder and the sensor.

## Components

### microcontroler/computer

Could be a RaspberryPi, an Intel Edison, or an Arduino. But for the reason down below, it will be an Edison -- plus would have to change the name of the project.

We need to support: 

- get messages from the GPS, 
- be able to write to GPIO (motor direction)
- be able to generate PWM signals (motor rotation)
- be able to act as a i2c master (compass and gyroscope)
- be able to read from GPIO (button)

We will use an Intel Edison for this project, and can write Golang for this platform. We will need to find a couple of libraries to help us.
The main reason for this choice is because we can write Golang. The second reason is because it is the first time I play with this platform. The third reason (which is the first reasonable reason) is beacuse the Linux, x86 architecture, 1GB of ram and on-board wifi plus all the IO pins make it a quite evolutive platform for the job. Would be relatively easy to add mapping, remote control, AIS traffic monitoring features, or a GUI...

Note the Edison's wifi has an [AP mode](https://software.intel.com/en-us/getting-started-with-ap-mode-for-intel-edison-board).

### Compass and gyro

For that we will use an HMC5883L as the compass and an MPU6050 for the gyro.

### GPIO pin multiplexing on Intel Edison
We are using a mini breakout board for the Intel Edison. This has a limited number of pins. But some of them are multiplexed and via configuration we can decide which pin does what [GPIO pin multiplexing guide](http://www.emutexlabs.com/project/215-intel-edison-gpio-pin-multiplexing-guide).

We need: 

- an i2c bus for the HMC5883 and the MPU6050
- a serial interface for the GPS
- a PWM pin for the motor rotation
- a GPIO output pin for motor direction
- a GPIO input pin for the hold heading button
- a couple of GPIOU output pin to control status LEDs

<!-- TODO(ssoudan) pin map -->
<!-- TODO(ssoudan) custom lib for GPIO/PWM + references -->

- [GPIO configuration](http://www.malinov.com/Home/sergey-s-blog/intelgalileo-programminggpiofromlinux)
- [GPIO and sysfs](https://www.kernel.org/doc/Documentation/gpio/sysfs.txt)

### HMC5883L and MPU6050

Arduino has a couple of libraries for these chips: [jarzebski/Arduino-HMC5883L](https://github.com/jarzebski/Arduino-HMC5883L) and 
[jarzebski/Arduino-MPU6050](https://github.com/jarzebski/Arduino-MPU6050).

We will need to write our own implementation of them in Go.
As a support library, [gmcbay/i2c](https://bitbucket.org/gmcbay/i2c) will be use to wrap the i2c buses.

<!-- TODO(ssoudan) which i2c bus do we use? -->

### GPS

The GPS is a [MTK3339 packaged by Adafruit](http://www.adafruit.com/products/746).
It provides NMEA messages via a serial interface at 9600 bauds.
We will use [adrianmo/go-nmea](https://github.com/adrianmo/go-nmea) library to decode the messages and use [tarm/serial](https://github.com/tarm/serial) to access the serial interface. We use `/dev/ttyMFD1` serial interface.

## PID controller

[felixge/pidctrl](https://github.com/felixge/pidctrl)
<!-- TODO(ssoudan) describe this -->