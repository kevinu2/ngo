package File

const (
	DefaultFileMode = uint32(0644)
	DefaultDirMode  = uint32(0755)
)

type FileMode uint8

const (
	ModeDefaultFile FileMode = iota + 1
	ModeDefaultDir
	ModeOwnerFile
	ModeOwnerDir
	ModeGroupFile
	ModeGroupDir
	ModeOtherFile
	ModeOtherDir
	ModeWriteOnly
	ModeFull
)

func (fm FileMode) Code() uint32 {
	switch fm {
	case ModeDefaultFile:
		return 0644
	case ModeDefaultDir:
		return 0755
	case ModeOwnerFile:
		return 0600
	case ModeOwnerDir:
		return 0755
	case ModeGroupFile:
		return 0660
	case ModeGroupDir:
		return 0770
	case ModeOtherFile:
		return 0666
	case ModeOtherDir:
		return 0777
	case ModeWriteOnly:
		return 0200
	case ModeFull:
		return 0777
	default:
		return 0777
	}
}
