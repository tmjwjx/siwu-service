package logics

import (
	"bytes"
	"fmt"
	"forum/internal/image/repositories"
	"forum/internal/image/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"github.com/disintegration/imaging"
	"github.com/ericpauley/go-quantize/quantize"
	"github.com/go-redis/redis/v8"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// CompressImageLogic 压缩图片
func CompressImageLogic(rdb *redis.Client, req requests.CompressImageReq) (string, []byte, error) {
	// 安全处理路径
	// 将上传图片的url路径处理成其在服务器中的路径
	cleanPath, err := internalUtils.ProcessImagePath(req.Path)
	if err != nil {
		return "", nil, fmt.Errorf("CompressImageLogic -> 将上传图片的url路径处理成其在服务器中的路径异常 -> %v", err)
	}

	// 获取响应类型
	contentType := getContentType(cleanPath)

	// 生成缓存键
	cacheKey := generateCacheKey(cleanPath, req.Width, req.Height, req.Level)

	// 检查缓存
	if cached, err := repositories.GetFromCache(rdb, cacheKey); err == nil {
		return contentType, cached, nil
	}

	// 处理图片
	processed, err := processImage(cleanPath, req.Width, req.Height, req.Level)
	if err != nil {
		return "", nil, fmt.Errorf("CompressImageLogic -> 处理图片异常 -> %v", err)
	}

	// 保存到缓存
	if err = repositories.SetToCache(rdb, cacheKey, processed); err != nil {
		return "", nil, fmt.Errorf("CompressImageLogic -> cache save failed -> %v", err)
	}

	return contentType, processed, nil
}

// 图片处理核心逻辑
func processImage(path string, width, height, level int) ([]byte, error) {
	// 读取原始文件
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("processImage -> 读取原始文件异常 -> %v", err)
	}
	// 关闭文件
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			globals.Log.Errorf("processImage -> 关闭文件失败 -> %v", err)
		}
	}(file)

	// 解码图片
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// 调整尺寸
	resizedImg := resizeImage(img, width, height)

	// 创建输出缓冲区
	var buf bytes.Buffer

	// 根据格式进行压缩
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		quality := calculateJPEGQuality(level)
		err = imaging.Encode(&buf, resizedImg, imaging.JPEG, imaging.JPEGQuality(quality))
	case "png":
		compressionLevel := calculatePNGCompression(level)
		err = encodePNG(&buf, resizedImg, compressionLevel)
	case "gif":
		err = processGIF(path, &buf, width, height, level)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 缓存相关函数
func generateCacheKey(path string, width, height, level int) string {
	return fmt.Sprintf("img:%s:%d:%d:%d", path, width, height, level)
}

