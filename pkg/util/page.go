package util

import "strconv"

// Page Page
type Page struct {
	PageSize  uint64      `json:"page_size"`
	PageIndex uint64      `json:"page_index"`
	Total     uint64      `json:"total"`
	Data      interface{} `json:"data"`
}

// NewPage NewPage
func NewPage(pageIndex uint64, pageSize uint64, total uint64, data interface{}) *Page {

	return &Page{PageIndex: pageIndex, PageSize: pageSize, Data: data, Total: total}
}

//ToPageNumOrDefault 将字符串转换为数字类型 如果字符串为空或无效 则赋值分页默认参数
func ToPageNumOrDefault(pageIndex string, pageSize string) (pIndex64 uint64, pSize64 uint64) {
	var pageIndex64 uint64
	var pageSize64 uint64

	pageIndex64, err := strconv.ParseUint(pageIndex, 10, 64)
	if err != nil || pageIndex64 == 0 {
		pageIndex64 = 1
	}

	pageSize64, err = strconv.ParseUint(pageSize, 10, 64)
	if err != nil || pageSize64 == 0 {
		pageSize64 = 10
	}

	return pageIndex64, pageSize64
}
