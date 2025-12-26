package translator

import (
	"context"

	// Import the generated code
	pb "github.com/justinhjy1004/sentenceminer/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 1. Define a struct to hold the persistent client
type TranslationService struct {
	client pb.TranslatorClient
	conn   *grpc.ClientConn
}

// 2. Initialize the connection ONCE (e.g., at app startup)
func NewTranslationService(address string) (*TranslationService, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &TranslationService{
		client: pb.NewTranslatorClient(conn),
		conn:   conn,
	}, nil
}

// 3. The method now only handles the request logic
func (s *TranslationService) TranslateText(ctx context.Context, text string) (string, error) {
	req := &pb.TranslateRequest{
		Text:       text,
		SourceLang: "de",
		TargetLang: "es",
	}

	r, err := s.client.Translate(ctx, req)
	if err != nil {
		return "", err
	}

	return r.TranslatedText, nil
}

// 4. Remember to close the connection when the app shuts down
func (s *TranslationService) Close() {
	s.conn.Close()
}
