package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Cutoff</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.main{padding:1.5rem;max-width:900px;margin:0 auto}
.stats{display:grid;grid-template-columns:1fr 1fr;gap:.8rem;margin-bottom:1.2rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem;text-align:center}
.st-v{font-size:1.4rem}.st-l{font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.1rem}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.6rem;padding:1rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.rule-name{font-size:.85rem}
.rule-config{font-size:.65rem;color:var(--cd);margin-top:.3rem;display:flex;gap:1rem;flex-wrap:wrap}
.rule-config span{background:var(--bg);padding:.15rem .4rem;border:1px solid var(--bg3)}
.toggle{position:relative;width:36px;height:18px;cursor:pointer;display:inline-block;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}
.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.test-box{background:#0d0b09;border:1px solid var(--bg3);padding:.8rem;margin-top:.5rem;font-size:.7rem}
.test-result{margin-top:.3rem;padding:.3rem .5rem}.test-pass{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}.test-fail{background:#c9444422;color:var(--red);border:1px solid #c9444444}
.btn{font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>CUTOFF</h1><button class="btn btn-p" onclick="openForm()">+ New Rule</button></div>
<div class="main" id="main"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let rules=[];
async function load(){
  const[r,s]=await Promise.all([fetch(A+'/rules').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
  rules=r.rules||[];render(s);
}
function render(s){
  const m=document.getElementById('main');
  let h='<div class="stats"><div class="st"><div class="st-v">'+s.rules+'</div><div class="st-l">Rules</div></div><div class="st"><div class="st-v">'+s.total_hits+'</div><div class="st-l">Total Hits</div></div></div>';
  if(!rules.length){h+='<div class="empty">No rate limit rules. Create one to start limiting.</div>';m.innerHTML=h;return;}
  rules.forEach(r=>{
    let win=r.window_sec<60?r.window_sec+'s':r.window_sec<3600?(r.window_sec/60)+'m':(r.window_sec/3600)+'h';
    h+='<div class="card"><div class="card-top"><div class="rule-name">'+esc(r.name)+'</div><div style="display:flex;gap:.5rem;align-items:center"><label class="toggle"><input type="checkbox" '+(r.enabled?'checked':'')+' onchange="tog(\''+r.id+'\')"><span class="sl"></span></label><button class="btn" onclick="del(\''+r.id+'\')" style="color:var(--red)">✕</button></div></div>';
    h+='<div class="rule-config"><span>Pattern: '+esc(r.key)+'</span><span>Limit: '+r.limit+' / '+win+'</span><span>Action: '+r.action+'</span></div>';
    h+='<div class="test-box">Test: <input id="test-'+r.id+'" value="test-user" style="width:120px;padding:2px 4px;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem"> <button class="btn" onclick="test(\''+r.id+'\')">Check</button><div id="result-'+r.id+'"></div></div>';
    h+='</div>';
  });
  m.innerHTML=h;
}
async function tog(id){await fetch(A+'/rules/'+id+'/toggle',{method:'POST'});load();}
async function del(id){if(confirm('Delete?')){await fetch(A+'/rules/'+id,{method:'DELETE'});load();}}
async function test(id){
  const key=document.getElementById('test-'+id).value;
  const r=await fetch(A+'/check/'+id+'?key='+encodeURIComponent(key)).then(r=>r.json());
  const el=document.getElementById('result-'+id);
  el.innerHTML='<div class="test-result '+(r.allowed?'test-pass':'test-fail')+'">'+(r.allowed?'✓ ALLOWED':'✗ RATE LIMITED')+' — '+r.remaining+'/'+r.limit+' remaining</div>';
}
function openForm(){
  document.getElementById('mdl').innerHTML='<h2>New Rate Limit Rule</h2><div class="fr"><label>Name</label><input id="f-n" placeholder="e.g. API requests"></div><div class="fr"><label>Key Pattern</label><input id="f-k" value="*" placeholder="* matches all"></div><div class="fr"><label>Requests per window</label><input id="f-l" type="number" value="100"></div><div class="fr"><label>Window (seconds)</label><input id="f-w" type="number" value="60"></div><div class="fr"><label>Action when exceeded</label><select id="f-a"><option value="reject">Reject (429)</option><option value="throttle">Throttle (delay)</option><option value="log">Log only</option></select></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Create</button></div>';
  document.getElementById('mbg').classList.add('open');
}
async function sub(){await fetch(A+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-n').value,key:document.getElementById('f-k').value,limit:parseInt(document.getElementById('f-l').value),window_sec:parseInt(document.getElementById('f-w').value),action:document.getElementById('f-a').value,enabled:true})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
