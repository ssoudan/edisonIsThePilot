/*
* @Author: Sebastien Soudan
* @Date:   2015-09-18 23:39:53
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-19 01:19:09
 */

package hmc

import (
	"bitbucket.org/gmcbay/i2c"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"

	"bytes"
	"encoding/binary"
)

var log = logger.Log("hmc")

const (
	HMC5883L_ADDRESS = (0x1E)
)

const (
	HMC5883L_REG_CONFIG_A = (0x00)
	HMC5883L_REG_CONFIG_B = (0x01)
	HMC5883L_REG_MODE     = (0x02)
	HMC5883L_REG_OUT_X_M  = (0x03)
	HMC5883L_REG_OUT_X_L  = (0x04)
	HMC5883L_REG_OUT_Z_M  = (0x05)
	HMC5883L_REG_OUT_Z_L  = (0x06)
	HMC5883L_REG_OUT_Y_M  = (0x07)
	HMC5883L_REG_OUT_Y_L  = (0x08)
	HMC5883L_REG_STATUS   = (0x09)
	HMC5883L_REG_IDENT_A  = (0x0A)
	HMC5883L_REG_IDENT_B  = (0x0B)
	HMC5883L_REG_IDENT_C  = (0x0C)
)

// hmc5883l_samples_t
const (
	HMC5883L_SAMPLES_8 = 3 //0b11
	HMC5883L_SAMPLES_4 = 2 //0b10
	HMC5883L_SAMPLES_2 = 1 //0b01
	HMC5883L_SAMPLES_1 = 0 //0b00
)

// hmc5883l_dataRate_t
const (
	HMC5883L_DATARATE_75HZ    = 6 //0b110
	HMC5883L_DATARATE_30HZ    = 5 //0b101
	HMC5883L_DATARATE_15HZ    = 4 //0b100
	HMC5883L_DATARATE_7_5HZ   = 3 //0b011
	HMC5883L_DATARATE_3HZ     = 2 //0b010
	HMC5883L_DATARATE_1_5HZ   = 1 //0b001
	HMC5883L_DATARATE_0_75_HZ = 0 //0b000
)

// hmc5883l_range_t
const (
	HMC5883L_RANGE_8_1GA  = 7 //0b111
	HMC5883L_RANGE_5_6GA  = 6 //0b110
	HMC5883L_RANGE_4_7GA  = 5 //0b101
	HMC5883L_RANGE_4GA    = 4 //0b100
	HMC5883L_RANGE_2_5GA  = 3 //0b011
	HMC5883L_RANGE_1_9GA  = 2 //0b010
	HMC5883L_RANGE_1_3GA  = 1 //0b001
	HMC5883L_RANGE_0_88GA = 0 //0b000
)

// hmc5883l_mode_t
const (
	HMC5883L_IDLE      = 2 //0b10
	HMC5883L_SINGLE    = 1 //0b01
	HMC5883L_CONTINOUS = 0 //0b00
)

type Vector struct {
	XAxis float32
	YAxis float32
	ZAxis float32
}

type HMC5883L struct {
	mgPerDigit float32
	v          Vector
	xOffset    float32
	yOffset    float32
	bus        byte

	/*
		// Vector readRaw(void);
	*/

	/*
		// hmc5883l_range_t getRange(void);

		// hmc5883l_mode_t getMeasurementMode(void);
	*/

	/*
		// hmc5883l_dataRate_t getDataRate(void);
	*/

	/*
		// hmc5883l_samples_t getSamples(void);
	*/
}

func (hmc HMC5883L) SetMeasurementMode(mode byte /*hmc5883l_mode_t*/) error {
	value, err := hmc.readRegister8(HMC5883L_REG_MODE)
	if err != nil {
		return err
	}
	value &= 0xFC // 0b11111100
	value |= mode

	return hmc.writeRegister8(HMC5883L_REG_MODE, value)
}

func (hmc HMC5883L) ReadNormalize() (Vector, error) {
	xValue, err := hmc.readRegister16(HMC5883L_REG_OUT_X_M)
	if err != nil {
		return Vector{}, err
	}
	yValue, err := hmc.readRegister16(HMC5883L_REG_OUT_Y_M)
	if err != nil {
		return Vector{}, err
	}
	zValue, err := hmc.readRegister16(HMC5883L_REG_OUT_Z_M)
	if err != nil {
		return Vector{}, err
	}

	hmc.v.XAxis = (float32(xValue) - hmc.xOffset) * hmc.mgPerDigit
	hmc.v.YAxis = (float32(yValue) - hmc.yOffset) * hmc.mgPerDigit
	hmc.v.ZAxis = float32(zValue) * hmc.mgPerDigit

	return Vector{hmc.v.XAxis, hmc.v.YAxis, hmc.v.ZAxis}, nil
}

