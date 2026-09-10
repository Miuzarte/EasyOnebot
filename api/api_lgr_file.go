package api

/*
GetPrivateFileUrl 获取私聊文件资源链接

https://lagrange-onebot.apifox.cn/241662646e0

参数:

	userId: 用户 Uin，接收文件用户的Uin
	fileId: 文件 ID
	fileHash(可选): 文件 Hash
*/
func GetPrivateFileUrl(userId int, fileId string, fileHash ...string) *Request {
	if len(fileHash) == 0 {
		return NewReq("get_private_file_url", map[string]any{
			"user_id": userId,
			"file_id": fileId,
		})
	} else {
		return NewReq("get_private_file_url", map[string]any{
			"user_id":   userId,
			"file_id":   fileId,
			"file_hash": fileHash[0],
		})
	}
}

type validateGetPrivateFileUrl struct {
	UserId   int      `validate:"gt=0"`
	FileId   string   `validate:"required"`
	FileHash []string `validate:"omitempty"`
}

type GetPrivateFileUrlResp struct {
	Url string `json:"url" mapstructure:"url"` // 文件下载链接
}

func (c LgrCaller) GetPrivateFileUrl(userId int, fileId string, fileHash ...string) (*GetPrivateFileUrlResp, error) {
	err := validate.Struct(&validateGetPrivateFileUrl{userId, fileId, fileHash})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetPrivateFileUrlResp](c.PostReq(GetPrivateFileUrl(userId, fileId, fileHash...)))
}

/*
GetGroupFileUrl 获取群文件资源链接

https://lagrange-onebot.apifox.cn/236972687e0

参数:

	groupId: 群 Uin
	fileId: 文件 ID
	busid(已废弃):
*/
func GetGroupFileUrl(groupId int, fileId, busid string) *Request {
	if busid == "" {
		return NewReq("get_group_file_url", map[string]any{
			"group_id": groupId,
			"file_id":  fileId,
		})
	} else {
		return NewReq("get_group_file_url", map[string]any{
			"group_id": groupId,
			"file_id":  fileId,
			"busid":    busid,
		})
	}
}

type validateGetGroupFileUrl struct {
	GroupId int    `validate:"gt=0"`
	FileId  string `validate:"required"`
	Busid   string `validate:"omitempty"`
}

type GetGroupFileUrlResp struct {
	Url string `json:"url" mapstructure:"url"` // 文件链接
}

// GetGroupFileUrl 获取群文件资源链接
func (c LgrCaller) GetGroupFileUrl(groupId int, fileId, busid string) (*GetGroupFileUrlResp, error) {
	err := validate.Struct(&validateGetGroupFileUrl{groupId, fileId, busid})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetGroupFileUrlResp](c.PostReq(GetGroupFileUrl(groupId, fileId, busid)))
}

/*
GetGroupRootFiles 获取群根目录文件列表

https://lagrange-onebot.apifox.cn/236973502e0

参数:

	groupId: 群 Uin
*/
func GetGroupRootFiles(groupId int) *Request {
	return NewReq("get_group_root_files", map[string]any{
		"group_id": groupId,
	})
}

type GetGroupRootFilesResp struct {
	Files   []LgrFile   `json:"files" mapstructure:"files"`     // 文件列表
	Folders []LgrFolder `json:"folders" mapstructure:"folders"` // 文件夹列表
}

// LgrFile 拓展消息链
type LgrFile struct {
	GroupId       int    `json:"group_id" mapstructure:"group_id"`             // 群号
	FileId        string `json:"file_id" mapstructure:"file_id"`               // 文件ID
	FileName      string `json:"file_name" mapstructure:"file_name"`           // 文件名
	Busid         int    `json:"busid" mapstructure:"busid"`                   // 文件类型
	FileSize      int    `json:"file_size" mapstructure:"file_size"`           // 文件大小
	UploadTime    int    `json:"upload_time" mapstructure:"upload_time"`       // 上传时间
	DeadTime      int    `json:"dead_time" mapstructure:"dead_time"`           // 过期时间,永久文件恒为0
	ModifyTime    int    `json:"modify_time" mapstructure:"modify_time"`       // 最后修改时间
	DownloadTimes int    `json:"download_times" mapstructure:"download_times"` // 下载次数
	Uploader      int    `json:"uploader" mapstructure:"uploader"`             // 上传者ID
	UploaderName  string `json:"uploader_name" mapstructure:"uploader_name"`   // 上传者名字
}

