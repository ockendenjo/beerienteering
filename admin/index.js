let stashFile = { demo: false, stashes: [] };
let editingId = null;

document.addEventListener('DOMContentLoaded', init);

async function init() {
  try {
    const res = await fetch(`https://${window.location.host}/stashes`);
    if (!res.ok) throw new Error(res.statusText);
    stashFile = await res.json();
  } catch {
    showStatus('Could not load stashes — starting with empty data', 'error');
  }
  render();

  const toggle = document.getElementById('show-instructions-toggle');
  toggle.checked = !!stashFile.showInstructions;
  toggle.addEventListener('change', () => {
    stashFile.showInstructions = toggle.checked;
  });

  document.getElementById('save-btn').addEventListener('click', saveStashes);
  document.getElementById('add-btn').addEventListener('click', openAdd);
  document.getElementById('cancel-btn').addEventListener('click', closeModal);
  document.getElementById('stash-form').addEventListener('submit', handleFormSubmit);
  document.getElementById('modal-overlay').addEventListener('click', e => {
    if (e.target === document.getElementById('modal-overlay')) closeModal();
  });
}

function render() {
  const tbody = document.getElementById('stashes-body');
  tbody.innerHTML = '';
  for (const stash of stashFile.stashes) {
    tbody.appendChild(makeRow(stash));
  }
}

function makeRow(stash) {
  const tr = document.createElement('tr');
  const fields = [
    stash.location,
    stash.type,
    String(stash.points),
    String(stash.lat),
    String(stash.lon),
    stash.w3w || '',
  ];
  for (const text of fields) {
    const td = document.createElement('td');
    td.textContent = text;
    tr.appendChild(td);
  }
  const td = document.createElement('td');
  const editBtn = document.createElement('button');
  editBtn.textContent = 'Edit';
  editBtn.addEventListener('click', () => openEdit(stash.id));
  const deleteBtn = document.createElement('button');
  deleteBtn.textContent = 'Delete';
  deleteBtn.className = 'delete-btn';
  deleteBtn.addEventListener('click', () => deleteStash(stash.id));
  td.append(editBtn, deleteBtn);
  tr.appendChild(td);
  return tr;
}

function openAdd() {
  editingId = null;
  document.getElementById('modal-title').textContent = 'Add Stash';
  document.getElementById('stash-form').reset();
  document.getElementById('modal-overlay').classList.remove('hidden');
}

function openEdit(id) {
  const stash = stashFile.stashes.find(s => s.id === id);
  if (!stash) return;
  editingId = id;
  document.getElementById('modal-title').textContent = 'Edit Stash';
  const form = document.getElementById('stash-form');
  form.elements['location'].value = stash.location;
  form.elements['type'].value = stash.type;
  form.elements['points'].value = stash.points;
  form.elements['lat'].value = stash.lat;
  form.elements['lon'].value = stash.lon;
  form.elements['w3w'].value = stash.w3w || '';
  document.getElementById('modal-overlay').classList.remove('hidden');
}

function closeModal() {
  document.getElementById('modal-overlay').classList.add('hidden');
  editingId = null;
}

function handleFormSubmit(e) {
  e.preventDefault();
  const form = e.target;
  const existing = editingId ? stashFile.stashes.find(s => s.id === editingId) : null;

  const stash = {
    id: editingId || crypto.randomUUID(),
    location: form.elements['location'].value.trim(),
    type: form.elements['type'].value,
    points: parseInt(form.elements['points'].value, 10),
    lat: parseFloat(form.elements['lat'].value),
    lon: parseFloat(form.elements['lon'].value),
    w3w: form.elements['w3w'].value.trim(),
    contents: existing ? (existing.contents || []) : [],
    hide: existing ? !!existing.hide : false,
  };

  if (editingId) {
    const idx = stashFile.stashes.findIndex(s => s.id === editingId);
    if (idx !== -1) stashFile.stashes[idx] = stash;
  } else {
    stashFile.stashes.push(stash);
  }

  closeModal();
  render();
}

function deleteStash(id) {
  if (!confirm('Delete this stash?')) return;
  stashFile.stashes = stashFile.stashes.filter(s => s.id !== id);
  render();
}

async function saveStashes() {
  const btn = document.getElementById('save-btn');
  btn.disabled = true;
  btn.textContent = 'Saving…';
  try {
    const res = await fetch(`https://${window.location.host}/stashes`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(stashFile),
    });
    if (res.ok) {
      showStatus('Saved successfully', 'success');
    } else {
      const text = await res.text();
      showStatus(`Save failed: ${text}`, 'error');
    }
  } catch (e) {
    showStatus(`Save failed: ${e.message}`, 'error');
  } finally {
    btn.disabled = false;
    btn.textContent = 'Save';
  }
}

let statusTimeout;

function showStatus(msg, type) {
  const el = document.getElementById('status');
  el.textContent = msg;
  el.className = type;
  clearTimeout(statusTimeout);
  statusTimeout = setTimeout(() => { el.className = ''; }, 4000);
}
