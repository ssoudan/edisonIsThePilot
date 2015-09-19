
# Autopilot Design

<a href="mailto:sebastien.soudan@gmail.com">Sebastien Soudan</a>

## Problem statement

For now, we want **a system that can hold the heading by acting on the rudder**.

## Requirements

- emergency disconnect
- powerful enough actuator
- as little button as possible
- ability to tune/calibrate the system on-board
  - sinusoidal steering wheel input of know amplitude and frequency -- for different frequencies
  - record/export track
  - ability to change the parameters

## Pitfalls?

- slackness in the steering wheel
- slackness in the rudder
- compass issues:
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

In this section we describe the system that we will build in the long term.


    [Autopilot]
               
                 | position                                 | propulsion
                 v                                          v
    waypoints |-----| target       |-----------| heading |------| position
    --------> |     | ---------->  | course    | ------> | boat | --->
              |-----| course       | autopilot |         |------|  |
               route               |-----------|                   |
               execution                   ^                       |
                                           \-----------------------/

Initially, we will concentrate on the right part of this diagram and design the course 
autopilot.

    [Course autopilot]
                   error              heading                          position
    course  /-----\     |------------|       |----------|    |--------| 
    ------> |error| --> | Controller | ----> | Steering |--->| vessel |--+-->
            \-----/     |------------|       |----------|    |--------|  |
               ^                                                         |
               |                                                         | 
               \---------------------------------------------------------/

Though, not quite sure how to compute the error and still have a LTI system - that we can study its stability.
For the inital phase we will thus focus on the heading control autopilot describe below.
                                                 
    [Heading control]
                                    
                                      
                          error       steering    rudder angle      actual 
    heading (sp)   /-----\      |-----|      |--------|    |------| heading
    -------------> | +/- | ---> | PID | ---> | rudder | -> | boat |--+---->
                   \-----/      |-----|      |--------|    |------|  |
                      ^                                              |
                      |                                              |
                      |       |--------------|                       |
                      \-------| compass/gps? |<----------------------/
                              |--------------|

// TODO(ssoudan) need to figure out how the steering wheel/rudder system works.

    [Rudder system]
    steering     |------| rudder angle
    -----------> |   K  | ----------->
    angle        |------|
                 

### Autopilot

For the first iteration, the heading set point will be defined as the current heading when a button is pressed (heading hold mode).

### Heading control
Seems that compass calibration might be required here to prevent non-linear behaviors. It is sensible to pitch and roll and would require a gyro to compensate for this. Also it is sensible to magnetic environnement and would require to be located as far as possible of the engine.

The utilisation of a GPS to get the heading also brings some constraints but we will investiguate this way.

*Note:* it would be nice to be able to disable the feedback to be able to experimentally identify the dynamic characteristic of rudder-boat system and tune the PID controller from that. 

**Requirements:**

- Be able to have a sinusoidal steering input of know amplitude
- Measure the track
- Be able to export data

### Rudder control
For now, we will assume we don't need a closed-loop control system here and the existing steering chain is fine. But we need to make sure there is as little  slackness as possible in the chain made of the motor, the steering wheel, and the rudder.

## Components

### microcontroler/computer

Could be a RaspberryPi, an Intel Edison, or an Arduino. But for the reason down below, it will be an Edison -- plus would have to change the name of the project.

We need to support the following: 

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