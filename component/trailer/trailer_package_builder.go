package trailer

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

/*
*
* APP包构建器, 用来解压ZIP包, 不过0.6.4暂不支持
*
 */
// 压缩多个文件到指定目录
func Zip(files []string, targetPath string) error {
	// 创建目标文件
	zipFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建 zip 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历并压缩文件列表
	for _, sourcePath := range files {
		// 打开要压缩的文件
		srcFile, err := os.Open(sourcePath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 获取文件信息
		srcInfo, err := srcFile.Stat()
		if err != nil {
			return err
		}

		// 创建 zip 文件条目
		zipFileWriter, err := zipWriter.Create(srcInfo.Name())
		if err != nil {
			return err
		}

		// 复制文件内容到 zip 条目
		_, err = io.Copy(zipFileWriter, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// 解压缩文件到指定目录
func Unzip(zipPath string, targetDir string) error {
	// 打开要解压的 zip 文件
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 遍历 zip 文件中的文件并解压
	for _, file := range zipFile.File {
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destPath := targetDir + "/" + file.Name
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func __test() {
	// 压缩文件
	if err := Zip(
		[]string{"manifest-native.json",
			"./app1.exe"}, "app-native.zip"); err != nil {
		fmt.Println("Error compressing file:", err)
		return
	}

	fmt.Println("File compressed successfully")

	// 解压文件
	if err := Unzip("./app-native.zip", "./app-native"); err != nil {
		fmt.Println("Error decompressing file:", err)
		return
	}

	fmt.Println("File decompressed successfully")
}
