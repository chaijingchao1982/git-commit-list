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

const metaSheet = "Meta"

// ReadCheckpointHash 从 Meta 工作表的 A1 单元格读取 checkpoint commit hash。
// 若 Meta sheet 不存在或为空则返回空字符串。
func ReadCheckpointHash(filename string) (string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return "", fmt.Errorf("打开文件 %s 失败: %w", filename, err)
	}
	defer f.Close()

	val, err := f.GetCellValue(metaSheet, "A1")
	if err != nil {
		// Meta sheet 不存在时 excelize 会返回错误，视为空 checkpoint
		return "", nil
	}
	return val, nil
}

// WriteCheckpointHash 将 checkpoint commit hash 写入 Meta 工作表的 A1 单元格。
func WriteCheckpointHash(f *excelize.File, hash string) error {
	if f == nil {
		return fmt.Errorf("excelize.File 为 nil")
	}

	// 若 Meta sheet 不存在则创建
	if idx, _ := f.GetSheetIndex(metaSheet); idx == -1 {
		if _, err := f.NewSheet(metaSheet); err != nil {
			return fmt.Errorf("创建 Meta 工作表失败: %w", err)
		}
	}

	if err := f.SetCellValue(metaSheet, "A1", hash); err != nil {
		return fmt.Errorf("写入 checkpoint hash 失败: %w", err)
	}

	if err := f.Save(); err != nil {
		return fmt.Errorf("保存文件失败：%w", err)
	}

	return nil
}
