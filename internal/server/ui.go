package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Cutoff</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.header{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.header h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.content{padding:1.5rem;max-width:900px;margin:0 auto}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.6rem;padding:1rem 1.2rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.card-name{font-family:var(--mono);font-size:.82rem}
.card-detail{display:grid;grid-template-columns:repeat(auto-fit,minmax(120px,1fr));gap:.5rem;margin-top:.5rem;font-family:var(--mono);font-size:.68rem}
.detail-item{background:var(--bg);padding:.4rem .6rem;border:1px solid var(--bg3)}
.detail-label{color:var(--cm);font-size:.55rem;text-transform:uppercase;letter-spacing:1px}
.detail-val{color:var(--cream);margin-top:.1rem}
.badge{font-family:var(--mono);font-size:.55rem;padding:.1rem .4rem;text-transform:uppercase;letter-spacing:1px}
.badge-enabled{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}
.badge-disabled{background:var(--bg3);color:var(--cm);border:1px solid var(--bg3)}
.badge-reject{background:#c9444422;color:var(--red)}.badge-throttle{background:#d4843a22;color:var(--orange)}.badge-log{background:var(--bg3);color:var(--cm)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-primary:hover{opacity:.85}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.test-box{background:#0d0b09;border:1px solid var(--bg3);padding:1rem;margin-top:1rem;font-family:var(--mono);font-size:.75rem}
.test-result{margin-top:.5rem;padding:.5rem;border:1px solid var(--bg3)}
.test-allowed{border-color:var(--green);color:var(--green)}.test-blocked{border-color:var(--red);color:var(--red)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.form-row{margin-bottom:.6rem}
.form-row label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.form-row input,.form-row select{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem}
.actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:2rem;color:var(--cm);font-style:italic}
</style></head><body>
<div class="header"><h1>CUTOFF</h1><button class="btn btn-primary" onclick="openForm()">+ New Rule</button></div>
<div class="content" id="main"></div>
<div class="modal-bg" id="modalBg" onclick="if(event.target===this)closeModal()"><div class="modal" id="modal"></div></div>

<script>
const API='/api';let rules=[];

async function load(){
  const r=await fetch(API+'/rules').then(r=>r.json());
  rules=r.rules||[];render();
}

function render(){
  const m=document.getElementById('main');
  if(!rules||!rules.length){m.innerHTML='<div class="empty">No rate limit rules yet. Create one to start limiting.</div><div class="test-box"><div style="color:var(--leather);font-size:.6rem;margin-bottom:.5rem;text-transform:uppercase;letter-spacing:1px">Test endpoint</div><code>GET /api/check/{rule_id}?key=user123</code><br><br>Returns whether the request is allowed and remaining quota.</div>';return;}
  let h='';
  (rules||[]).forEach(r=>{
    h+='<div class="card"><div class="card-top"><div><span class="card-name">'+esc(r.name)+'</span> <span class="badge badge-'+(r.enabled?'enabled':'disabled')+'">'+(r.enabled?'enabled':'disabled')+'</span></div><div style="display:flex;gap:.3rem"><button class="btn btn-sm" onclick="testRule(\''+r.id+'\')">Test</button><button class="btn btn-sm" onclick="delRule(\''+r.id+'\')" style="color:var(--red)">Delete</button></div></div>';
    h+='<div class="card-detail"><div class="detail-item"><div class="detail-label">Pattern</div><div class="detail-val">'+esc(r.key||'*')+'</div></div><div class="detail-item"><div class="detail-label">Limit</div><div class="detail-val">'+r.limit+' req</div></div><div class="detail-item"><div class="detail-label">Window</div><div class="detail-val">'+r.window_sec+'s</div></div><div class="detail-item"><div class="detail-label">Action</div><div class="detail-val"><span class="badge badge-'+r.action+'">'+r.action+'</span></div></div></div></div>';
  });
  h+='<div class="test-box" id="testArea"><div style="color:var(--leather);font-size:.6rem;margin-bottom:.5rem;text-transform:uppercase;letter-spacing:1px">Test a rule</div>Click "Test" on any rule above, or call: <code>GET /api/check/{rule_id}?key=user123</code></div>';
  m.innerHTML=h;
}

async function testRule(id){
  const key='test-user-'+Math.random().toString(36).slice(2,6);
  const r=await fetch(API+'/check/'+id+'?key='+key).then(r=>r.json());
  const area=document.getElementById('testArea');
  area.innerHTML='<div style="color:var(--leather);font-size:.6rem;margin-bottom:.5rem;text-transform:uppercase;letter-spacing:1px">Test Result</div><div class="test-result '+(r.allowed?'test-allowed':'test-blocked')+'"><strong>'+(r.allowed?'ALLOWED':'BLOCKED')+'</strong><br>Key: '+esc(r.key||key)+'<br>Remaining: '+r.remaining+'/'+r.limit+'<br>Resets: '+r.reset_at+'</div>';
}

async function delRule(id){if(confirm('Delete?')){await fetch(API+'/rules/'+id,{method:'DELETE'});load();}}

function openForm(){
  document.getElementById('modal').innerHTML='<h2>New Rate Limit Rule</h2><div class="form-row"><label>Name</label><input id="f-name" placeholder="e.g. API requests per minute"></div><div class="form-row"><label>Key pattern</label><input id="f-key" placeholder="e.g. api-* or * for all" value="*"></div><div class="form-row"><label>Limit (requests)</label><input id="f-limit" type="number" value="100"></div><div class="form-row"><label>Window (seconds)</label><input id="f-window" type="number" value="60"></div><div class="form-row"><label>Action when exceeded</label><select id="f-action"><option value="reject">Reject (429)</option><option value="throttle">Throttle (delay)</option><option value="log">Log only</option></select></div><div class="actions"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-primary" onclick="submitRule()">Create</button></div>';
  document.getElementById('modalBg').classList.add('open');
}
async function submitRule(){
  await fetch(API+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-name').value,key:document.getElementById('f-key').value,limit:parseInt(document.getElementById('f-limit').value),window_sec:parseInt(document.getElementById('f-window').value),action:document.getElementById('f-action').value,enabled:true})});
  closeModal();load();
}
function closeModal(){document.getElementById('modalBg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
