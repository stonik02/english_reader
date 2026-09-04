package uploadbook

import (
	"errors"
	"io"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/library/upload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	usecase UseCase
	tokens  TokenParser
}

func New(usecase UseCase, tokens TokenParser) *Handler {
	return &Handler{usecase: usecase, tokens: tokens}
}

func (h *Handler) UploadBook(stream grpc.ClientStreamingServer[readerv1.UploadBookRequest, readerv1.Book]) error {
	first, err := stream.Recv()
	if err == io.EOF {
		return status.Error(codes.InvalidArgument, "upload metadata is required")
	}
	if err != nil {
		return status.Error(codes.InvalidArgument, "read upload metadata")
	}
	metadata := first.GetMetadata()
	if metadata == nil || metadata.GetFilename() == "" {
		return status.Error(codes.InvalidArgument, "upload metadata is required")
	}
	userID, err := h.tokens.Parse(metadata.GetAccessToken())
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid access token")
	}

	reader, writer := io.Pipe()
	streamErr := make(chan error, 1)
	go copyChunks(stream, writer, streamErr)

	book, uploadErr := h.usecase.Execute(stream.Context(), uc.Request{
		UserID:   userID,
		Filename: metadata.GetFilename(),
		File:     reader,
	})
	_ = reader.Close()
	receiveErr := <-streamErr
	if uploadErr != nil {
		return response.Error(uploadErr)
	}
	if receiveErr != nil {
		return status.Error(codes.InvalidArgument, "invalid upload stream")
	}

	return stream.SendAndClose(response.Book(book))
}

func copyChunks(stream grpc.ClientStreamingServer[readerv1.UploadBookRequest, readerv1.Book], writer *io.PipeWriter, result chan<- error) {
	defer close(result)
	defer writer.Close()

	for {
		request, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			result <- nil
			return
		}
		if err != nil {
			_ = writer.CloseWithError(err)
			result <- err
			return
		}
		chunk := request.GetChunk()
		if chunk == nil {
			err := errors.New("upload chunk is required after metadata")
			_ = writer.CloseWithError(err)
			result <- err
			return
		}
		if _, err := writer.Write(chunk); err != nil {
			result <- err
			return
		}
	}
}
