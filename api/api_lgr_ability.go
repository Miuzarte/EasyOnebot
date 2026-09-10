package api

// func CanSendImage() *Request
// std: [CanSendImage] https://lagrange-onebot.apifox.cn/236970338e0
// func CanSendRecord() *Request
// std: [CanSendRecord] https://lagrange-onebot.apifox.cn/236970561e0

/*
UploadImage 上传图片

https://lagrange-onebot.apifox.cn/236970730e0

参数:

	file: file 链接, 支持 http/https/file/base64
*/
func UploadImage(file any) *Request {
	return NewReq("upload_image", map[string]any{
		"file": file,
	})
}

type UploadImageResp = string // 文件 Url

func (c LgrCaller) UploadImage(file any) (UploadImageResp, error) {
	err := validate.require(file)
	if err != nil {
		return "", err
	}
	return DecodeResponseAssertion[UploadImageResp](c.PostReq(UploadImage(file)))
}
