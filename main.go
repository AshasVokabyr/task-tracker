package main


import (
        "context"
        "fmt"
        "log"
        "os"
        "time"

        "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
        databaseURL := os.Getenv("DATABASE_URL")
        if databaseURL == "" {
                log.Fatal("DATABASE_URL is not set")
        }

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        pool, err := pgxpool.New(ctx, databaseURL)
        if err != nil {
                log.Fatal(err)
        }
        defer pool.Close()

        if err := pool.Ping(ctx); err != nil {
                log.Fatal(err)
        }
        var store TaskStorage = NewPostgresTaskStorage(pool)
        // Create
	    task, err := store.CreateTask(ctx, "Сделать ДЗ", "CRUD на Postgres", "todo", nil)
	    if err != nil {
		    log.Fatal(err)
	    }
	    fmt.Printf("CreateTask: %+v\n", task)

    	// Read
    	found, err := store.GetTask(ctx, task.id)
    	if err != nil {
    		log.Fatal(err)
    	}
    	fmt.Printf("GetTask: %+v\n", found)

    	// Update
    	updated, err := store.UpdateTask(ctx, task.id, "Сделать ДЗ по Go", "CRUD на Postgres через pgx", "in_progress", nil)
    	if err != nil {
    		log.Fatal(err)
    	}
    	fmt.Printf("UpdateTask: %+v\n", updated)

    	// Delete
    	if err := store.DeleteTask(ctx, task.id); err != nil {
    		log.Fatal(err)
    	}
    	fmt.Println("DeleteTask: ok")

}
