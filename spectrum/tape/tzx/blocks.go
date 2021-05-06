package tzx

import "github.com/voytas/z80-go-zx/spectrum/helpers"

type HeaderDataBlock struct {
	signature string // TZX signature
	eot       byte   // End of text file marker
	VerMajor  byte   // TZX major revision number
	VerMinor  byte   // TZX minor revision number
}

type StandardSpeedDataBlock struct {
	PauseAfter uint16 // Pause after this block (ms.) {1000}
	Length     uint16 // Length of data that follow
	Data       []byte // Data as in .TAP files
}

type TurboSpeedDataBlock struct {
	PilotPulseLength   uint16 // Length of PILOT pulse {2168}
	SyncPulse1Length   uint16 // Length of SYNC first pulse {667}
	SyncPulse2Length   uint16 // Length of SYNC second pulse {735}
	ZeroBitPulseLength uint16 // Length of ZERO bit pulse {855}
	OneBitPulseLength  uint16 // Length of ONE bit pulse {1710}
	PilotToneLength    uint16 // Length of PILOT tone (number of pulses) {8063 header (flag<128), 3223 data (flag>=128)}
	UsedBitsLastByte   byte   // Used bits in the last byte (other bits should be 0) {8}
	PauseAfter         uint16 // Pause after this block (ms.) {1000}
	Length             []byte // Length of data that follow
	Data               []byte // Data as in .TAP files
}

type PureToneDataBlock struct {
	PulseLength uint16 // Length of one pulse in T-states
	PulseCount  uint16 // Number of pulses
}

type PulseSequenceDataBlock struct {
	PulseCount   byte     // Number of pulses
	PulseLengths []uint16 // Pulses' lengths
}

type PureDataBlock struct {
	ZeroBitPulseLength uint16 // Length of ZERO bit pulse
	OneBitPulseLength  uint16 // Length of ONE bit pulse
	UsedBitsLastByte   byte   // Used bits in last byte (other bits should be 0)
	PauseAfter         uint16 // Pause after this block (ms.)
	Length             []byte // Length of data that follow
	Data               []byte // Data as in .TAP files
}

type DirectRecordingDataBlock struct {
	TStatesPerSample uint16 // Number of T-states per sample (bit of data)
	PauseAfter       uint16 // Pause after this block in milliseconds (ms.)
	UsedBitsLastByte byte   // Used bits (samples) in last byte of data (1-8)
	Length           []byte // Length of samples' data
	Data             []byte // Samples data. Each bit represents a state on the EAR port (i.e. one sample)
}

type CswRecordingDataBlock struct {
	Length          uint32 // Block length (without these four bytes)
	PauseAfter      uint16 // Pause after this block (in ms).
	SamplingRate    []byte // Sampling rate
	CompressionType byte   // Compression type
	PulseCount      uint32 // Number of stored pulses (after decompression, for validation purposes)
	Data            []byte // CSW data, encoded according to the CSW file format specification
}

type GeneralizedDataBlock struct {
	Length                 uint32       // Block length (without these four bytes)
	PauseAfter             uint16       // Pause after this block (ms)
	PilotTotalSymbols      uint32       // TOTP - Total number of symbols in pilot/sync block (can be 0)
	PilotMaxPulses         byte         // NPP - Maximum number of pulses per pilot/sync symbol
	PilotSymbolsCountAlpha byte         // ASP - Number of pilot/sync symbols in the alphabet table (0=256)
	DataTotalSymbols       uint32       // TOTD - Total number of symbols in data stream (can be 0)
	DataMaxPulses          byte         // NPD - Maximum number of pulses per data symbol
	DataSymbolsCountAlpha  byte         // ASD - Number of data symbols in the alphabet table (0=256)
	PilotSymbols           []*SymbolDef // Pilot and sync symbols definition table
	PilotDataStream        []*Prle      // Pilot and sync data stream
	DataSymbols            []*SymbolDef // Data symbols definition table
	DataStream             []byte       // Data stream
}

type SymbolDef struct {
	Flags        byte     // Symbol flags
	PulseLengths []uint16 // Array of pulse lengths
}

type Prle struct {
	Symbol      byte   // Symbol to be represented
	RepeatCount uint16 // Number of repetitions
}

type SilenceDataBlock struct {
	PauseDuration uint16 // Pause duration (ms.)
}

type GroupStartDataBlock struct {
	Length byte   // Length of the group name string
	Chars  []byte // Group name in ASCII format
}

type GroupEndDataBlock struct {
}

type JumpToDataBlock struct {
	Jump uint16 // Relative jump value
}

type LoopStartDataBlock struct {
	Count uint16 // Number of repetitions
}

