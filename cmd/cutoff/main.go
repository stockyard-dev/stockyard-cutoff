package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-cutoff/internal/server";"github.com/stockyard-dev/stockyard-cutoff/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9090"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./cutoff-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("cutoff: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Cutoff — Self-hosted API rate limiter\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Check:      POST http://localhost:%s/api/check\n  Data:       %s\n  ─────────────────────────────────\n  Questions? hello@stockyard.dev\n\n",port,port,port,dataDir)
log.Printf("cutoff: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
