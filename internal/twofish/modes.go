package twofish

import "errors"

type Mode int

const (
	ECB Mode = iota
	CBC
	CFB
	OFB
)

func (m Mode) String() string {
	switch m {
	case ECB:
		return "ECB"
	case CBC:
		return "CBC"
	case CFB:
		return "CFB"
	case OFB:
		return "OFB"
	}

	return ""
}

func ParseMode(mode string) (Mode, error) {
	switch mode {
	case "ECB":
		return ECB, nil
	case "CBC":
		return CBC, nil
	case "CFB":
		return CFB, nil
	case "OFB":
		return OFB, nil
	}

	return 0, errors.New("unsupported mode")
}
