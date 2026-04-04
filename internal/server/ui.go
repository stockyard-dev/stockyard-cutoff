package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Cutoff</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.hdr-stats{font-family:var(--mono);font-size:.7rem;color:var(--cm)}
.wrap{padding:1.5rem;max-width:800px;margin:0 auto}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.6rem;padding:1rem 1.2rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.card-name{font-family:var(--mono);font-size:.82rem}
.card-limits{font-family:var(--mono);font-size:.72rem;color:var(--gold)}
.card-meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.3rem;display:flex;gap:1rem}
.toggle{position:relative;width:36px;height:18px;cursor:pointer;display:inline-block;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}
.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.test-box{background:#0d0b09;border:1px solid var(--bg3);padding:.8rem;margin-top:.5rem;font-family:var(--mono);font-size:.72rem}
.test-result{margin-top:.3rem;font-size:.7rem}
.test-result.allowed{color:var(--green)}.test-result.denied{color:var(--red)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem}
.actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic}
</style></head><body>
<div class="hdr"><h1>CUTOFF</h1><div class="hdr-stats" id="stats"></div></div>
<div class="wrap">
<div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:1rem">
<h2 style="font-family:var(--mono);font-size:.75rem;color:var(--leather)">RATE LIMIT RULES</h2>
<button class="btn btn-primary" onclick="openForm()">+ New Rule</button>
</div>
<div id="main"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let rules=[];
async function load(){const[r,s]=await Promise.all([fetch(A+'/rules').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
rules=r.rules||[];document.getElementById('stats').textContent=s.rules+' rules · '+s.total_hits+' total hits';render();}
function render(){const m=document.getElementById('main');
if(!rules||!rules.length){m.innerHTML='<div class="empty">No rate limit rules. Create one to start throttling.</div>';return;}
let h='';rules.forEach(r=>{
const win=r.window_sec>=3600?(r.window_sec/3600)+'h':r.window_sec>=60?(r.window_sec/60)+'m':r.window_sec+'s';
h+='<div class="card"><div class="card-top"><div><label class="toggle"><input type="checkbox" '+(r.enabled?'checked':'')+' onchange="toggle(\''+r.id+'\')"><span class="sl"></span></label> <span class="card-name">'+esc(r.name)+'</span></div><span class="card-limits">'+r.limit+' req / '+win+'</span></div>';
h+='<div class="card-meta"><span>Pattern: <code>'+esc(r.key)+'</code></span><span>Action: '+r.action+'</span></div>';
h+='<div class="test-box"><div style="display:flex;gap:.5rem;align-items:center"><span style="color:var(--cm)">Test:</span><input id="test-'+r.id+'" placeholder="client key" style="flex:1;padding:.2rem .4rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem"><button class="btn" onclick="testRule(\''+r.id+'\')" style="font-size:.55rem">Check</button></div><div id="result-'+r.id+'" class="test-result"></div></div>';
h+='<div style="display:flex;gap:.3rem;margin-top:.5rem"><button class="btn" onclick="del(\''+r.id+'\')" style="font-size:.55rem;color:var(--red)">Delete</button></div>';
h+='</div>';});
m.innerHTML=h;}
async function toggle(id){await fetch(A+'/rules/'+id+'/toggle',{method:'POST'});load();}
async function testRule(id){const key=document.getElementById('test-'+id).value||'test-user';
const r=await fetch(A+'/check?rule_id='+id+'&key='+encodeURIComponent(key)).then(r=>r.json());
const el=document.getElementById('result-'+id);
el.className='test-result '+(r.allowed?'allowed':'denied');
el.textContent=r.allowed?'✓ Allowed — '+r.remaining+' remaining of '+r.limit:'✗ Rate limited — resets '+new Date(r.reset_at).toLocaleTimeString();}
async function del(id){if(confirm('Delete?')){await fetch(A+'/rules/'+id,{method:'DELETE'});load();}}
function openForm(){document.getElementById('mdl').innerHTML='<h2>New Rate Limit Rule</h2><div class="fr"><label>Name</label><input id="f-name" placeholder="e.g. API requests"></div><div class="fr"><label>Key Pattern</label><input id="f-key" placeholder="e.g. api:* or user:123" value="*"></div><div class="fr"><label>Limit (requests)</label><input id="f-limit" type="number" value="100"></div><div class="fr"><label>Window (seconds)</label><input id="f-window" type="number" value="60"></div><div class="fr"><label>Action</label><select id="f-action"><option value="reject">Reject (429)</option><option value="throttle">Throttle (delay)</option><option value="log">Log only</option></select></div><div class="actions"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-primary" onclick="submitRule()">Create</button></div>';
document.getElementById('mbg').classList.add('open');}
async function submitRule(){await fetch(A+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-name').value,key:document.getElementById('f-key').value,limit:parseInt(document.getElementById('f-limit').value),window_sec:parseInt(document.getElementById('f-window').value),action:document.getElementById('f-action').value,enabled:true})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
