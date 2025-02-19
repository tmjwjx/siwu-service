package logics

import (
	"fmt"
	"forum/internal/image/requests"
	"forum/pkg/globals"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"image"
	"strconv"
	"strings"
)

// 解码上传的多张图片
func decodeImages(c *gin.Context) ([]image.Image, []string, error) {
	// 使用 MultipartForm 提取所有字段
	// 获取上传的所有文件
	form, err := c.MultipartForm()
	if err != nil {
		return nil, nil, fmt.Errorf("decodeImage -> 获取上传的所有文件失败 -> %s", err)
	}

	// 提取文件
	files := form.File["files"]

	// 如果文件中没有图片，直接返回nil。
	if len(files) == 0 {
		return nil, nil, fmt.Errorf("文件中没有图片")
	}

	var images []image.Image
	var formats []string

	// 遍历每个上传的文件
	for _, fileHeader := range files {

		//// 如果文件中没有图片，直接返回nil。
		//if fileHeader.Size == 0 {
		//	return nil, nil, fmt.Errorf("文件中没有图片")
		//}

		// 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			return nil, nil, fmt.Errorf("decodeImage -> 打开文件失败 -> %s", err)
		}

		defer file.Close()

		// 解码图片
		img, format, err := image.Decode(file)
		if err != nil {
			return nil, nil, fmt.Errorf("decodeImage -> 解码图片失败 -> %s", err)
		}

		// 保存解码后的图片和格式
		images = append(images, img)
		formats = append(formats, format)
	}

	return images, formats, nil

}

// ProduceUrlLogic 图片文件的逻辑处理
func ProduceUrlLogic(c *gin.Context, watermarkParam requests.WatermarkParam) (*requests.ImageUrl, error) {

	var finalImage image.Image
	// 读取图片并进行处理
	images, formats, err := decodeImages(c)
	if err != nil {
		if strings.Contains(err.Error(), "文件中没有图片") {
			return nil, fmt.Errorf("ProduceUrlLogic -> 文件中没有图片")
		} else {
			return nil, fmt.Errorf("ProduceUrlLogic -> 生成图片的url路径失败 -> %s", err)
		}
	}

	// 获取上传的图片属于的类型（文章，用户，标签，评论)
	typeParam := c.PostForm("type")

	var urls []*requests.UrlPath
	// 处理解码后的图片（例如压缩、加水印等）
	for i, img := range images {

		if typeParam == "文章封面" {
			if float64(img.Bounds().Dx()/img.Bounds().Dy()) != watermarkParam.Scale {
				return nil, fmt.Errorf("上传的图片比例不对")
			}
		}

		if typeParam == "文章" {
			// 添加水印
			imgWithWatermark := addWatermark(img, watermarkParam)
			finalImage = imgWithWatermark
		} else {
			finalImage = img
		}

		// 生成唯一的文件名
		uniqueFilename := generateUniqueFilename(formats[i])

		// 生成url
		imageUrl := "http://" + globals.AppConfig.App.Domain + ":" + strconv.Itoa(globals.AppConfig.App.Port) + globals.SConfig.Prefix + "/" + uniqueFilename

		url := &requests.UrlPath{
			Url: imageUrl,
		}
		urls = append(urls, url)

		// 定义存储的完整路径，确保以原始格式的扩展名结尾
		outFile := fmt.Sprintf(globals.SConfig.Path + "/" + uniqueFilename)

		// 根据原始格式保存文件
		switch formats[i] {
		case "jpeg", "jpg":
			err = imaging.Save(finalImage, outFile, imaging.JPEGQuality(80))
		case "png":
			err = imaging.Save(finalImage, outFile, imaging.PNGCompressionLevel(5))
		case "gif":
			err = imaging.Save(finalImage, outFile)
		default:
			return nil, fmt.Errorf("ProduceUrlLogic -> 不支持的格式: %s -> %s", formats[i], err)
		}
		if err != nil {
			return nil, fmt.Errorf("ProduceUrlLogic -> 根据原始格式保存图片文件失败 -> %s", err)
		}

	}

	res := &requests.ImageUrl{
		Data:  urls,
		Errno: 0,
	}

	return res, nil
}

// 生成唯一文件名
func generateUniqueFilename(format string) string {
	// 生成一个唯一的 UUID
	uniqueID := uuid.New().String()

	// 创建一个新的唯一文件名
	return fmt.Sprintf("%s.%s", uniqueID, format)
}

// 给图片添加水印
func addWatermark(img image.Image, watermarkParam requests.WatermarkParam) image.Image {
	// 使用 gg 库进行绘图
	dc := gg.NewContextForImage(img)

	// 设置字体大小和水印颜色
	if err := dc.LoadFontFace("pkg/font/dingliezhuhaifont-20240831GengXinBan)-2.ttf", watermarkParam.Size); err != nil {
		globals.Log.Errorf("加载字体失败: %v", err)
	}

	dc.SetRGBA(0.4431, 0.4667, 0.4902, 1) //半透明白色

	// 获取图片尺寸
	w, h := dc.Width(), dc.Height()

	// 计算水印的位置，稍微向内偏移，防止水印超出图片边界
	margin := watermarkParam.Margin // 边距，避免水印靠得太边
	x := float64(w) - margin        // x 位置从右边向左偏移一些
	y := float64(h) - 3*margin      // y 位置从下方向上偏移一些
	// 添加文本水印
	dc.DrawStringAnchored(watermarkParam.Watermark, x, y, 1.0, 1.0)

	// 返回处理后的图片
	return dc.Image()
}
