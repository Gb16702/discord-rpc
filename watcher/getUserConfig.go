package watcher

func GetUserConfig(item []Item) *Item {
	for _, i := range item {
		if(i.IsDefault) {
			return &i
		}
	}

	return nil
}