type LoopEndDataBlock struct {
}

type CallSequenceDataBlock struct {
	Count   uint16   // Number of calls to be made
	Offsets []uint16 // Array of call block numbers (relative-signed offsets)
}

type ReturnFromSequenceDataBlock struct {
}

type SelectDataBlock struct {
	Length     uint16    // Length of the whole block (without these two bytes)
	Count      byte      // Number of selections
	Selections []*Select // List of selections
}

type Select struct {
	Offset uint16 // Relative Offset
	Length byte   // Length of description text
	Chars  []byte // Description text
}

type StopTheTape48kDataBlock struct {
	Length uint16 // Length of the block without these four bytes (0)
}

type SetSignalLevelDataBlock struct {
	Length uint16 // lock length (without these four bytes)
	Level  byte   // Signal level (0=low, 1=high)
}

type TextDescriptionDataBlock struct {
	Length      byte   // Length of the text description
	Description []byte // Text description in ASCII format
}

type MessageDataBlock struct {
	Time    byte   // Time (in seconds) for which the message should be displayed
	Length  byte   // Length of the text message
	Message []byte // Message that should be displayed in ASCII format
}

type ArchiveInfoDataBlock struct {
	Length uint16  // Length of the whole block (without these two bytes)
	Count  byte    // Number of text strings
	Texts  []*Text // List of text strings
}

type Text struct {
	Id     byte   // Text identification byte
	Length byte   // Length of text string
	Text   []byte // Text string in ASCII format
}

type HardwareTypeDataBlock struct {
	Count byte            // Number of machines and hardware types for which info is supplied
	Infos []*HardwareInfo // List of machines and hardware
}

type HardwareInfo struct {
	Type byte // Hardware type
	Id   byte // Hardware ID
	Info byte // Hardware information
}

type CustomInfoDataBlock struct {
	IdText     []byte // Identification string (in ASCII)
	Length     uint16 // Length of the custom info
	CustomInfo []byte // Custom info
}

type GlueDataBlock struct {
	Value []byte // Value: { "XTape!",0x1A,MajR,MinR }
}

func (b *StandardSpeedDataBlock) populate(r *helpers.BinaryReader) {
	b.PauseAfter = r.ReadWord()
	b.Length = r.ReadWord()
	b.Data = r.ReadBytes(int(b.Length))
}

func (b *TurboSpeedDataBlock) populate(r *helpers.BinaryReader) {
	b.PilotPulseLength = r.ReadWord()
	b.SyncPulse1Length = r.ReadWord()
	b.SyncPulse2Length = r.ReadWord()
	b.ZeroBitPulseLength = r.ReadWord()
	b.OneBitPulseLength = r.ReadWord()
	b.PilotToneLength = r.ReadWord()
	b.UsedBitsLastByte = r.ReadByte()
	b.PauseAfter = r.ReadWord()
	b.Length = r.ReadBytes(3)
	b.Data = r.ReadBytes(bytesToInt(b.Length))
}

func (b *PureToneDataBlock) populate(r *helpers.BinaryReader) {
	b.PulseLength = r.ReadWord()
	b.PulseCount = r.ReadWord()
}

func (b *PulseSequenceDataBlock) populate(r *helpers.BinaryReader) {
	b.PulseCount = r.ReadByte()
	b.PulseLengths = r.ReadWords(int(b.PulseCount))
}

func (b *PureDataBlock) populate(r *helpers.BinaryReader) {
	b.ZeroBitPulseLength = r.ReadWord()
	b.OneBitPulseLength = r.ReadWord()
	b.UsedBitsLastByte = r.ReadByte()
	b.PauseAfter = r.ReadWord()
	b.Length = r.ReadBytes(3)
	b.Data = r.ReadBytes(bytesToInt(b.Length))
}

func (b *DirectRecordingDataBlock) populate(r *helpers.BinaryReader) {
	b.TStatesPerSample = r.ReadWord()
	b.PauseAfter = r.ReadWord()
	b.UsedBitsLastByte = r.ReadByte()
	b.Length = r.ReadBytes(3)
	b.Data = r.ReadBytes(bytesToInt(b.Length))
}

func (b *CswRecordingDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadDWord()
	b.PauseAfter = r.ReadWord()
	b.SamplingRate = r.ReadBytes(3)
	b.CompressionType = r.ReadByte()
	b.PulseCount = r.ReadDWord()
	b.Data = r.ReadBytes(int(b.Length - 2 - 3 - 1 - 4))
}

