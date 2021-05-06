package tzx

import (
	"fmt"
	"io/ioutil"

	"github.com/voytas/z80-go-zx/spectrum/helpers"
)

const (
	signature = "ZXTape!"
	eot       = 0x1A

	StandardSpeedDataBlockId      = 0x10
	TurboSpeedDataBlockId         = 0x11
	PureToneDataBlockId           = 0x12
	PulseSequenceDataBlockId      = 0x13
	PureDataBlockId               = 0x14
	DirectRecordingDataBlockId    = 0x15
	CswRecordingDataBlockId       = 0x18
	GeneralizedDataBlockId        = 0x19
	SilenceDataBlockId            = 0x20
	GroupStartDataBlockId         = 0x21
	GroupEndDataBlockId           = 0x22
	JumpToDataBlockId             = 0x23
	LoopStartDataBlockId          = 0x24
	LoopEndDataBlockId            = 0x25
	CallSequenceDataBlockId       = 0x26
	ReturnFromSequenceDataBlockId = 0x27
	SelectDataBlockId             = 0x28
	StopTheTape48kDataBlockId     = 0x2A
	SetSignalLevelDataBlockId     = 0x2B
	TextDescriptionDataBlockId    = 0x30
	MessageDataBlockId            = 0x31
	ArchiveInfoDataBlockId        = 0x32
	HardwareTypeDataBlockId       = 0x33
	CustomInfoDataBlockId         = 0x35
	GlueDataBlockId               = 0x5A
)

type Tzx struct {
	Header HeaderDataBlock
	Blocks []Block
}

type Block interface {
	populate(r *helpers.BinaryReader)
}

func Load(file string) (*Tzx, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	r := helpers.NewBinaryReader(data)
	header := readHeader(r)
	if header == nil || header.signature != signature || header.eot != eot {
		return nil, fmt.Errorf("'%s' is not a valid TZX file.", file)
	}

	tzx := &Tzx{
		Header: *header,
	}

	for {
		block := readBlock(r)
		if block == nil {
			break
		}
		tzx.Blocks = append(tzx.Blocks, block)
	}

	return tzx, nil
}

func readHeader(r *helpers.BinaryReader) *HeaderDataBlock {
	return &HeaderDataBlock{
		signature: r.ReadString(7),
		eot:       r.ReadByte(),
		VerMajor:  r.ReadByte(),
		VerMinor:  r.ReadByte(),
	}
}

func readBlock(r *helpers.BinaryReader) Block {
	var b Block
	id := r.ReadByte()
	switch id {
	case StandardSpeedDataBlockId:
		b = &StandardSpeedDataBlock{}
	case TurboSpeedDataBlockId:
		b = &TurboSpeedDataBlock{}
	case PureToneDataBlockId:
		b = &PureToneDataBlock{}
	case PulseSequenceDataBlockId:
		b = &PulseSequenceDataBlock{}
	case PureDataBlockId:
		b = &PureDataBlock{}
	case DirectRecordingDataBlockId:
		b = &DirectRecordingDataBlock{}
	case CswRecordingDataBlockId:
		b = &CswRecordingDataBlock{}
	case GeneralizedDataBlockId:
		b = &GeneralizedDataBlock{}
	case SilenceDataBlockId:
		b = &SilenceDataBlock{}
	case GroupStartDataBlockId:
		b = &GroupStartDataBlock{}
	case GroupEndDataBlockId:
		b = &GroupEndDataBlock{}
	case JumpToDataBlockId:
		b = &JumpToDataBlock{}
	case LoopStartDataBlockId:
		b = &LoopStartDataBlock{}
	case LoopEndDataBlockId:
		b = &LoopEndDataBlock{}
	case CallSequenceDataBlockId:
		b = &CallSequenceDataBlock{}
	case ReturnFromSequenceDataBlockId:
		b = &ReturnFromSequenceDataBlock{}
	case SelectDataBlockId:
		b = &SelectDataBlock{}
	case StopTheTape48kDataBlockId:
		b = &StopTheTape48kDataBlock{}
	case SetSignalLevelDataBlockId:
		b = &SetSignalLevelDataBlock{}
	case TextDescriptionDataBlockId:
		b = &TextDescriptionDataBlock{}
	case MessageDataBlockId:
		b = &MessageDataBlock{}
	case ArchiveInfoDataBlockId:
		b = &ArchiveInfoDataBlock{}
	case HardwareTypeDataBlockId:
		b = &HardwareTypeDataBlock{}
	case CustomInfoDataBlockId:
		b = &CustomInfoDataBlock{}
	case GlueDataBlockId:
		b = &GlueDataBlock{}
	}

	if b != nil {
		b.populate(r)
	}

	return b
}
