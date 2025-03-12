package sendblobhandle

import (
	constant "app/internal/constants"
	httpresponse "app/pkg/http_response"
	logapp "app/pkg/log"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func (h *sendblobHandle) StreamEncoding(ctx *gin.Context) {
	uuid := ctx.Query("uuid")

	dirOutput := fmt.Sprintf("data/video/%s/360", uuid)
	err := os.MkdirAll(dirOutput, os.ModePerm)
	if err != nil {
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("init-file-m3u8", err.Error(), constant.ERROR_LOG)
		return
	}
	fileM3U8 := fmt.Sprintf("%s/index.m3u8", dirOutput)

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("up-connection", err.Error(), constant.ERROR_LOG)
		return
	}

	// io
	inputReader, inputWriter := io.Pipe()
	chanBlob := make(chan []byte, 1*100*100)

	// Tạo tiến trình ffmpeg
	cmd := exec.Command("ffmpeg",
		"-f", "webm", // Định dạng đầu vào
		"-i", "pipe:0", // Nhận từ stdin
		"-vf", "scale=-2:360", // Giảm độ phân giải
		// "-vcodec", "copy",
		"-preset", "ultrafast", // Cấu hình preset nhanh
		"-vcodec", "libx264", // Bộ mã hóa video H.264
		"-acodec", "aac", // Bộ mã hóa âm thanh AAC
		"-hls_time", "5", // Mỗi segment .ts có độ dài 4 giây
		"-hls_list_size", "5", // Giữ tối đa 10 segment trong playlist
		"-hls_flags", "delete_segments", // Xóa segment cũ để tiết kiệm bộ nhớ
		"-fflags", "+genpts", // Đảm bảo timestamp chính xác
		"-f", "hls", // Xuất dưới định dạng HLS
		fileM3U8, // Lưu danh sách phát và các file .ts
	)

	cmd.Stdin = inputReader
	cmd.Stderr = os.Stderr

	// Start ffmpeg
	err = cmd.Start()
	if err != nil {
		log.Fatalf("error start ffmpeg: %v", err)
	}
	defer func() {
		// Đóng WebSocket và dừng tiến trình ffmpeg
		conn.Close()
		close(chanBlob)
		inputWriter.Close() // Đảm bảo ffmpeg nhận EOF và thoát

		if cmd.Process != nil {
			cmd.Process.Kill() // Dừng tiến trình ffmpeg
		}
	}()

	// Goroutine để push dữ liệu vào ffmpeg
	go func() {
		for blob := range chanBlob {
			_, err = inputWriter.Write(blob)
			if err != nil {
				log.Println("Error encoding: ", err)
				break
			}
		}
	}()

	// Nhận dữ liệu từ WebSocket và gửi vào `chanBlob`
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Lỗi khi nhận tin nhắn từ WebSocket: %v", err)
			break
		}
		chanBlob <- data
	}
}
