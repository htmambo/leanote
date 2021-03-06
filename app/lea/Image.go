package lea

import (
	"context"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/revel/revel"
)

func waterJpeg(path string) (ok bool, transPath string) {
	// 水印图片文件
	water := revel.BasePath + "/public/images/watermark.png"
	img_file, err := os.Open(path)
	defer img_file.Close()
	if err != nil {
		//fmt.Println("打开图片出错")
		//fmt.Println(err)
		return ok, path
	}
	img, err := jpeg.Decode(img_file)
	if err != nil {
		//fmt.Println("把图片解码为结构体时出错")
		//fmt.Println(img)
		return ok, path
	}

	wmb_file, err := os.Open(water)
	if err != nil {
		//fmt.Println("打开水印图片" + water + "出错")
		//fmt.Println(err)
		return ok, path
	}
	wmb_img, err := png.Decode(wmb_file)
	if err != nil {
		defer wmb_file.Close()
		//fmt.Println("把水印图片解码为结构体时出错")
		//fmt.Println(err)
		return ok, path
	}

	//把水印写在右下角，并向0坐标偏移10个像素
	offset := image.Pt(img.Bounds().Dx()-wmb_img.Bounds().Dx()-10, img.Bounds().Dy()-wmb_img.Bounds().Dy()-10)
	b := img.Bounds()
	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(b)

	//image.ZP代表Point结构体，目标的源点，即(0,0)
	//draw.Src源图像透过遮罩后，替换掉目标图像
	//draw.Over源图像透过遮罩后，覆盖在目标图像上（类似图层）
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, wmb_img.Bounds().Add(offset), wmb_img, image.ZP, draw.Over)

	//生成新图片new.jpg,并设置图片质量
	//fmt.Println("写入文件")
	img_sfile, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0)
	//n, _ := img_file.Seek(0, os.SEEK_END)
	jpeg.Encode(img_sfile, m, &jpeg.Options{100})
	defer img_sfile.Close()
	defer wmb_file.Close()
	return ok, path
}

func waterPng(path string) (ok bool, transPath string) {
	// 水印图片文件
	water := revel.BasePath + "/public/images/watermark.png"
	img_file, err := os.Open(path)
	defer img_file.Close()
	if err != nil {
		//fmt.Println("打开图片出错")
		//fmt.Println(err)
		return ok, path
	}
	img, err := png.Decode(img_file)
	if err != nil {
		//fmt.Println("把图片解码为结构体时出错")
		//fmt.Println(img)
		return ok, path
	}

	wmb_file, err := os.Open(water)
	if err != nil {
		//fmt.Println("打开水印图片" + water + "出错")
		//fmt.Println(err)
		return ok, path
	}
	wmb_img, err := png.Decode(wmb_file)
	if err != nil {
		defer wmb_file.Close()
		//fmt.Println("把水印图片解码为结构体时出错")
		//fmt.Println(err)
		return ok, path
	}

	//把水印写在右下角，并向0坐标偏移10个像素
	offset := image.Pt(img.Bounds().Dx()-wmb_img.Bounds().Dx()-10, img.Bounds().Dy()-wmb_img.Bounds().Dy()-10)
	b := img.Bounds()
	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(b)

	//image.ZP代表Point结构体，目标的源点，即(0,0)
	//draw.Src源图像透过遮罩后，替换掉目标图像
	//draw.Over源图像透过遮罩后，覆盖在目标图像上（类似图层）
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, wmb_img.Bounds().Add(offset), wmb_img, image.ZP, draw.Over)

	//生成新图片new.jpg,并设置图片质量
	//fmt.Println("写入文件")
	img_sfile, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0)
	//n, _ := img_file.Seek(0, os.SEEK_END)
	png.Encode(img_sfile, m)
	defer img_sfile.Close()
	defer wmb_file.Close()
	return ok, path
}
func TransToGif(path string, maxWidth uint, afterDelete bool) (ok bool, transPath string) {
	//图片，网上随便找了一张
	img_file, err := os.Open(path)
	defer img_file.Close()
	if err != nil {
		//fmt.Println("打开图片出错")
		//fmt.Println(err)
		return ok, path
	}

	buff := make([]byte, 512)

	_, err = img_file.Read(buff)
	if err != nil {
		//fmt.Println("读取源文件" + path + "时出错")
		//fmt.Println(err)
		return ok, path
	}
	//水印的活交给七牛来完成
	if revel.Config.BoolDefault("qiniu.enabled", false) {
		//使用七牛云
		ok, fileurl := upload_qiniu(path)
		if ok {
			//fmt.Println("七牛云图片链接 ：" + fileurl)
			return ok, fileurl
		}
	}
	imgType := http.DetectContentType(buff)
	if imgType == "image/jpeg" {
		_, toPathGif := waterJpeg(path)
		return ok, toPathGif
	} else if imgType == "image/png" {
		_, toPathGif := waterPng(path)
		return ok, toPathGif
	} else {
		//fmt.Println("不支持的图片类型" + imgType)
		return ok, path
	}
}

func upload_qiniu(filePath string) (ok bool, transPath string) {
	//上传凭证,关于凭证这块大家可以去看看官方文档
	putPolicy := storage.PutPolicy{
		Scope: revel.Config.StringDefault("qiniu.bucket", ""),
	}
	mac := qbox.NewMac(revel.Config.StringDefault("qiniu.access_key", ""), revel.Config.StringDefault("qiniu.secret_key", ""))
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	//空间对应机房
	//其中关于Zone对象和机房的关系如下：
	//    机房    Zone对象
	//    华东    storage.ZoneHuadong
	//    华北    storage.ZoneHuabei
	//    华南    storage.ZoneHuanan
	//    北美    storage.ZoneBeimei
	//七牛云存储空间设置首页有存储区域
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	//构建上传表单对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	pos := len(revel.BasePath) + 35
	key := strings.TrimLeft(filePath[pos:], "/")
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra)
	if err != nil {
		//fmt.Println("上传七牛云失败")
		//fmt.Println(err)
		return
	}
	//fmt.Println("上传" + filePath + "到七牛云成功：" + key)
	return true, revel.Config.StringDefault("qiniu.domain", "https://img.imzhp.com/") + key
}
