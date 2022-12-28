package model

func GetAllBytlies() ([]Bytly, error) {
	var bytlies []Bytly

	result := db.Order("ID").Find(&bytlies)

	if result.Error != nil {
		return []Bytly{}, result.Error
	}

	return bytlies, nil
}

func GetBytlyById(id uint64) (Bytly, error) {
	var bytly Bytly

	result := db.Where("id = ?", id).First(&bytly)

	if result.Error != nil {
		return Bytly{}, result.Error
	}

	return bytly, nil
}

func GetBytlyByShortcut(shortcut string) (Bytly, error) {
	var bytly Bytly

	result := db.Where("bytly = ?", shortcut).First(&bytly)

	if result.Error != nil {
		return Bytly{}, result.Error
	}

	return bytly, nil
}

func CreateBytly(bytly Bytly) (Bytly, error) {
	result := db.Create(&bytly)

	if result.Error != nil {
		return Bytly{}, result.Error
	}

	return bytly, nil
}

func DeleteBytly(id uint64) error {
	result := db.Delete(&Bytly{}, id)

	return result.Error
}

func UpdateBytly(bytly Bytly) (Bytly, error) {
	result := db.Updates(bytly)

	if result.Error != nil {
		return Bytly{}, result.Error
	}

	return bytly, nil
}
