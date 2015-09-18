
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

TODO(ssoudan) could be a RaspberryPi, an Intel Edison, or an Arduino.


