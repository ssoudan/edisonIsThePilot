
# Autopilot Design

<a href="mailto:sebastien.soudan@gmail.com">Sebastien Soudan</a>

## Table of Contrent

<!-- MarkdownTOC -->

- 1. Problem statement
- 2. System Design
- 3. Heading Control Autopilot

<!-- /MarkdownTOC -->

## 1. Problem statement

For now, we want **a system that can hold the heading by acting on the rudder**.

### 1.1 Requirements

- emergency disconnect
- powerful enough actuator
- as little button as possible
- ability to tune/calibrate the system on-board
  - sinusoidal steering wheel input of know amplitude and frequency -- for different frequencies
  - record/export track
  - ability to change the parameters

### 1.2 Pitfalls?

- slackness in the steering wheel
- slackness in the rudder
- compass issues:
  - non-linearity of the compass
  - tilt compensation(?)
- extreme positions?
- salty env.
- power stability

### 1.3 Resources

The Internet plus few other things.

#### 1.3.1 GPS or compas?

[Sparkfun Forum](https://forum.sparkfun.com/viewtopic.php?f=14&t=31443)
We don't know the declination so the GPS heading and the compass heading will be slightly different.

The compass is sensitive to angular variation though the gyroscope can help to compensate for that.
The compass is also sensitive to the alignement with the movement direction.
The GPS does not provide a meangingful heading when the speed is not enough.


#### 1.3.2 HMC5883L

- [Arduino drivers for HMC5883L](https://github.com/jarzebski/Arduino-HMC5883L)
- [HMC5883L with a RaspberryPi](http://www.instructables.com/id/Interfacing-Digital-Compass-HMC5883L-with-Raspberr/#step1)
- See `HMC5883L_compensation_MPU6050.ino` for a tilt-compensated compass

#### 1.3.3 PID tuning
"Re-tune your PID.  Sounds like one of the parameters is way off.

My understanding is that you start with all parameters 0 and the crank up the P parameter until you get a fast response with minimal overshoot.  Then increase D to eliminate the overshoot without slowing the response too much.  Then increase I to eliminate offset (where it settles close to but not on the desired setpoint).

There is also a PID auto-tune library that might help."
From [Arduino Forum](http://forum.arduino.cc/index.php?topic=232450.0)

#### 1.3.4 Extra

- [AIS on Pi](http://publiclab.org/notes/ajawitz/06-11-2015/raspberry-pi-as-marine-traffic-radar)
- [AIS soft](http://hackaday.com/2013/05/06/tracking-ships-using-software-defined-radio-sdr/)

## 2. System Design

In this section we describe the system that we will build in the long term.


    [Autopilot]
               
                 | position                                 | propulsion
                 v                                          v
    waypoints |-----| target       |-----------| heading |--------| position
    --------> |     | ---------->  | course    | ------> | vessel | -+->
              |-----| course       | autopilot |         |--------|  |
               route               |-----------|                     |
               execution                   ^                         |
                                           \-------------------------/

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

Not quite sure how to compute the error and still have a LTI system -- that we can study its stability.
Thus, for the inital iteration, we will focus on the heading control autopilot describe below.
                                                 
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
                 

## 3. Heading Control Autopilot

For the first iteration, the heading set point will be defined as the current heading when a button is pressed (heading hold mode).

### 3.1 Requirements

<!-- TODO(ssoudan) copy requirements of section 1.1 -->

*Note:* it would be nice to be able to disable the feedback to be able to experimentally identify the dynamic characteristic of rudder-boat system and tune the PID controller from that. 

- Be able to have a sinusoidal steering input of know amplitude
- Measure the track
- Be able to export data

TODO(ssoudan)

Seems that compass calibration might be required here to prevent non-linear behaviors.
Would be nice to be able to disable the feedback to be able to experimentally identify the rudder-boat system and tune the PID controller from that. Which means being able to export data (serial interface?).

### 3.2 Subsystems

#### 3.2.1 Heading measurement
Seems that compass calibration might be required here to prevent non-linear behaviors. It is sensible to pitch and roll and would require a gyro to compensate for this. Also it is sensible to magnetic environnement and would require to be located as far as possible of the engine.

The utilisation of a GPS to get the heading also brings some constraints but we will investiguate this way.


#### 3.2.2 Rudder control
For now, we will assume we don't need a closed-loop control system here and the existing steering chain is fine. But we need to make sure there is as little  slackness as possible in the chain made of the motor, the steering wheel, and the rudder.

<!-- TODO(ssoudan) mechanical interface? -->

### 3.3 Platform and components

### 3.3.1 GPS

The GPS is a [MTK3339 packaged by Adafruit](http://www.adafruit.com/products/746).
It provides NMEA messages via a serial interface at 9600 bauds.

#### 3.3.2 Compass and gyro
If we decide to integrate a compass and gyro, we will use a HMC5883L compass and a MPU6050 gyro.

#### 3.3.3 Platform

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

### 3.4 Detailed Design

#### 3.4.1 GPIO pin multiplexing on Intel Edison
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

Because once configured as output, GPIO and PWM are only 1.8V, we need to level shift it to 3.3V to drive motor driver.
For this, we use the following circuit: 

    -------------+---- 3.3v (J20[pin2])
                 |
                 /
                 \ 10k
    J19[pin2]    /
    1.8v ref     +--------> output 0-3.3v
      _www_____|/
               |\   2N3904
    input        |
    0-1.8v ------+


- [GPIO configuration](http://www.malinov.com/Home/sergey-s-blog/intelgalileo-programminggpiofromlinux)
- [GPIO and sysfs](https://www.kernel.org/doc/Documentation/gpio/sysfs.txt)

#### 3.4.2 Interfacing with HMC5883L and MPU6050 

Arduino has a couple of libraries for these chips: [jarzebski/Arduino-HMC5883L](https://github.com/jarzebski/Arduino-HMC5883L) and 
[jarzebski/Arduino-MPU6050](https://github.com/jarzebski/Arduino-MPU6050).

We will need to write our own implementation of them in Go.
As a support library, [gmcbay/i2c](https://bitbucket.org/gmcbay/i2c) will be use to wrap the i2c buses.

<!-- TODO(ssoudan) which i2c bus do we use? -->

#### 3.4.3 Interfacing with the GPS 
We will use [adrianmo/go-nmea](https://github.com/adrianmo/go-nmea) library to decode the messages and use [tarm/serial](https://github.com/tarm/serial) to access the serial interface. We use `/dev/ttyMFD1` serial interface.

For now, only the [GPRMC](http://aprs.gids.nl/nmea/#rmc) sentence will be used:

    $GPRMC,hhmmss.ss,A,llll.ll,a,yyyyy.yy,a,x.x,x.x,ddmmyy,x.x,a*hh
    1    = UTC of position fix
    2    = Data status (V=navigation receiver warning)
    3    = Latitude of fix
    4    = N or S
    5    = Longitude of fix
    6    = E or W
    7    = Speed over ground in knots
    8    = Track made good in degrees True
    9    = UT date
    10   = Magnetic variation degrees (Easterly var. subtracts from true course)
    11   = E or W
    12   = Checksum

<!-- TODO(ssoudan) need to also consider the sentence with signal quality for the disengagement feature -->

#### 3.4.4 PID controller

[felixge/pidctrl](https://github.com/felixge/pidctrl)
<!-- TODO(ssoudan) describe this -->

#### 3.4.5 Software Architecture

<!-- TODO(ssoudan) -->

#### 3.4.6 Mechanical integration

<!-- TODO(ssoudan) -->

#### 3.4.6 Power source

<!-- TODO(ssoudan) -->

#### 3.4.7 Security

<!-- TODO(ssoudan) -->
- operating conditions
- disengagement
- alarm condition
- handling of recoverable errors

### 3.5 Tests and Validation
<!-- TODO(ssoudan) -->

### 3.5.1 Boundaries 

<!-- TODO(ssoudan) -->
