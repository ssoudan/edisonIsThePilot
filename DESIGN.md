
# Autopilot Design

You can find the latest version <a href="https://github.com/ssoudan/edisonIsThePilot">here</a>.

<a href="https://github.com/ssoudan">Sebastien Soudan</a> --
<a href="https://github.com/philixxx">Philippe Martinez</a>

## Table of Content

<!-- MarkdownTOC -->

- 0 -- To do
- 1 -- Problem statement
- 2 -- System Design
- 3 -- Heading Control Autopilot

<!-- /MarkdownTOC -->

## 0 -- To do

- FUTURE(?) LED that tell the system is powered
- FUTURE(?) LED that tell the system is running (heartbeat)
- FUTURE(?) Proper AP mode with the wifi  

## 1 -- Problem statement

For now, we want **a system that can hold the heading by acting on the rudder**.

DISCLAIMER -- this can hurt and/or cause plenty of other things you don't want. Don't use it unless that's what your are after!

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
The compass is also sensitive to the alignment with the movement direction.
The GPS does not provide a meaningful heading when the speed is too low.


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

## 2 -- System Design

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

We assume the rudder system to be linear.

    [Rudder system]
    steering     |------| rudder angle
    -----------> |   K  | ----------->
    angle        |------|
                 

The boat is assumed to an integrating system:

    [Boat]
    rudder angle |-------|  heading
    -----------> |  Kb/s | ----------->
                 |-------|

## 3 -- Heading Control Autopilot

For the first iteration, the heading set point will be defined as the current heading when a button is pressed (heading hold mode).

### 3.1 Requirements

Let's go over requirements of section 1.1.

- emergency disconnect: the stepper motor will be 'sleeping' mode (no holding torque) all the time but when it is actually driving the motor - two switches will be available: one to enable/disable the autopilot (soft), one to power on/off the system. When the system is powered off, there is no holding torque from the motor either.
- powerful enough actuator: stepper motor with adjustable torque (current limiting on the stepper driver board)
- as little button as possible: 2 ON/OFF switches
- ability to tune/calibrate the system on-board: 
  - sinusoidal steering wheel input of know amplitude and frequency -- for different frequencies
  - record/export track
  - ability to change the parameters

*Note:* it would be nice to be able to disable the feedback to be able to experimentally identify the dynamic characteristic of rudder-boat system and tune the PID controller from that. 

- Be able to have a step steering input of know amplitude
- Measure the track
- Be able to export data

This is implemented in 'cmd/systemCalibration' and exports a JSON with the results to the filesystem or an HTTP endpoint.
The output need to be further processed to extract the characteristics of the steering+boat system in order to tune the PID controller.

### 3.2 Subsystems

#### 3.2.1 Heading measurement
Seems that compass calibration might be required here to prevent non-linear behaviors. It is sensible to pitch and roll and would require a gyro to compensate for this. Also it is sensible to magnetic environment and would require to be located as far as possible of the engine.

The utilisation of a GPS to get the heading also brings some constraints but we will investigate this way.

#### 3.2.2 Rudder control
For now, we will assume we don't need a closed-loop control system here and the existing steering chain is fine. But we need to make sure there is as little play as possible in the chain made of the motor, the steering wheel, and the rudder.

The stepper motor is driven at constant speed for a duration which depends on the requested rotation. Direction of the rotation is defined when the movement is requested. Positive rotation are made in clockwise direction (for the motor). 

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
The main reason for this choice is because we can write Golang. The second reason is because it is the first time I play with this platform. The third reason (which is the first reasonable reason) is because the Linux, x86 architecture, 1GB of ram and on-board wifi plus all the IO pins make it a quite evolutive platform for the job. Would be relatively easy to add mapping, remote control, AIS traffic monitoring features, or a GUI...

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

Currently, we use the following pins:

    For the dashboard informations:  
    NoGPSFix                 gpio40 --> J19 - pin 10
    InvalidGPSData           gpio43 --> J19 - pin 11
    SpeedTooLow              gpio48 --> J19 - pin 6
    HeadingErrorOutOfBounds  gpio82 --> J19 - pin 13
    CorrectionAtLimit        gpio83 --> J19 - pin 14  

    For the alarm:

    alarmGpio                gpio183 --> J18 - pin 8 -- which is also pwm3
    motorDir                 gpio165 --> J18 - pin 2
    motorSleep               gpio12  --> J18 - pin 7
    motorStep                gpio182 --> J17 - pin 1 -- which is pwm2
)

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

Pins driving LEDs uses a more standard NPN based-driver.

    -------------+---- 5 V
                 |
                 /
                 \ 10k
                 /
                 |
                 _
                 V LED
                ---
    pin          |
      _www_____|/
               |\   2N3904
                 |
    0 -----------+  gnd

