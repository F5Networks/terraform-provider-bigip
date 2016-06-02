package bigip

import ()

var (
	uriUpload = "shared/file-transfer/uploads"
)

//func (b *BigIP) UploadFile(name string, data []byte) error  {
//	req := &APIRequest{
//		Method:      "post",
//		URL:         fmt.Sprintf("%s/%s", uriUpload, name),
//		Body:        data,
//		ContentType: "application/octet-stream",
//	}
//	_, callErr := b.APICall(req)
//	return callErr
//}
