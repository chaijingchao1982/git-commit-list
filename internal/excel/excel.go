package excel

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

// RemoveFile 删除指定的文件，若文件不存在则忽略。
func RemoveFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("未指定文件名")
	}
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件 %s 失败: %w", filename, err)
	}
	return nil
}

// ReadLastCommitHash 从 Excel 最后一行的 B 列读取上次的 commit hash。
// 返回空字符串表示表格为空或缺少 B 列数据。
func ReadLastCommitHash(filename, sheet string) (string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return "", fmt.Errorf("打开文件 %s 失败: %w", filename, err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return "", fmt.Errorf("读取工作表 %s 失败: %w", sheet, err)
	}

	if len(rows) == 0 {
		return "", nil
	}

	lastRow := rows[len(rows)-1]
	if len(lastRow) >= 2 {
		return lastRow[1], nil
	}

	return "", nil
}

// AppendRow 向指定工作表追加一行，A 列为 subject，B 列为 commitHash。
func AppendRow(f *excelize.File, sheet, subject, commitHash string) error {
	if f == nil {
		return fmt.Errorf("excelize.File 为 nil")
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return fmt.Errorf("读取工作表 %s 失败: %w", sheet, err)
	}

	rowNum := len(rows) + 1
	f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), subject)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), commitHash)

	if err := f.Save(); err != nil {
		return fmt.Errorf("保存文件失败：%w", err)
	}

	return nil
}
