package sound

// Code based on mame implementation, I wouldn't know how to do it otherwise:
// https://github.com/mamedev/mame/blob/master/src/devices/sound/ay8910.cpp
// https://github.com/jsanchezv/JSpeccy/blob/master/src/machine/AY8912.java
// https://github.com/gasman/jsspeccy2/blob/master/core/sound.js
// http://f.rdw.se/AY-3-8910-datasheet.pdf
// http://map.grauw.nl/resources/sound/generalinstrument_ay-3-8910.pdf

// AY registers
const (
	regFineA      = 0x00 // Channel A fine tune
	regCoarseA    = 0x01 // Channel A coarse tune
	regFineB      = 0x02 // Channel B fine tune
	regCoarseB    = 0x03 // Channel B coarse tune
	regFineC      = 0x04 // Channel C fine tune
	regCoarseC    = 0x05 // Channel C coarse tune
	regNoise      = 0x06 // Noise
	regEnable     = 0x07 // Multi-function enable
	regAmplitudeA = 0x08 // Channel A amplitude
	regAmplitudeB = 0x09 // Channel B amplitude
	regAmplitudeC = 0x0A // Channel C amplitude
	regFineE      = 0x0B // Envelope fine tune
	regCoarseE    = 0x0C // Envelope coarse tune
	regShapeE     = 0x0D // Envelope shape
	regPortA      = 0x0E // I/O port A
	regPortB      = 0x0F // I/O port B

	channelA        = 0
	channelB        = 1
	channelC        = 2
	numChannels     = 3
	maxAmplitude    = 10900
	ampModeFixed    = 0
	ampModeVariable = 1
)

var volRates = [16]float32{
	0.0000, 0.0137, 0.0205, 0.0291, 0.0423, 0.0618, 0.0847, 0.1369,
	0.1691, 0.2647, 0.3527, 0.4499, 0.5704, 0.6873, 0.8482, 1.0000,
}
var volLevels [16]float32

type ayCommand struct {
}

type AY8910 struct {
	reg      byte     // currently selected register
	regs     [16]byte // all registers
	tones    [numChannels]tone
	envelope envelope
}

func NewAY8910() *AY8910 {
	ay := &AY8910{}
	for i, v := range volRates {
		volLevels[i] = v * maxAmplitude
	}
	return ay
}

func (ay *AY8910) writeReg(val byte) {
	ay.regs[ay.reg] = val

	switch ay.reg {
	case regFineA, regCoarseA:
		ay.tones[channelA].setPeriod(ay.regs[regFineA], ay.regs[regCoarseA])
	case regFineB, regCoarseB:
		ay.tones[channelB].setPeriod(ay.regs[regFineB], ay.regs[regCoarseB])
	case regFineC, regCoarseC:
		ay.tones[channelC].setPeriod(ay.regs[regFineC], ay.regs[regCoarseC])
	case regNoise:
		ay.regs[regNoise] &= 0x1F
	case regEnable:
		ay.tones[channelA].toneEnabled = val&0x01 == 0
		ay.tones[channelA].noiseEnabled = val&0x08 == 0
		ay.tones[channelB].toneEnabled = val&0x02 == 0
		ay.tones[channelB].noiseEnabled = val&0x10 == 0
		ay.tones[channelC].toneEnabled = val&0x04 == 0
		ay.tones[channelC].noiseEnabled = val&0x20 == 0
	case regAmplitudeA:
		ay.tones[channelA].envEnabled = val&0x10 != 0
		if ay.tones[channelA].envEnabled {
			ay.tones[channelA].amplitude = volLevels[ay.envelope.amplitude]
		} else {
			ay.tones[channelA].amplitude = volLevels[val&0x0F]
		}
	case regAmplitudeB:
		ay.tones[channelB].envEnabled = val&0x10 != 0
	case regAmplitudeC:
		ay.tones[channelC].envEnabled = val&0x10 != 0
	case regFineE, regCoarseE:
		ay.envelope.period = uint16(ay.regs[regCoarseE])<<8 | uint16(ay.regs[regFineE])
		if ay.envelope.period == 0 {
			ay.envelope.period = 2
		} else {
			ay.envelope.period <<= 1
		}
	case regShapeE:
		ay.envelope.hold = val&0x01 != 0
		ay.envelope.alternate = val&0x02 != 0
		ay.envelope.attack = val&0x04 != 0
		ay.envelope.cont = val&0x08 != 0
		if ay.envelope.attack {
			ay.envelope.amplitude = 0
		} else {
			ay.envelope.amplitude = 15
		}

		if ay.tones[channelA].envEnabled {
			ay.tones[channelA].amplitude = 0
		}
		if ay.tones[channelB].envEnabled {
			ay.tones[channelB].amplitude = 0
		}
		if ay.tones[channelC].envEnabled {
			ay.tones[channelC].amplitude = 0
		}
	}
}

func (ay *AY8910) SelectReg(reg byte) {
	ay.reg = reg
}

func (ay *AY8910) WriteReg(val byte, t int64) {

}

func Update(t int64) {

}

func StartPlay() {

}

type tone struct {
	period       uint16
	toneEnabled  bool
	noiseEnabled bool
	envEnabled   bool
	amplitude    float32
}

type envelope struct {
	amplitude byte
	period    uint16
	hold      bool
	alternate bool
	attack    bool
	cont      bool
}

func (t *tone) setPeriod(fine, coarse byte) {
	t.period = uint16(fine) | uint16(coarse&0x0F)<<8
}
