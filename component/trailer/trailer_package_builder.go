package trailer

/*
*
* 包验证器, 未来会增强应用管理功能, 允许上传ZIP包
*
 */
import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

/*
*
* 包内部配置
*
 */
type AppManifest struct {
	Native     bool     `json:"native"`
	ScriptHost string   `json:"scripthost"`
	Executable string   `json:"executable"`
	Env        []string `json:"env"`
}

/*
*
* 针对脚本语言
*
 */
var ExecuteType = map[string]string{
	".jar": "JAVA",
	".exe": "EXE",
	".py":  "PYTHON",
	".js":  "NODEJS",
	".lua": "LUA",
}

func ValidatePackage(mf AppManifest) error {
	return nil
}

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

func ReadManifestFromZip(zipPath string, targetFile string) (AppManifest, error) {
	var manifest AppManifest

	// 打开 ZIP 文件
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return manifest, err
	}
	defer zipFile.Close()

	// 查找并打开目标文件
	var found bool
	var targetReader io.ReadCloser
	for _, file := range zipFile.File {
		if file.Name == targetFile {
			targetReader, err = file.Open()
			if err != nil {
				return manifest, err
			}
			found = true
			break
		}
	}
	if !found {
		return manifest, errors.New("AppManifest not found in the ZIP archive")
	}
	defer targetReader.Close()

	// 解析 JSON 数据
	decoder := json.NewDecoder(targetReader)
	if err := decoder.Decode(&manifest); err != nil {
		return manifest, err
	}

	return manifest, nil
}

/*
*
  - 解压文件的路径: unzip -d ./upload/TrailerGoods/GOODSQUNK3Z/ app.zip
    路径下应该包含了所有解压的文件
    ll ./upload/TrailerGoods/GOODSQUNK3Z
  - manifest.json
  - app.exe

*
*/
func CreateUnzipPath(dir string) error {
	if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
		return err
	}
	return nil
}

/*
*
* 删除解压文件夹 rm -r ./upload/TrailerGoods/GOODSQUNK3Z
*
 */
func RemoveUnzipPath(dir string) error {
	if err := os.RemoveAll(filepath.Dir(dir)); err != nil {
		return err
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
