package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

type ProtoDeserializer struct {
	Message *dynamic.Message
}

func (pd *ProtoDeserializer) Deserialize(data []byte) error {
	return pd.Message.Unmarshal(data)
}

func (pd *ProtoDeserializer) ToJSON() (map[string]interface{}, error) {
	jsonBytes, err := pd.Message.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return result, nil
}

func (pd *ProtoDeserializer) GetFieldMappings() map[string]string {
	mappings := make(map[string]string)
	fields := pd.Message.GetMessageDescriptor().GetFields()

	for _, field := range fields {
		mappings[field.GetJSONName()] = string(field.GetName())
	}

	return mappings
}

func parseProtoSchema(schema string, messageType string) (*desc.MessageDescriptor, error) {
	tempDir, err := os.MkdirTemp("", "proto-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	protoFile := filepath.Join(tempDir, "schema.proto")
	if err := os.WriteFile(protoFile, []byte(schema), 0644); err != nil {
		return nil, fmt.Errorf("failed to write schema file: %v", err)
	}

	parser := protoparse.Parser{
		ImportPaths: []string{tempDir},
	}
	files, err := parser.ParseFiles("schema.proto")
	if err != nil {
		return nil, fmt.Errorf("failed to parse proto schema: %v", err)
	}

	var messageDesc *desc.MessageDescriptor
	for _, file := range files {
		messageDesc = file.FindMessage(messageType)
		if messageDesc != nil {
			break
		}
	}
	if messageDesc == nil {
		return nil, fmt.Errorf("message type not found in schema")
	}

	return messageDesc, nil
}

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/deserialize", func(c *gin.Context) {
		schema := c.PostForm("schema")
		messageType := c.PostForm("messageType")
		hexMessage := c.PostForm("message")

		if schema == "" || messageType == "" || hexMessage == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Schema, message type, and hex message are required"})
			return
		}

		messageDesc, err := parseProtoSchema(schema, messageType)
		if err != nil {
			log.Printf("Failed to parse schema: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		message := dynamic.NewMessage(messageDesc)
		deserializer := &ProtoDeserializer{Message: message}

		hexStr := strings.TrimSpace(hexMessage)
		decodedData, err := hex.DecodeString(hexStr)
		if err != nil {
			log.Printf("Failed to decode hex string: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hex string: " + err.Error()})
			return
		}

		if err := deserializer.Deserialize(decodedData); err != nil {
			log.Printf("Failed to deserialize message: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to deserialize message: " + err.Error()})
			return
		}

		jsonData, err := deserializer.ToJSON()
		if err != nil {
			log.Printf("Failed to convert to JSON: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert to JSON: " + err.Error()})
			return
		}

		mappings := deserializer.GetFieldMappings()

		c.JSON(http.StatusOK, gin.H{
			"json":     jsonData,
			"mappings": mappings,
		})
	})

	r.POST("/generate", func(c *gin.Context) {
		schema := c.PostForm("schema")
		messageType := c.PostForm("messageType")
		jsonData := c.PostForm("jsonData")

		if schema == "" || messageType == "" || jsonData == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Schema, message type, and JSON data are required"})
			return
		}

		messageDesc, err := parseProtoSchema(schema, messageType)
		if err != nil {
			log.Printf("Failed to parse schema: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		message := dynamic.NewMessage(messageDesc)

		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(jsonData), &jsonMap); err != nil {
			log.Printf("Failed to parse JSON data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data: " + err.Error()})
			return
		}

		if err := message.UnmarshalJSON([]byte(jsonData)); err != nil {
			log.Printf("Failed to unmarshal JSON to message: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert JSON to message: " + err.Error()})
			return
		}

		protoBytes, err := message.Marshal()
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate protobuf message: " + err.Error()})
			return
		}

		hexMessage := hex.EncodeToString(protoBytes)

		c.JSON(http.StatusOK, gin.H{
			"hex": hexMessage,
		})
	})

	log.Printf("Starting server on :8080")
	r.Run(":8080")
}
