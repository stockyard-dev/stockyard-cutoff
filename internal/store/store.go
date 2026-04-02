package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Rule struct{ID string `json:"id"`;Name string `json:"name"`;Key string `json:"key"`;Limit int `json:"limit"`;WindowSec int `json:"window_sec"`;Action string `json:"action"`;Enabled bool `json:"enabled"`;CreatedAt string `json:"created_at"`}
type CheckResult struct{Allowed bool `json:"allowed"`;Remaining int `json:"remaining"`;Limit int `json:"limit"`;ResetAt string `json:"reset_at"`;Key string `json:"key"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"cutoff.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
for _,q:=range[]string{
`CREATE TABLE IF NOT EXISTS rules(id TEXT PRIMARY KEY,name TEXT NOT NULL,key_pattern TEXT DEFAULT '*',lmt INTEGER DEFAULT 100,window_sec INTEGER DEFAULT 60,action TEXT DEFAULT 'reject',enabled INTEGER DEFAULT 1,created_at TEXT DEFAULT(datetime('now')))`,
`CREATE TABLE IF NOT EXISTS hits(id TEXT PRIMARY KEY,rule_id TEXT NOT NULL,client_key TEXT NOT NULL,created_at TEXT DEFAULT(datetime('now')))`,
`CREATE INDEX IF NOT EXISTS idx_hits_rule ON hits(rule_id,client_key,created_at)`,
}{if _,err:=db.Exec(q);err!=nil{return nil,fmt.Errorf("migrate: %w",err)}};return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)CreateRule(r *Rule)error{r.ID=genID();r.CreatedAt=now();en:=1;if!r.Enabled{en=0};if r.Action==""{r.Action="reject"};if r.Limit==0{r.Limit=100};if r.WindowSec==0{r.WindowSec=60}
_,err:=d.db.Exec(`INSERT INTO rules(id,name,key_pattern,lmt,window_sec,action,enabled,created_at)VALUES(?,?,?,?,?,?,?,?)`,r.ID,r.Name,r.Key,r.Limit,r.WindowSec,r.Action,en,r.CreatedAt);return err}
func(d *DB)GetRule(id string)*Rule{var r Rule;var en int;if d.db.QueryRow(`SELECT id,name,key_pattern,lmt,window_sec,action,enabled,created_at FROM rules WHERE id=?`,id).Scan(&r.ID,&r.Name,&r.Key,&r.Limit,&r.WindowSec,&r.Action,&en,&r.CreatedAt)!=nil{return nil};r.Enabled=en==1;return &r}
func(d *DB)ListRules()[]Rule{rows,_:=d.db.Query(`SELECT id,name,key_pattern,lmt,window_sec,action,enabled,created_at FROM rules ORDER BY name`);if rows==nil{return nil};defer rows.Close()
var o []Rule;for rows.Next(){var r Rule;var en int;rows.Scan(&r.ID,&r.Name,&r.Key,&r.Limit,&r.WindowSec,&r.Action,&en,&r.CreatedAt);r.Enabled=en==1;o=append(o,r)};return o}
func(d *DB)DeleteRule(id string)error{d.db.Exec(`DELETE FROM hits WHERE rule_id=?`,id);_,err:=d.db.Exec(`DELETE FROM rules WHERE id=?`,id);return err}
func(d *DB)ToggleRule(id string)error{_,err:=d.db.Exec(`UPDATE rules SET enabled=1-enabled WHERE id=?`,id);return err}
func(d *DB)Check(ruleID,clientKey string)CheckResult{r:=d.GetRule(ruleID);if r==nil{return CheckResult{Allowed:true,Key:clientKey}}
if !r.Enabled{return CheckResult{Allowed:true,Remaining:r.Limit,Limit:r.Limit,Key:clientKey}}
window:=time.Now().Add(-time.Duration(r.WindowSec)*time.Second).UTC().Format(time.RFC3339)
var count int;d.db.QueryRow(`SELECT COUNT(*) FROM hits WHERE rule_id=? AND client_key=? AND created_at>=?`,ruleID,clientKey,window).Scan(&count)
remaining:=r.Limit-count-1;if remaining<0{remaining=0}
allowed:=count<r.Limit
if allowed{d.db.Exec(`INSERT INTO hits(id,rule_id,client_key,created_at)VALUES(?,?,?,?)`,genID(),ruleID,clientKey,now())}
resetAt:=time.Now().Add(time.Duration(r.WindowSec)*time.Second).UTC().Format(time.RFC3339)
return CheckResult{Allowed:allowed,Remaining:remaining,Limit:r.Limit,ResetAt:resetAt,Key:clientKey}}
func(d *DB)Cleanup(){cutoff:=time.Now().Add(-24*time.Hour).UTC().Format(time.RFC3339);d.db.Exec(`DELETE FROM hits WHERE created_at<?`,cutoff)}
type Stats struct{Rules int `json:"rules"`;TotalHits int `json:"total_hits"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM rules`).Scan(&s.Rules);d.db.QueryRow(`SELECT COUNT(*) FROM hits`).Scan(&s.TotalHits);return s}
