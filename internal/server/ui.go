package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Cutoff</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.ct{padding:1.5rem;max-width:900px;margin:0 auto}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.8rem;padding:1rem 1.2rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.rule-name{font-family:var(--mono);font-size:.85rem}
.rule-config{font-family:var(--mono);font-size:.7rem;color:var(--cd);margin-top:.3rem;display:flex;gap:1.5rem}
.rule-config span{color:var(--cm)}
.rule-config strong{color:var(--cream)}
.toggle{position:relative;width:36px;height:18px;cursor:pointer;display:inline-block;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-p:hover{opacity:.85}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.test-box{background:#0d0b09;border:1px solid var(--bg3);padding:1rem;margin-top:1.5rem}
.test-result{font-family:var(--mono);font-size:.75rem;margin-top:.5rem;padding:.5rem;border:1px solid var(--bg3)}
.test-allow{color:var(--green);border-color:#4a9e5c44}.test-deny{color:var(--red);border-color:#c9444444}
.empty{text-align:center;padding:2rem;color:var(--cm);font-style:italic}
</style></head><body>
<div class="hdr"><h1>CUTOFF</h1><span style="font-family:var(--mono);font-size:.7rem;color:var(--cm)" id="st"></span></div>
<div class="ct">
<div style="display:flex;justify-content:space-between;margin-bottom:1rem"><span style="font-family:var(--mono);font-size:.75rem;color:var(--leather)">RATE LIMIT RULES</span><button class="btn btn-p" onclick="oRule()">+ New Rule</button></div>
<div id="rules"></div>
<div class="test-box"><div style="font-family:var(--mono);font-size:.7rem;color:var(--leather);margin-bottom:.5rem">TEST A RULE</div>
<div style="display:flex;gap:.5rem;align-items:flex-end"><div class="fr" style="flex:1;margin:0"><label>Rule</label><select id="tr"></select></div><div class="fr" style="flex:1;margin:0"><label>Client Key</label><input id="tk" value="test-client"></div><button class="btn btn-p" onclick="testRule()">Check</button></div>
<div id="testResult"></div></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let rules=[];
async function ld(){const[r,s]=await Promise.all([fetch(A+'/rules').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);rules=r.rules||[];
document.getElementById('st').textContent=s.rules+' rules, '+s.total_hits+' total checks';rn();}
function rn(){
  const m=document.getElementById('rules');
  let sel='';
  if(!rules.length){m.innerHTML='<div class="empty">No rules yet. Create your first rate limit rule.</div>';} else {
  let h='';rules.forEach(r=>{
    h+='<div class="card"><div class="card-top"><div><span class="rule-name">'+esc(r.name)+'</span> <label class="toggle"><input type="checkbox" '+(r.enabled?'checked':'')+' onchange="tog(\''+r.id+'\')"><span class="sl"></span></label></div><button class="btn" onclick="del(\''+r.id+'\')" style="font-size:.55rem;color:var(--red)">Delete</button></div>';
    h+='<div class="rule-config"><span>Key: <strong>'+esc(r.key)+'</strong></span><span>Limit: <strong>'+r.limit+' / '+r.window_sec+'s</strong></span><span>Action: <strong>'+r.action+'</strong></span></div></div>';
    sel+='<option value="'+r.id+'">'+esc(r.name)+'</option>';
  });m.innerHTML=h;}
  document.getElementById('tr').innerHTML='<option value="">Select rule</option>'+sel;
}
async function tog(id){await fetch(A+'/rules/'+id+'/toggle',{method:'PATCH'});ld();}
async function del(id){if(confirm('Delete rule and all hits?')){await fetch(A+'/rules/'+id,{method:'DELETE'});ld();}}
function oRule(){document.getElementById('mdl').innerHTML='<h2>New Rate Limit Rule</h2><div class="fr"><label>Name</label><input id="rn" placeholder="e.g. API Rate Limit"></div><div class="fr"><label>Key Pattern</label><input id="rk" value="*" placeholder="* or specific key"></div><div class="fr"><label>Limit (requests)</label><input id="rl" type="number" value="100"></div><div class="fr"><label>Window (seconds)</label><input id="rw" type="number" value="60"></div><div class="fr"><label>Action when exceeded</label><select id="ra"><option value="reject">Reject (429)</option><option value="throttle">Throttle (delay)</option><option value="log">Log only</option></select></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sRule()">Create</button></div>';document.getElementById('mbg').classList.add('open');}
async function sRule(){await fetch(A+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('rn').value,key:document.getElementById('rk').value,limit:parseInt(document.getElementById('rl').value),window_sec:parseInt(document.getElementById('rw').value),action:document.getElementById('ra').value})});cm();ld();}
async function testRule(){
  const ruleId=document.getElementById('tr').value;const key=document.getElementById('tk').value;
  if(!ruleId){document.getElementById('testResult').innerHTML='<div class="test-result" style="color:var(--cm)">Select a rule first</div>';return;}
  const r=await fetch(A+'/check',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({rule_id:ruleId,key:key})}).then(r=>r.json());
  document.getElementById('testResult').innerHTML='<div class="test-result '+(r.allowed?'test-allow':'test-deny')+'">'+(r.allowed?'✓ ALLOWED':'✗ RATE LIMITED')+' — '+r.remaining+'/'+r.limit+' remaining</div>';
}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
ld();
</script></body></html>`
