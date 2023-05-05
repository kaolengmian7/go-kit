package csvutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestGetByCoordinate(t *testing.T) {
	testCase := []struct {
		rowNum    int
		columnNum int
		input     [][]string
		expect    string
	}{
		{
			rowNum:    0,
			columnNum: 0,
			input:     [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}},
			expect:    "1",
		},
		{
			rowNum:    2,
			columnNum: 3,
			input:     [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}},
			expect:    "",
		},
		{
			rowNum:    3,
			columnNum: 2,
			input:     [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}},
			expect:    "",
		},
	}
	csv := New()
	for index, tc := range testCase {
		csv.data = tc.input
		output, err := csv.GetByCoordinate(tc.rowNum, tc.columnNum)
		assert.Equal(t, output, tc.expect)
		if index == 0 {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
			log.Printf("%+v", err)
		}
	}
}

func TestAddDataToRow(t *testing.T) {
	content := "test"
	testCase := []struct {
		rowNum int
		input  [][]string
		expect [][]string
	}{
		{
			rowNum: 0,
			input:  [][]string{{"1", "1", "1"}, {"2", "2", "2"}, {"3", "3", "3"}},
			expect: [][]string{{"1", "1", "1", content}, {"2", "2", "2"}, {"3", "3", "3"}},
		},
		{
			rowNum: 3,
			input:  [][]string{{"1", "1", "1"}, {"2", "2", "2"}, {"3", "3", "3"}},
			expect: [][]string{{"1", "1", "1"}, {"2", "2", "2"}, {"3", "3", "3"}},
		},
	}

	csv := New()
	for index, tc := range testCase {
		csv.data = tc.input
		err := csv.AddToRow(tc.rowNum, content)
		assert.ElementsMatch(t, csv.data, tc.expect)
		if index == 0 {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
			log.Printf("%+v", err)
		}
	}
}

func TestRead(t *testing.T) {
	csv := New()
	err := csv.Read("./test.csv")
	expect := [][]string{{"1", "2", "中文3"}, {"4", "5", "中文6"}, {"7", "8", "中文9"}}
	assert.Nil(t, err)
	assert.ElementsMatch(t, expect, csv.data)
}

func TestExport(t *testing.T) {
	csv := New()
	csv.data = [][]string{{"1", "2", "3", "test"}, {"4", "5", "6", ""}, {"7", "8", "9", ""}}
	err := csv.Export("test_export")
	assert.Nil(t, err)

	fp := fmt.Sprintf("%s-%s.csv", "test_export", time.Now().Format("2006-01-02"))
	log.Println(fp)
	// 验证文件存在
	_, err = os.Stat(fp)
	assert.Nil(t, err)
	// 验证内容一致
	expectCsv := New()
	err = expectCsv.Read(fp)
	expect := [][]string{{"\ufeff1", "2", "3", "test"}, {"4", "5", "6", ""}, {"7", "8", "9", ""}} // 因为写入了 UTF-8 BOM 所以要加上 \ufeff
	log.Printf("expectCSVData:%+v", expectCsv.data)
	assert.ElementsMatch(t, expect, expectCsv.data)
	// 删除测试文件
	err = os.Remove(fp)
	if err != nil {
		panic(fmt.Sprintf("csv_util2 TestExport failed! os.Remove(%s) failed!", fp))
	}
}