// LgrFolder 拓展消息链
type LgrFolder struct {
	GroupId        int    `json:"group_id" mapstructure:"group_id"`                 // 群号
	FolderId       string `json:"folder_id" mapstructure:"folder_id"`               // 文件夹ID
	FolderName     string `json:"folder_name" mapstructure:"folder_name"`           // 文件名
	CreateTime     int    `json:"create_time" mapstructure:"create_time"`           // 创建时间
	Creator        int    `json:"creator" mapstructure:"creator"`                   // 创建者
	CreatorName    string `json:"creator_name" mapstructure:"creator_name"`         // 创建者名字
	TotalFileCount int    `json:"total_file_count" mapstructure:"total_file_count"` // 子文件数量
}

// GetGroupRootFiles 获取群根目录文件列表
func (c LgrCaller) GetGroupRootFiles(groupId int) (*GetGroupRootFilesResp, error) {
	err := validate.idGt0(groupId)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetGroupRootFilesResp](c.PostReq(GetGroupRootFiles(groupId)))
}

/*
GetGroupFilesByFolder 获取群子目录文件列表

https://lagrange-onebot.apifox.cn/236974042e0

参数:

	groupId: 群 Uin
	folderId: 文件夹 ID
*/
func GetGroupFilesByFolder(groupId int, folderId string) *Request {
	return NewReq("get_group_files_by_folder", map[string]any{
		"group_id": groupId,
		"folder":   folderId,
	})
}

type validateGetGroupFilesByFolder struct {
	GroupId  int    `validate:"gt=0"`
	FolderId string `validate:"required"`
}

type GetGroupFilesByFolderResp struct {
	Files   []LgrFile   `json:"files" mapstructure:"files"`     // 文件列表
	Folders []LgrFolder `json:"folders" mapstructure:"folders"` // 文件夹列表
}

// GetGroupFilesByFolder 获取群子目录文件列表
func (c LgrCaller) GetGroupFilesByFolder(groupId int, folderId string) (*GetGroupFilesByFolderResp, error) {
	err := validate.Struct(&validateGetGroupFilesByFolder{groupId, folderId})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetGroupFilesByFolderResp](c.PostReq(GetGroupFilesByFolder(groupId, folderId)))
}

/*
MoveGroupFile 移动群文件

https://lagrange-onebot.apifox.cn/236974078e0

参数:

	groupId: 群 Uin
	fildId: 文件 ID
	parentDir: 当前文件夹 ID
	targetDir: 目标文件夹 ID
*/
func MoveGroupFile(groupId int, fileId, parentDir, targetDir string) *Request {
	return NewReq("move_group_file", map[string]any{
		"group_id":         groupId,
		"file_id":          fileId,
		"parent_directory": parentDir,
		"target_directory": targetDir,
	})
}

type validateMoveGroupFile struct {
	GroupId   int    `validate:"gt=0"`
	FileId    string `validate:"required"`
	ParentDir string `validate:"required"`
	TargetDir string `validate:"required"`
}

func (c LgrCaller) MoveGroupFile(groupId int, fileId, parentDir, targetDir string) error {
	err := validate.Struct(&validateMoveGroupFile{groupId, fileId, parentDir, targetDir})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(MoveGroupFile(groupId, fileId, parentDir, targetDir))
}

/*
DeleteGroupFile 删除群文件

https://lagrange-onebot.apifox.cn/236974086e0

参数:

	groupId: 群 Uin
	fildId: 文件 ID
*/
func DeleteGroupFile(groupId int, fileId string) *Request {
	return NewReq("delete_group_file", map[string]any{
		"group_id": groupId,
		"file_id":  fileId,
	})
}

type validateDeleteGroupFile struct {
	GroupId int    `validate:"gt=0"`
	FileId  string `validate:"required"`
}

