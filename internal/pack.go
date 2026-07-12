package internal

// TogglePack switches ON/OFF state.
//
// SYSTEM and ERROR packs are ignored.
func TogglePack(pack *Pack) bool {

	switch pack.Status {

	case StatusOn:
		pack.Status = StatusOff
		return true

	case StatusOff:
		pack.Status = StatusOn
		return true

	default:
		return false
	}
}
