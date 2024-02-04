package t

type SettingFolderValidation struct {
	Folder string `form:"Folder" validate:"required,dirpath"`
}

type SettingFolderDeleteValidation struct {
	FolderId uint `form:"FolderId" validate:"required,number,gte=1"`
}
