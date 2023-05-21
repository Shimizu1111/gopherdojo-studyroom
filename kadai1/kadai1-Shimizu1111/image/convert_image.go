package image

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	InputPath  string
	Name       string
	OutputPath string
	Extension  string
}

type ImageExtension string

var (
	imageExtensionPng ImageExtension = "png"
	imageExtensionJpg ImageExtension = "jpg"
)

func ConvertImage(path string) {
	var inputImagePath string
	if path == "" {
		i, err := filepath.Abs("./image/image/input")
		if err != nil {
			panic(err)
		}
		inputImagePath = i
	} else {
		inputImagePath = path
	}
	outputImagePath, err := filepath.Abs("./image/image/output")
	if err != nil {
		panic(err)
	}

	inputImagepaths, err := getImagePaths(inputImagePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(inputImagepaths)

	// 対象ディレクトリ内の画像を再帰的にconvert
	for _, path := range inputImagepaths {
		// TODO: 複数の拡張子を選択できるようにする
		img := &Image{
			InputPath:  path,
			Name:       filepath.Base(path),
			OutputPath: outputImagePath,
			Extension:  filepath.Ext(path),
		}
		switch img.Extension {
		case "." + string(imageExtensionJpg):
			if err := img.convert(); err != nil {
				panic(err)
			}
		default:
			continue
		}
	}
}

// 対象ディレクトリ配下の全てのディレクトリpathを取得
func getImagePaths(path string) ([]string, error) {
	var paths []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func (i *Image) convert() error {
	// 画像ファイルのopen
	file, err := os.Open(i.InputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	// jpegファイルをimage型に変換
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("画像のデコードに失敗しました")
		return err
	}
	// image型にしたjpeg画像をpng画像にしてbufに書き込み
	var buf bytes.Buffer
	if err := png.Encode(&buf, jpegImage); err != nil {
		return err
	}
	// 対象の拡張子の画像を作成
	fileName := filepath.Base(i.InputPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	outputImage, err := os.Create(i.OutputPath + "/" + fileNameWithoutExt + "." + string(imageExtensionPng))
	if err != nil {
		return err
	}
	// 対象画像ファイルにbufの内容を書き込む
	_, err = outputImage.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
