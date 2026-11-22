package models

type ReplayFrame interface {
	ButtonState() uint8
	MouseX() float32
	MouseY() float32
	Time() int32
}

type ReplayFrameBundle interface {
	Action() uint8
	Extra() int32
	Sequence() uint16
	FrameCount() int
	Frame(index int) ReplayFrame

	ScoreFrame() (frame ScoreFrame, present bool)
}

type ScoreFrame interface {
	Time() int32
	Id() uint8
	Total300() uint16
	Total100() uint16
	Total50() uint16
	TotalGeki() uint16
	TotalKatu() uint16
	TotalMiss() uint16
	TotalScore() uint32
	MaxCombo() uint16
	CurrentCombo() uint16
	Perfect() bool
	Hp() uint8
	TagByte() uint8
	UsingScoreV2() bool
	ComboPortion() float64
	BonusPortion() float64
}