func (b *GeneralizedDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadDWord()
	b.PauseAfter = r.ReadWord()
	b.PilotTotalSymbols = r.ReadDWord()
	b.PilotMaxPulses = r.ReadByte()
	b.PilotSymbolsCountAlpha = r.ReadByte()
	b.DataTotalSymbols = r.ReadDWord()
	b.DataMaxPulses = r.ReadByte()
	b.DataSymbolsCountAlpha = r.ReadByte()
	if b.PilotTotalSymbols > 0 {
		b.PilotSymbols = make([]*SymbolDef, int(b.PilotSymbolsCountAlpha))
		for i := 0; i < len(b.PilotSymbols); i++ {
			b.PilotSymbols[i] = &SymbolDef{
				Flags:        r.ReadByte(),
				PulseLengths: r.ReadWords(int(b.PilotMaxPulses)),
			}
		}
		b.PilotDataStream = make([]*Prle, int(b.PilotTotalSymbols))
		for i := 0; i < len(b.PilotDataStream); i++ {
			b.PilotDataStream[i] = &Prle{
				Symbol:      r.ReadByte(),
				RepeatCount: r.ReadWord(),
			}
		}
	}
	if b.DataTotalSymbols > 0 {
		b.DataSymbols = make([]*SymbolDef, int(b.DataSymbolsCountAlpha))
		for i := 0; i < len(b.DataSymbols); i++ {
			b.DataSymbols[i] = &SymbolDef{
				Flags:        r.ReadByte(),
				PulseLengths: r.ReadWords(int(b.DataMaxPulses)),
			}
		}
		b.DataStream = r.ReadBytes(int(b.DataTotalSymbols / 8))
	}
}

func (b *SilenceDataBlock) populate(r *helpers.BinaryReader) {
	b.PauseDuration = r.ReadWord()
}

func (b *GroupStartDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadByte()
	b.Chars = r.ReadBytes(int(b.Length))
}

func (b *GroupEndDataBlock) populate(r *helpers.BinaryReader) {
}

func (b *JumpToDataBlock) populate(r *helpers.BinaryReader) {
	b.Jump = r.ReadWord()
}

func (b *LoopStartDataBlock) populate(r *helpers.BinaryReader) {
	b.Count = r.ReadWord()
}

func (b *LoopEndDataBlock) populate(r *helpers.BinaryReader) {
}

func (b *CallSequenceDataBlock) populate(r *helpers.BinaryReader) {
	b.Count = r.ReadWord()
	b.Offsets = r.ReadWords(int(b.Count))
}

func (b *ReturnFromSequenceDataBlock) populate(r *helpers.BinaryReader) {
}

func (b *SelectDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadWord()
	b.Count = r.ReadByte()
	b.Selections = make([]*Select, b.Count)
	for i := 0; i < len(b.Selections); i++ {
		s := &Select{
			Offset: r.ReadWord(),
		}
		s.Length = r.ReadByte()
		s.Chars = r.ReadBytes(int(b.Length))
		b.Selections[i] = s
	}
}

func (b *StopTheTape48kDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadWord()
}

func (b *SetSignalLevelDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadWord()
	b.Level = r.ReadByte()
}

func (b *TextDescriptionDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadByte()
	b.Description = r.ReadBytes(int(b.Length))
}

func (b *MessageDataBlock) populate(r *helpers.BinaryReader) {
	b.Time = r.ReadByte()
	b.Length = r.ReadByte()
	b.Message = r.ReadBytes(int(b.Length))
}

func (b *ArchiveInfoDataBlock) populate(r *helpers.BinaryReader) {
	b.Length = r.ReadWord()
	b.Count = r.ReadByte()
	b.Texts = make([]*Text, b.Count)
	for i := 0; i < len(b.Texts); i++ {
		t := &Text{
			Id: r.ReadByte(),
		}
		t.Length = r.ReadByte()
		t.Text = r.ReadBytes(int(t.Length))
		b.Texts[i] = t
	}
}

func (b *HardwareTypeDataBlock) populate(r *helpers.BinaryReader) {
	b.Count = r.ReadByte()
	b.Infos = make([]*HardwareInfo, b.Count)
	for i := 0; i < len(b.Infos); i++ {
		b.Infos[i] = &HardwareInfo{
			Type: r.ReadByte(),
			Id:   r.ReadByte(),
			Info: r.ReadByte(),
		}
	}
}

func (b *CustomInfoDataBlock) populate(r *helpers.BinaryReader) {
	b.IdText = r.ReadBytes(10)
	b.Length = r.ReadWord()
	b.CustomInfo = r.ReadBytes(int(b.Length))
}

func (b *GlueDataBlock) populate(r *helpers.BinaryReader) {
	b.Value = r.ReadBytes(9)
}
