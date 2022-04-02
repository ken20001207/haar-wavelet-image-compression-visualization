package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/ken20001207/image-compressor/haar"
	"github.com/ken20001207/image-compressor/services"
	"github.com/ken20001207/image-compressor/utils"
	"image/jpeg"
	"strconv"
)

func DownloadImage(c *gin.Context) {

	// level means the number of haar transform
	level := 2
	if c.Query("level") != "" {
		l, _ := c.GetQuery("level")
		ls, err := strconv.Atoi(l)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "level must be a number",
			})
			return
		}
		level = ls
	}

	// step means the number of transform step
	step := 0
	if c.Query("step") != "" {
		s, _ := c.GetQuery("step")
		ss, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "step must be a number",
			})
			return
		}
		step = ss
	}

	image, err := services.GetImage(c.Request.Context(), c.Query("uid"))
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Image not found",
		})
		c.Abort()
		return
	}

	array, err := utils.ConvertImageTo2dArray(*image)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	// execute haar transform to the specified step

	s := 0
	for s <= step {

		// transform
		if s != 0 && s < 2*level+1 {
			horizontal := s%2 == 1
			if horizontal {
				haar.Horizontal(array, (s-1)/2+1)
			} else {
				haar.Vertical(array, (s-1)/2+1)
			}
		}

		// compress
		if s == 2*level+1 {
			haar.Compress(array, 0.1)
		}

		// reverse
		if s > 2*level+1 {
			horizontal := s%2 == 1
			if horizontal {
				haar.ReverseHaarHorizontal(array, level-(s-level)/2+2)
			} else {
				haar.ReverseHaarVertical(array, level-(s-level)/2+2)
			}
		}

		s++
	}

	// 1 - haar.HaarHorizontal(array, 1)
	// 2 - haar.Vertical(array, 1)
	// 3 - haar.HaarHorizontal(array, 2)
	// 4 - haar.Vertical(array, 2)
	// 5 - haar.HaarHorizontal(array, 3)
	// 6 - haar.Vertical(array, 3)
	//haar.Compress(array, 0.1)
	//haar.ReverseHaarVertical(array, 3)
	//haar.ReverseHaarHorizontal(array, 3)
	//haar.ReverseHaarVertical(array, 2)
	//haar.ReverseHaarHorizontal(array, 2)
	//haar.ReverseHaarVertical(array, 1)
	//haar.ReverseHaarHorizontal(array, 1)

	// encode result into jpeg for display

	result := utils.Convert2DArrayToImage(array)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, result, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	c.Header("Content-Type", "image/jpeg")

	_, err = c.Writer.Write(buf.Bytes())
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
}