func (hmc *HMC5883L) SetOffset(xo float32, yo float32) {
	hmc.xOffset = xo
	hmc.yOffset = yo
}

func (hmc *HMC5883L) SetRange(r byte /* hmc5883l_range_t */) {

	switch r {
	case HMC5883L_RANGE_0_88GA:
		hmc.mgPerDigit = 0.073
		break

	case HMC5883L_RANGE_1_3GA:
		hmc.mgPerDigit = 0.92
		break

	case HMC5883L_RANGE_1_9GA:
		hmc.mgPerDigit = 1.22
		break

	case HMC5883L_RANGE_2_5GA:
		hmc.mgPerDigit = 1.52
		break

	case HMC5883L_RANGE_4GA:
		hmc.mgPerDigit = 2.27
		break

	case HMC5883L_RANGE_4_7GA:
		hmc.mgPerDigit = 2.56
		break

	case HMC5883L_RANGE_5_6GA:
		hmc.mgPerDigit = 3.03
		break

	case HMC5883L_RANGE_8_1GA:
		hmc.mgPerDigit = 4.35
		break

	default:
		break
	}

	hmc.writeRegister8(HMC5883L_REG_CONFIG_B, r<<5)
}

func (hmc HMC5883L) SetDataRate(dataRate byte /*hmc5883l_dataRate_t */) error {

	value, err := hmc.readRegister8(HMC5883L_REG_CONFIG_A)
	if err != nil {
		return err
	}
	value &= 0xE3 // 0b11100011
	value |= (dataRate << 2)

	return hmc.writeRegister8(HMC5883L_REG_CONFIG_A, value)
}

func (hmc HMC5883L) SetSamples(samples byte /*hmc5883l_samples_t*/) error {

	value, err := hmc.readRegister8(HMC5883L_REG_CONFIG_A)
	if err != nil {
		return err
	}
	value &= 0x9F // 0b10011111
	value |= (samples << 5)

	return hmc.writeRegister8(HMC5883L_REG_CONFIG_A, value)
}

func New(bus byte) HMC5883L {
	return HMC5883L{bus: bus}
}

func (hmc HMC5883L) Begin() bool {

	a, err := hmc.fastRegister8(HMC5883L_REG_IDENT_A)
	if err != nil {
		log.Error("Failed to read IDENT_A: %v", err)
		return false
	}
	b, err := hmc.fastRegister8(HMC5883L_REG_IDENT_B)
	if err != nil {
		log.Error("Failed to read IDENT_B: %v", err)
		return false
	}
	c, err := hmc.fastRegister8(HMC5883L_REG_IDENT_C)
	if err != nil {
		log.Error("Failed to read IDENT_C: %v", err)
		return false
	}
	if (a != 0x48) || (b != 0x34) || (c != 0x33) {
		return false
	}

	hmc.SetRange(HMC5883L_RANGE_1_3GA)
	hmc.SetMeasurementMode(HMC5883L_CONTINOUS)
	hmc.SetDataRate(HMC5883L_DATARATE_15HZ)
	hmc.SetSamples(HMC5883L_SAMPLES_1)

	hmc.mgPerDigit = 0.92

	return true
}

func (hmc HMC5883L) writeRegister8(reg uint8, value uint8) error {
	i2c, err := i2c.Bus(hmc.bus)
	if err != nil {
		return err
	}

	return i2c.WriteByte(HMC5883L_ADDRESS, reg, value)

}
func (hmc HMC5883L) readRegister8(reg uint8) (uint8, error) {
	i2c, err := i2c.Bus(hmc.bus)
	if err != nil {
		return 0, err
	}

	block, err := i2c.ReadByteBlock(HMC5883L_ADDRESS, reg, 1)
	if err != nil {
		return 0, err
	}
	return block[0], nil
}

func (hmc HMC5883L) fastRegister8(reg uint8) (uint8, error) {
	i2c, err := i2c.Bus(hmc.bus)
	if err != nil {
		return 0, err
	}

	block, err := i2c.ReadByteBlock(HMC5883L_ADDRESS, reg, 1)
	if err != nil {
		return 0, err
	}
	return block[0], nil
}

func (hmc HMC5883L) readRegister16(reg uint8) (int16, error) {
	i2c, err := i2c.Bus(hmc.bus)
	if err != nil {
		return 0, err
	}

	block, err := i2c.ReadByteBlock(HMC5883L_ADDRESS, reg, 2)
	if err != nil {
		return 0, err
	}
	return read_int16(block), nil
}

func read_int16(data []byte) (ret int16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}
