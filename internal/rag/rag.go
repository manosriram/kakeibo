package rag

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"log"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores"

	"github.com/tmc/langchaingo/vectorstores/qdrant"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

type RAG struct {
	// Docs schema.Document
}

func NewRAG() RAG {
	return RAG{}
}

func (r *RAG) Query(query string) (string, error) {
	docs, err := r.loadDocuments()
	if err != nil {
		fmt.Println("Error ", err.Error())
		return "", err
	}

	splitter := textsplitter.NewRecursiveCharacter(textsplitter.WithChunkOverlap(150), textsplitter.WithChunkSize(1000))

	var chunks []schema.Document
	for _, doc := range docs {
		texts, err := splitter.SplitText(doc.PageContent)
		if err != nil {
			fmt.Println("Error ", err.Error())
			return "", err
		}

		for _, text := range texts {
			chunks = append(chunks, schema.Document{
				PageContent: text,
				Metadata:    doc.Metadata,
			})
		}
	}

	o, err := openai.New(openai.WithEmbeddingModel("text-embedding-3-small"))
	if err != nil {
		fmt.Println("Error ", err.Error())
		return "", err
	}

	embedder, err := embeddings.NewEmbedder(o)
	if err != nil {
		fmt.Println("Error ", err.Error())
		return "", err
	}

	qdrantURL, err := url.Parse("http://qdrant_server:6333")
	if err != nil {
		log.Fatal(err)
	}

	store, err := qdrant.New(
		qdrant.WithURL(*qdrantURL),
		qdrant.WithCollectionName("kakeibo_knowledge_base"),
		qdrant.WithEmbedder(embedder),
	)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	_, err = store.AddDocuments(context.Background(), chunks)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	retriever := vectorstores.ToRetriever(store, 1) // Get top 3 similar docs
	relevantDocs, err := retriever.GetRelevantDocuments(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	c := r.buildContext(relevantDocs)

	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	prompt := fmt.Sprintf(`You are an expert expense analyzer. You are given a csv file with expense data in it.
		Answer the questions and be as accurate as possible,

Context (expense csv):
%s

Question: %s

Answer:`, c, query)

	answer, err := llm.Call(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}

	return answer, nil
}

func (r *RAG) loadDocuments() ([]schema.Document, error) {
	file, err := os.Open("spends.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loader := documentloaders.NewText(file)
	docs, err := loader.Load(context.Background())
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *RAG) buildContext(docs []schema.Document) string {
	var context strings.Builder
	for i, doc := range docs {
		fmt.Fprintf(&context, "Document %d:\n%s\n\n", i+1, doc.PageContent)
	}
	return context.String()
}
