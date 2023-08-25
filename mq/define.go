package mq

// TransferData 转移队列中小型载体的结构格式
type TransferData struct {
	FileHash     string
	CurLocation  string
	DestLocation string
	//DestStoreType common.StoreType  // 用于扩展
}
