package csvutil

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path"
	"time"
)

type CSV struct {
	data [][]string
}

func New() *CSV {
	c := &CSV{}
	c.data = make([][]string, 0)
	return c
}

// 根据坐标返回数据，rowNum == 0, column == 1 则返回坐标为(0,1)的数据
func (c *CSV) GetByCoordinate(rowNum int, columnNum int) (res string, err error) {
	if rowNum >= len(c.data) {
		return "", errors.WithStack(errors.New(fmt.Sprintf("rowNum(%d) invalid! sliceLen(%d)", rowNum, len(c.data))))
	}
	if columnNum >= len(c.data[rowNum]) {
		return "", errors.WithStack(errors.New(fmt.Sprintf("columnNum(%d) invalid! sliceLen(%d)", columnNum, len(c.data[rowNum]))))
	}

	return c.data[rowNum][columnNum], nil
}

// 在某一行末尾追加数据
func (c *CSV) AddToRow(rowNum int, data string) (err error) {
	if rowNum >= len(c.data) {
		return errors.WithStack(errors.New(fmt.Sprintf("rowNum(%d) invalid! sliceLen(%d)", rowNum, len(c.data))))
	}

	c.data[rowNum] = append(c.data[rowNum], data)
	return
}

// 读取 CSV 文件
func (c *CSV) Read(filePath string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	c.data, err = r.ReadAll()
	if err != nil {
		return
	}

	return
}

// 导出 CSV 文件到当前目录
func (c *CSV) Export(fileName string) (err error) {
	fileName = fmt.Sprintf("%s-%s.csv", fileName, time.Now().Format("2006-01-02"))

	savePath := "./"
	fp, err := os.Create(path.Join(savePath, fileName)) // 创建文件句柄
	if err != nil {
		return
	}
	defer fp.Close()

	fp.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(fp)
	w.WriteAll(c.data)
	w.Flush()

	return
}