- [GPIO configuration](http://www.malinov.com/Home/sergey-s-blog/intelgalileo-programminggpiofromlinux)
- [GPIO and sysfs](https://www.kernel.org/doc/Documentation/gpio/sysfs.txt)

#### 3.4.2 Interfacing with HMC5883L and MPU6050 

Arduino has a couple of libraries for these chips: [jarzebski/Arduino-HMC5883L](https://github.com/jarzebski/Arduino-HMC5883L) and 
[jarzebski/Arduino-MPU6050](https://github.com/jarzebski/Arduino-MPU6050).

We will need to write our own implementation of them in Go.
As a support library, [gmcbay/i2c](https://bitbucket.org/gmcbay/i2c) will be use to wrap the i2c buses.

We don't use that yet.

#### 3.4.3 Interfacing with the GPS 
We will use [adrianmo/go-nmea](https://github.com/adrianmo/go-nmea) library to decode the messages and use [tarm/serial](https://github.com/tarm/serial) to access the serial interface. We use `/dev/ttyMFD1` serial interface.

The [GPRMC](http://aprs.gids.nl/nmea/#rmc) sentence will be used:

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

As well as the [GPGGA](http://aprs.gids.nl/nmea/#gga):
    $GPGGA,hhmmss.ss,llll.ll,a,yyyyy.yy,a,x,xx,x.x,x.x,M,x.x,M,x.x,xxxx*hh
    1    = UTC of Position
    2    = Latitude
    3    = N or S
    4    = Longitude
    5    = E or W
    6    = GPS quality indicator (0=invalid; 1=GPS fix; 2=Diff. GPS fix)
    7    = Number of satellites in use [not those in view]
    8    = Horizontal dilution of position
    9    = Antenna altitude above/below mean sea level (geoid)
    10   = Meters  (Antenna height unit)
    11   = Geoidal separation (Diff. between WGS-84 earth ellipsoid and
           mean sea level.  -=geoid is below WGS-84 ellipsoid)
    12   = Meters  (Units of geoidal separation)
    13   = Age in seconds since last update from diff. reference station
    14   = Diff. reference station ID#
    15   = Checksum

First is used for the heading and speed, second is used for the fix quality (field #6).

#### 3.4.4 PID controller

The input of the PID is the error defined as the difference between the current heading as provided by the GPS and the reference heading we have saved right after the autopilot has been enabled.
The error is centered on 0 and varies from -180 (excluded) to 180 (included).
The output of the PID is fed to the steering module which interpret this as the 
rotation to be done in one direction or the other. 

We use the following out own implementation of a PID with filtered derivative.

#### 3.4.5 Software Architecture

The software is architectured around 6 components: 

- gps -- which streams the position, heading, speed and signal quality
- control -- which collects user inputs 
- pilot -- which determine the heading error and correction to apply to the steering 
- steering -- which controls the steering of the boat
- dashboard -- which display notifications
- alarm - which controls the sound alarm

Each component as a single event loop implemented as a go routine.
It reads on the input channels, does what it has to do and send messages to another component. Components are created, wired, started and shutdown in `cmd/edisonIsThePilot.go`.

`conf/conf.go` contains the pin mapping and definition of constants.

`drivers` folder contains the drivers for the I/O subsystem used in this project: gpio, pwm, stepper motor, serial-attached gps.

#### 3.4.6 OS integration

##### Integration with the OS boot
At boot the alarm is on until the edisonIsThePilot is started.
For that we will use systemd:

    # cp edisonIsThePilot.service /lib/systemd/system/
    # systemctl enable edisonIsThePilot

To check the status of the service: 

    # systemctl status edisonIsThePilot -l

<!-- ##### Integration with the OS watchdog -->
<!-- FUTURE(ssoudan) integration with watchdog -->

#### Log rotation
<!-- `edisonIsThePilot` writes both to stderr and '/var/log/edisonIsThePilot.log'. When the program is started or when the size of the file gets greater than 40MB (check performed every `logger.maxWriteCountWithoutCheck` writes), 
the file is rotated to /var/log/edisonIsThePilot.log.old (previous edisonIsThePilot.log.old is deleted) and a new '/var/log/edisonIsThePilot.log' is created. -->
Logs are managed by journalctl. It's configuration is changed to limit the maximum amount of logs it keeps.

They can be watched with: 

    # journalctl -u edisonIsThePilot

#### 3.4.7 Mechanical integration

<!-- TODO(phi) -->

#### 3.4.8 Power source
We use a 12v power source. This is the fed directly to the Vin of the Edison and the Big EasyDriver which control the stepper motor. 
For the other components (GPS and level shifters), we have 2 regulators on the board to obtain 3.3V and 5v regulated.

#### 3.4.9 Security

When the pilot is enabled and detect an error or an over limit condition, the alarm is raised and the (autopilot) steering is disabled. Whenever the system is rebooted/restarted, the alarm is raised before the autopilot is operational (continuous beep) and continues to beep if the autopilot enabled button is ON when it starts. 

### 3.5 Tests and Validation

We have:

- unit/behavorial tests
- standalone programs to test different subsystems that have been used to test the board and its actuators on a bench
- matlab simulations to validate the feasibility of the entire system under some assumptions about the boat and steering chain behavior.

### 3.5.1 Boundaries 

We have different thresholds for that:

- minimum speed -> to cover for inaccurate gps heading
- maximum control angle -> to prevent to rapid correction which could be dangerous
- maximum allowable error -> to detect instabilities and alert the pilot.

### 3.6 Calibration procedure

The purpose of this calibration is to measure the behavior of the controlled system, assess its linearity, and find the parameters of the model that describe it.

#### 3.6.1 Step response

Using 'systemCalibration' which is made of the 'steering' and 'gps' components only.

We first need to make sure 'edisonIsThePilot' service is down.
Then ensure we have enough place for the operation.

Once the boat is going straight at constant speed (cruising speed):
- note the heading of the boat
- start a stopwatch when X degree of steering is added as fast a possible - hold this steering
- every seconds, note the heading of the boat

'systemCalibration' does this procedure automatically.

To test the linearity of the system, multiple such campaign need to be performed for different values of X, on both side and multiple speed if that's relevant.

#### 3.6.2 Frequency response

FUTURE(ssoudan) in the future we might want to do that, but since we don't know the range of frequency of the perturbation we can see we will delay that.





