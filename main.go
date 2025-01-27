package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
)

var db *pgx.Conn

type QueryRequest struct {
    SQL   string        `json:"sql"`  
    Args  []interface{} `json:"args"`  
}

func main() {
    var err error
    dbURL := "postgres://postgres:20032003@40.76.113.138:5432/medicalrecord"
    db, err = pgx.Connect(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Error conectando a la BD: %v", err)
    }
    defer db.Close(context.Background())

    r := gin.Default()

    r.POST("/query", executeQuery)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8001"
    }

    log.Printf("Servidor corriendo en :%s", port)
    r.Run(":" + port)
}

func executeQuery(c *gin.Context) {
    var req QueryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
        return
    }

    rows, err := db.Query(context.Background(), req.SQL, req.Args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    columns := rows.FieldDescriptions()
    results := []map[string]interface{}{}

    for rows.Next() {
        values, err := rows.Values()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        rowMap := map[string]interface{}{}
        for i, col := range columns {
            rowMap[string(col.Name)] = values[i]
        }
        results = append(results, rowMap)
    }

    c.JSON(http.StatusOK, gin.H{"data": results})
}