// 辅助函数
func getContentType(filePath string) string {
	// 获得图片的扩展名(带.的)
	ext := filepath.Ext(filePath)
	// 将扩展名的每个字母都变成小写字母
	switch strings.ToLower(ext) {
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// 调整图片尺寸(放大或缩小)
func resizeImage(img image.Image, width, height int) image.Image {
	switch {
	case width == 0 && height == 0:
		return img
	case width == 0:
		return imaging.Resize(img, 0, height, imaging.Lanczos)
	case height == 0:
		return imaging.Resize(img, width, 0, imaging.Lanczos)
	default:
		return imaging.Resize(img, width, height, imaging.Lanczos)
	}
}

// JPEG质量计算
func calculateJPEGQuality(level int) int {
	if level == 0 {
		level = 9
	} else if level == 1 {
		level = 8
	} else if level == 2 {
		level = 7
	} else if level == 3 {
		level = 6
	} else if level == 4 {
		level = 5
	} else if level == 5 {
		level = 4
	} else if level == 6 {
		level = 3
	} else if level == 7 {
		level = 2
	} else if level == 8 {
		level = 1
	} else if level == 9 {
		level = 0
	} else {
		level = 9
	}
	return level
}

// PNG压缩级别计算
func calculatePNGCompression(level int) int {
	return level // 0-9直接对应
}

// 自定义PNG编码
// w *bytes.Buffer：压缩后数据的输出目标。
// img image.Image：要编码的原始图片数据。
func encodePNG(w *bytes.Buffer, img image.Image, level int) error {
	// 配置 PNG 编码器的结构体
	encoder := png.Encoder{
		CompressionLevel: png.CompressionLevel(level),
	}
	// 将图片数据编码为 PNG 格式
	return encoder.Encode(w, img)
}

// 压缩GIF格式的图片
// path：GIF 图片的文件路径。
// w：指向 bytes.Buffer 的指针，用于存储压缩后的 GIF 数据。
// width 和 height：目标图片的宽度和高度。
// level：压缩级别（0-9）。
func processGIF(path string, w *bytes.Buffer, width, height, level int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			globals.Log.Errorf("processGIF -> 关闭文件失败 -> %v", err)
		}
	}(file)

	// 解码原始GIF
	// 使用 gif.DecodeAll 函数解码 GIF 文件，将所有帧和相关元数据（如延迟时间、循环次数）存储在 origGif 中
	origGif, err := gif.DecodeAll(file)
	if err != nil {
		return err
	}

	// 初始化处理后的GIF
	// 初始化一个新的 gif.GIF 结构体，用于存储处理后的 GIF 数据。
	// 保留原始 GIF 的 LoopCount（循环次数）。
	processedGif := &gif.GIF{
		LoopCount: origGif.LoopCount,
		Delay:     make([]int, 0, len(origGif.Delay)),
	}

	// 计算压缩参数
	// colorCount：根据压缩级别 level 计算颜色数量。
	// level=0：256 - 0*25 = 256 → 256 色。
	// level=9：256 - 9*25 = 31 → 31 色。
	// frameStep：根据压缩级别 level 计算帧间隔。
	// level=0-2：math.Pow(2, 0) = 1 → 每帧都保留。
	// level=3-5：math.Pow(2, 1) = 2 → 每隔一帧保留。
	// level=6-8：math.Pow(2, 2) = 4 → 每隔三帧保留。
	// level=9：math.Pow(2, 3) = 8 → 每隔七帧保留。
	colorCount := 256 - level*25                    // 颜色数量（255到31）level 0=256色，level9=31色
	frameStep := int(math.Pow(2, float64(level/3))) // 帧间隔（1到8）level 0-2=1帧，3-5=2帧，6-8=4帧，9=8帧

	// 处理每一帧
	for i := 0; i < len(origGif.Image); i += frameStep {
		// 调整尺寸
		resizedFrame := resizeImage(origGif.Image[i], width, height)

		// 颜色量化
		palettedFrame := quantizeImage(resizedFrame, colorCount)

		// 添加处理后的帧
		processedGif.Image = append(processedGif.Image, palettedFrame)
		processedGif.Delay = append(processedGif.Delay, origGif.Delay[i])
	}

	// 设置全局调色板（使用第一帧的调色板）
	// 如果处理后的 GIF 图片包含至少一帧，设置全局调色板为第一帧的调色板。
	if len(processedGif.Image) > 0 {
		processedGif.Config = image.Config{
			ColorModel: processedGif.Image[0].Palette,
			Width:      processedGif.Image[0].Bounds().Dx(),
			Height:     processedGif.Image[0].Bounds().Dy(),
		}
	}

	// 使用 gif.EncodeAll 函数将处理后的 GIF 数据编码为二进制数据，并写入到缓冲区 w
	return gif.EncodeAll(w, processedGif)
}

// 颜色量化实现
// img image.Image：要进行颜色量化的原始图像。
// maxColors int：指定生成的调色板中允许的最大颜色数量。
// *image.Paletted：返回一个 image.Paletted 类型的图像，包含生成的调色板和量化后的颜色索引。
func quantizeImage(img image.Image, maxColors int) *image.Paletted {
	// 使用专业量化算法
	// MedianCutQuantizer 是一种基于中值分割的颜色量化算法。
	// 它通过将颜色空间递归分割成多个子盒子，每个子盒子里包含一定数量的颜色，然后选择代表颜色来生成调色板。
	q := quantize.MedianCutQuantizer{}

	// 调用 Quantize 方法生成调色板。
	// 第一个参数是一个初始的颜色切片，这里创建了一个容量为 maxColors 的空切片。
	// 第二个参数是原始图像。
	// 返回值 palette 是生成的调色板。
	palette := q.Quantize(make([]color.Color, 0, maxColors), img)

	// 创建调色板图像
	//创建一个新的 image.Paletted 类型的图像。
	//第一个参数是图像的边界（img.Bounds()）。
	//第二个参数是生成的调色板。
	paletted := image.NewPaletted(img.Bounds(), palette)

	// 转换颜色空间
	// 循环遍历图像的每个像素：
	// 使用 img.Bounds() 获取图像的边界。
	// 使用双层循环遍历图像的每个像素。
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {

			// paletted.Set: 将像素颜色设置为调色板中的颜色。
			// color.RGBAModel.Convert(img.At(x, y)) 将原始图像的颜色转换为 RGBA 格式。
			paletted.Set(x, y, color.RGBAModel.Convert(img.At(x, y)))
		}
	}
	return paletted
}
