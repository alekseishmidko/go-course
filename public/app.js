'use strict';
const $ = s => document.querySelector(s);
const state = {tasks: [], users: [], filter: 'all', view: 'tasks', editor: null, loading: false};
const escapeHTML = value => String(value ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
const date = value => new Date(value).toLocaleDateString('ru-RU', {day:'numeric',month:'short'});
const userName = id => state.users.find(u => u.id === id)?.full_name || `Автор #${id}`;
const initials = name => name.split(/\s+/).slice(0,2).map(x => x[0] || '').join('').toUpperCase();
let noticeTimer;
function notify(message, success = false) { clearTimeout(noticeTimer); const n = $('#notice'); n.textContent = message; n.className = 'notice' + (success ? ' success' : ''); n.hidden = false; if(success) noticeTimer = setTimeout(() => n.hidden = true, 4500); }
async function api(path, options = {}) {
 const response = await fetch('/api/v1' + path, {...options, headers:{Accept:'application/json', ...(options.body ? {'Content-Type':'application/json'} : {}), ...options.headers}});
 if (!response.ok) { const body = await response.json().catch(() => ({})); throw new Error(body.error || body.message || `Ошибка HTTP ${response.status}`); }
 return response.status === 204 ? null : response.json();
}
function authorOptions() { return state.users.map(u => `<option value="${u.id}">${escapeHTML(u.full_name)} · #${u.id}</option>`).join(''); }
async function load() {
 state.loading = true; $('#refresh').disabled = true;
 try {
  const [users,tasks] = await Promise.all([api('/users'),api('/tasks')]);
  state.users = users || []; state.tasks = tasks || [];
  for(const selector of ['#author-filter','#stats-author']) { const el=$(selector), value=el.value; el.innerHTML='<option value="">Все авторы</option>'+authorOptions(); if([...el.options].some(o=>o.value===value)) el.value=value; }
  renderTasks(); renderUsers();
 } catch(error) { notify('Не удалось загрузить данные: ' + error.message); if(!state.tasks.length) $('#task-list').innerHTML='<div class="empty"><b>Нет соединения с API</b>Проверь, что приложение и база данных запущены. Нажми «Обновить», чтобы повторить.</div>'; }
 finally { state.loading=false; $('#refresh').disabled=false; }
}
function renderTasks() {
 const author = $('#author-filter').value;
 const scoped = state.tasks.filter(t => !author || t.author_user_id === Number(author));
 const done = scoped.filter(t => t.completed).length;
 const percentage = scoped.length ? Math.round(done/scoped.length*100) : 0;
 $('#total').textContent=scoped.length; $('#pending').textContent=scoped.length-done; $('#done').textContent=done;
 $('#progress-value').textContent=percentage+'%'; $('#progress-bar').style.width=percentage+'%'; $('#nav-count').textContent=state.tasks.length;
 const query=$('#search').value.trim().toLocaleLowerCase();
 const tasks=scoped.filter(t => (state.filter==='all' || t.completed===(state.filter==='done')) && `${t.title} ${t.description || ''}`.toLocaleLowerCase().includes(query));
 $('#task-count').textContent=`Показано ${tasks.length} из ${scoped.length}`;
 $('#task-list').innerHTML=tasks.length ? tasks.map(t => `<article class="task-row ${t.completed?'completed':''}"><input class="task-check" type="checkbox" data-toggle="${t.id}" ${t.completed?'checked':''} aria-label="${escapeHTML((t.completed?'Вернуть в работу: ':'Завершить: ')+t.title)}"><div class="task-info"><div class="task-title">${escapeHTML(t.title)}</div><div class="task-description">${escapeHTML(t.description || 'Без описания')}</div></div><div class="task-meta"><span class="avatar">${escapeHTML(initials(userName(t.author_user_id)))}</span><div>${escapeHTML(userName(t.author_user_id))}<small>${date(t.created_at)}</small></div></div><div><button class="icon-button" data-edit-task="${t.id}" aria-label="Редактировать задачу">✎</button><button class="icon-button" data-delete-task="${t.id}" aria-label="Удалить задачу">×</button></div></article>`).join('') : '<div class="empty"><b>Здесь пока нет задач</b>Создай новую задачу или измени фильтры.</div>';
}
function renderUsers() { $('#user-list').innerHTML=state.users.length ? state.users.map(u=>`<article class="user-card"><span class="avatar">${escapeHTML(initials(u.full_name))}</span><h3>${escapeHTML(u.full_name)}</h3><p>${escapeHTML(u.phone_number || 'Телефон не указан')} · #${u.id}</p><footer><button class="text-button" data-edit-user="${u.id}">Изменить</button><button class="icon-button" data-delete-user="${u.id}" aria-label="Удалить пользователя">×</button></footer></article>`).join('') : '<div class="empty"><b>Начни с первого пользователя</b>Затем можно будет добавить его задачи.</div>'; }
function field(name,label,value='',extra='') { return `<label>${label}<input name="${name}" value="${escapeHTML(value)}" ${extra}></label>`; }
function openEditor(type,id=null) {
 if(type==='task' && !state.users.length) { notify('Сначала добавь пользователя — автора задачи.'); return; }
 const item=id ? state[type==='task'?'tasks':'users'].find(x=>x.id===id) : null;
 state.editor={type,id,item}; $('#editor-title').textContent=(id?'Редактировать ':'Добавить ')+(type==='task'?'задачу':'пользователя'); $('#form-error').textContent='';
 if(type==='task') {
  $('#editor-fields').innerHTML=field('title','Название',item?.title,'required minlength="1" maxlength="100"')+field('description','Описание (необязательно)',item?.description,'maxlength="1000"')+(!id?`<label>Автор<select name="author_user_id" required>${authorOptions()}</select></label>`:`<label><input name="completed" type="checkbox" ${item.completed?'checked':''}> Задача завершена</label><p class="hint">Текущий API не позволяет очистить уже заданное описание.</p>`);
 } else $('#editor-fields').innerHTML=field('full_name','Полное имя',item?.full_name,'required minlength="3" maxlength="100"')+field('phone_number','Телефон (необязательно)',item?.phone_number,'type="tel" pattern="\\+.{9,14}" placeholder="+79997654321"');
 $('#save').disabled=false; $('#editor').showModal();
}
async function save(event) {
 event.preventDefault(); const {type,id,item}=state.editor, data=new FormData(event.target); let body;
 if(type==='task') {
  const description=data.get('description').trim(); body={title:data.get('title').trim()};
  if(!body.title) { $('#form-error').textContent='Название не должно быть пустым.'; return; }
  if(id) { body.completed=data.has('completed'); if(description) body.description=description; else if(item.description) { $('#form-error').textContent='Описание нельзя очистить: введи текст или отмени изменения.'; return; } }
  else { body.author_user_id=Number(data.get('author_user_id')); if(description) body.description=description; }
 } else { body={full_name:data.get('full_name').trim(),phone_number:data.get('phone_number').trim() || null}; if([...body.full_name].length<3) { $('#form-error').textContent='Имя должно содержать минимум 3 символа.'; return; } }
 $('#save').disabled=true;
 try { await api('/'+(type==='task'?'tasks':'users')+(id?'/'+id:''),{method:id?'PATCH':'POST',body:JSON.stringify(body)}); $('#editor').close(); notify('Изменения сохранены',true); await load(); }
 catch(error) { $('#form-error').textContent=error.message; } finally { $('#save').disabled=false; }
}
async function remove(type,id) {
 const dialog=$('#confirm-dialog'); dialog.returnValue=''; dialog.showModal();
 dialog.addEventListener('close',async()=>{ if(dialog.returnValue!=='delete') return; try { await api(`/${type}/${id}`,{method:'DELETE'}); notify('Запись удалена',true); await load(); } catch(error) { notify('Не удалось удалить запись: '+error.message); } },{once:true});
}
async function loadStats(event) {
 event?.preventDefault(); const params=new URLSearchParams();
 if($('#stats-author').value) params.set('userId',$('#stats-author').value);
 if($('#stats-from').value) params.set('from',$('#stats-from').value);
 if($('#stats-to').value) params.set('to',$('#stats-to').value);
 if(params.has('from')&&params.has('to')&&params.get('from')>=params.get('to')) { notify('Конец периода должен быть позже начала.'); return; }
 const button=$('#stats-form button'); button.disabled=true; $('#stats-result').textContent='Загружаем статистику…';
 try { const s=await api('/statistics?'+params); const metrics=[['Всего задач',s.tasks_created],['Завершено',s.tasks_completed],['Процент выполнения',s.tasks_completed_percent==null?'—':Number(s.tasks_completed_percent).toFixed(1)+'%'],['Среднее время',s.task_complete_average_time || '—']]; $('#stats-result').innerHTML=metrics.map(([label,value])=>`<article><span>${label}</span><strong>${escapeHTML(value)}</strong></article>`).join(''); }
 catch(error) { $('#stats-result').textContent='Статистика недоступна'; notify(error.message); } finally { button.disabled=false; }
}
document.addEventListener('click',async event=>{
 const nav=event.target.closest('[data-view]'); if(nav) { state.view=nav.dataset.view; document.querySelectorAll('.view').forEach(v=>v.hidden=v.id!==state.view+'-view'); document.querySelectorAll('.nav').forEach(n=>n.classList.toggle('active',n===nav)); $('#crumb').textContent=nav.textContent.trim().replace(/\d+$/,''); if(state.view==='statistics') loadStats(); }
 const filter=event.target.closest('[data-filter]'); if(filter) { state.filter=filter.dataset.filter; document.querySelectorAll('[data-filter]').forEach(b=>b.classList.toggle('selected',b===filter)); renderTasks(); }
 for(const [attr,type] of [['editTask','task'],['editUser','user']]) { const el=event.target.closest(type==='task'?'[data-edit-task]':'[data-edit-user]'); if(el) openEditor(type,Number(el.dataset[attr])); }
 for(const [attr,type] of [['deleteTask','tasks'],['deleteUser','users']]) { const el=event.target.closest(type==='tasks'?'[data-delete-task]':'[data-delete-user]'); if(el) remove(type,Number(el.dataset[attr])); }
});
$('#task-list').addEventListener('change',async event=>{ const el=event.target.closest('[data-toggle]'); if(!el) return; el.disabled=true; try { await api('/tasks/'+el.dataset.toggle,{method:'PATCH',body:JSON.stringify({completed:el.checked})}); await load(); } catch(error) { el.checked=!el.checked; notify(error.message); } finally { el.disabled=false; } });
$('#new-task').onclick=()=>openEditor('task'); $('#new-user').onclick=()=>openEditor('user');
$('#close-editor').onclick=$('#cancel-editor').onclick=()=>$('#editor').close();
$('#editor-form').addEventListener('submit',save); $('#stats-form').addEventListener('submit',loadStats);
$('#search').addEventListener('input',renderTasks); $('#author-filter').addEventListener('change',renderTasks); $('#refresh').onclick=load;
$('#today').textContent=new Date().toLocaleDateString('ru-RU',{weekday:'long',day:'numeric',month:'long'});
$('#task-list').innerHTML='<div class="empty">Загружаем задачи…</div>';
load();
