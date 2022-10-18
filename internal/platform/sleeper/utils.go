package sleeper

type ByUserID []interface {
	UserID() userID
}

func (us ByUserID) Len() int           { return len(us) }
func (us ByUserID) Less(i, j int) bool { return us[i].UserID() < us[j].UserID() }
func (us ByUserID) Swap(i, j int)      { us[i], us[j] = us[j], us[i] }
