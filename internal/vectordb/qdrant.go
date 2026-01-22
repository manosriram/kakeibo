package vectordb

import (
	"context"
	"fmt"
	"time"

	"github.com/manosriram/kakeibo/internal/llm"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantVectorDB struct {
	Client    *qdrant.Client
	LlmClient *llm.OpenAI // TODO: update this to use llm interface
}

func NewQdrantVectorDB() (*QdrantVectorDB, error) {
	llmClient := llm.NewOpenAI()
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		return nil, err
	}

	return &QdrantVectorDB{
		Client:    client,
		LlmClient: &llmClient,
	}, nil
}

func (q *QdrantVectorDB) CreateCollection() {
	q.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "statements",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1536,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func (q *QdrantVectorDB) AddPayload(payload map[string]any) error {
	vector, err := q.LlmClient.CreateEmbedding(payload)
	if err != nil {
		return err
	}

	vectorFloats := []float32{}
	for _, v := range vector {
		vectorFloats = append(vectorFloats, v.Embedding...)
	}

	result, err := q.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "statements",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(uint64(time.Now().UnixNano())), // Unique ID
				Vectors: qdrant.NewVectors(vectorFloats...),
				Payload: qdrant.NewValueMap(map[string]any{
					"amount":           payload["amount"],
					"tag":              payload["tag"],
					"transaction_type": payload["transaction_type"],
					"created_at":       payload["created_at"],
					"description":      payload["description"],
				}),
			},
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("result = ", result)
	return nil
}

func (q *QdrantVectorDB) Query(qry string) error {
	return nil
}