func (c LgrCaller) DeleteGroupFile(groupId int, fileId string) error {
	err := validate.Struct(&validateDeleteGroupFile{groupId, fileId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(DeleteGroupFile(groupId, fileId))
}

/*
CreateGroupFileFolder 创建群文件文件夹

https://lagrange-onebot.apifox.cn/236974237e0

参数:

	groupId: 群 Uin
	name: 文件夹名字
	parentId(已废弃): 父文件夹 ID，tx不再允许在非根目录创建文件夹了，该值废弃，请直接传递"/"
*/
func CreateGroupFileFolder(groupId int, fileId, parentId string) *Request {
	return NewReq("create_group_file_folder", map[string]any{
		"group_id": groupId,
		"file_id":  fileId,
		// "parant_id": parentId,
		"parant_id": "/",
	})
}

type validateCreateGroupFileFolder struct {
	GroupId  int    `validate:"gt=0"`
	FileId   string `validate:"required"`
	ParentId string `validate:"omitempty"`
}

func (c LgrCaller) CreateGroupFileFolder(groupId int, fileId, parentId string) error {
	err := validate.Struct(&validateCreateGroupFileFolder{groupId, fileId, parentId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(CreateGroupFileFolder(groupId, fileId, parentId))
}

/*
DeleteGroupFileFolder 删除群文件文件夹

https://lagrange-onebot.apifox.cn/236974248e0

参数:

	groupId: 群 Uin
	folderId: 文件夹 ID
*/
func DeleteGroupFileFolder(groupId int, folderId string) *Request {
	return NewReq("delete_group_file_folder", map[string]any{
		"group_id":  groupId,
		"folder_id": folderId,
	})
}

type validateDeleteGroupFileFolder struct {
	GroupId  int    `validate:"gt=0"`
	FolderId string `validate:"required"`
}

func (c LgrCaller) DeleteGroupFileFolder(groupId int, folderId string) error {
	err := validate.Struct(&validateDeleteGroupFileFolder{groupId, folderId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(DeleteGroupFileFolder(groupId, folderId))
}

/*
RenameGroupFileFolder 重命名群文件文件夹名

https://lagrange-onebot.apifox.cn/236974260e0

参数:

	groupId: 群 Uin
	folderId: 文件夹 ID
	newFolderName: 新文件夹名称
*/
func RenameGroupFileFolder(groupId int, folderId, newFolderName string) *Request {
	return NewReq("rename_group_file_folder", map[string]any{
		"group_id":        groupId,
		"folder_id":       folderId,
		"new_folder_name": newFolderName,
	})
}

type validateRenameGroupFileFolder struct {
	GroupId       int    `validate:"gt=0"`
	FolderId      string `validate:"required"`
	NewFolderName string `validate:"required"`
}

func (c LgrCaller) RenameGroupFileFolder(groupId int, folderId, newFolderName string) error {
	err := validate.Struct(&validateRenameGroupFileFolder{groupId, folderId, newFolderName})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(RenameGroupFileFolder(groupId, folderId, newFolderName))
}

/*
UploadGroupFile 上传群文件

https://lagrange-onebot.apifox.cn/236974303e0

参数:

	groupId: 群 Uin
	file: file 链接, 仅支持本地Path
	name: 文件名称
	folder: 文件夹 ID
*/
func UploadGroupFile(groupId int, file, name, folder string) *Request {
	return NewReq("upload_group_file", map[string]any{
		"group_id": groupId,
		"file":     file,
		"name":     name,
		"folder":   folder,
	})
}

type validateUploadGroupFile struct {
	GroupId int    `validate:"gt=0"`
	File    string `validate:"required"`
	Name    string `validate:"required"`
	Folder  string `validate:"required"`
}

// UploadGroupFile 上传群文件
func (c LgrCaller) UploadGroupFile(groupId int, file, name, folder string) error {
	err := validate.Struct(&validateUploadGroupFile{groupId, file, name, folder})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(UploadGroupFile(groupId, file, name, folder))
}

/*
UploadPrivateFile 私聊发送文件

https://lagrange-onebot.apifox.cn/236974322e0

参数:

	userId: 用户 Uin
	file: file 链接, 仅支持本地Path
	name: 文件名称
*/
func UploadPrivateFile(userId int, file, name string) *Request {
	return NewReq("upload_private_file", map[string]any{
		"user_id": userId,
		"file":    file,
		"name":    name,
	})
}

type validateUploadPrivateFile struct {
	UserId int    `validate:"gt=0"`
	File   string `validate:"required"`
	Name   string `validate:"required"`
}

// UploadPrivateFile 私聊发送文件
func (c LgrCaller) UploadPrivateFile(userId int, file, name string) error {
	err := validate.Struct(&validateUploadPrivateFile{userId, file, name})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(UploadPrivateFile(userId, file, name))
}
