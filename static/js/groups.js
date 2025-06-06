// Global Variables
let authToken = localStorage.getItem('authToken');
let currentGroupId = null;
let allGroups = [];
let allStaff = [];
let currentGroupMembers = [];
const API_BASE = '/api';

// DOM Elements
const groupForm = document.getElementById('groupForm');
const groupSearch = document.getElementById('groupSearch');
const groupsList = document.querySelector('.groups-list');
const membersList = document.querySelector('.members-list');
const staffSelect = document.getElementById('staffSelect');
const addMemberBtn = document.querySelector('.add-members button');
const groupModal = document.getElementById('groupModal');
const closeModalBtn = document.querySelector('#groupModal .close');
const cancelModalBtn = document.querySelector('#groupModal .btn-secondary');
const addGroupBtn = document.getElementById('addGroupBtn');

// Event Listeners
document.addEventListener('DOMContentLoaded', initGroups);
groupForm.addEventListener('submit', handleGroupSubmit);
groupSearch.addEventListener('input', handleSearch);
closeModalBtn.addEventListener('click', closeModal);
cancelModalBtn.addEventListener('click', closeModal);
addGroupBtn.addEventListener('click', () => openModal('add'));
addMemberBtn.addEventListener('click', addMembersToGroup);

// Initialize Groups Dashboard
function initGroups() {
    if (!authToken) {
        window.location.href = '/login';
        return;
    }
    
    loadGroups();
    loadAllStaff();
}

// Modal Functions
function openModal(action, group = null) {
    currentGroupId = group ? group.id : null;
    document.querySelector('#groupModal h2').textContent = 
        action === 'edit' ? 'Редактировать группу' : 'Создать новую группу';
    
    if (group) {
        document.getElementById('groupName').value = group.name;
        document.getElementById('groupDescription').value = group.description || '';
    } else {
        groupForm.reset();
    }
    
    groupModal.style.display = 'flex';
}

function closeModal() {
    groupModal.style.display = 'none';
    currentGroupId = null;
}

// Group CRUD Operations
async function loadGroups() {
    try {
        const response = await fetch(`${API_BASE}/groups`, {
            headers: {
                'Authorization': authToken
            }
        });

        if (response.status === 401) {
            window.location.href = '/login';
            return;
        }

        const groups = await response.json();
        allGroups = groups || [];
        displayGroups(allGroups);
        
        if (currentGroupId) {
            loadGroupMembers(currentGroupId);
        }
        else if (allGroups.length > 0) {
            loadGroupMembers(allGroups[0].id);
        }
    } catch (error) {
        showNotification('Ошибка загрузки групп', 'error');
    }
}

async function loadAllStaff() {
    try {
        const response = await fetch(`${API_BASE}/staff`, {
            headers: {
                'Authorization': authToken
            }
        });

        if (!response.ok) return;

        const staff = await response.json();
        allStaff = staff || [];
        populateStaffSelect();
    } catch (error) {
        console.error('Error loading staff:', error);
    }
}

function populateStaffSelect() {
    staffSelect.innerHTML = '';
    allStaff.forEach(staff => {
        const option = document.createElement('option');
        option.value = staff.id;
        option.textContent = `${staff.full_name} (${staff.email || 'нет email'})`;
        staffSelect.appendChild(option);
    });
}

function displayGroups(groups) {
    groupsList.innerHTML = '';

    groups.forEach(group => {
        const groupCard = document.createElement('div');
        groupCard.className = 'group-card';
        groupCard.innerHTML = `
            <div class="group-header">
                <h3>${group.name}</h3>
                <span class="group-meta">${group.member_count || 0} сотрудников</span>
            </div>
            <div class="group-actions">
                <button class="btn btn-sm btn-primary" onclick="editGroup(${group.id})">Редактировать</button>
                <button class="btn btn-sm btn-danger" onclick="deleteGroup(${group.id})">Удалить</button>
            </div>
        `;
        
        groupCard.addEventListener('click', () => loadGroupMembers(group.id));
        groupsList.appendChild(groupCard);
    });
}

async function loadGroupMembers(groupId) {
    try {
        const response = await fetch(`${API_BASE}/groups/${groupId}/members`, {
            headers: {
                'Authorization': authToken
            }
        });

        if (!response.ok) return;

        const members = await response.json();
        currentGroupMembers = members || [];
        currentGroupId = groupId;
        displayGroupMembers(currentGroupMembers);
    } catch (error) {
        showNotification('Ошибка загрузки участников группы', 'error');
    }
}

function displayGroupMembers(members) {
    const membersHeader = document.querySelector('.members-header');
    const groupDescriptionClass = document.querySelector('.group-description');
    const currentGroup = allGroups.find(group => group.id === currentGroupId);
    const groupName = currentGroup ? currentGroup.name : '';
    const groupDescription = currentGroup ? currentGroup.description : null;
    
    membersHeader.innerHTML = `Участники группы: ${groupName}`;
    groupDescriptionClass.innerHTML = `Описание: ${groupDescription || 'отсутствует'}`;
    
    membersList.innerHTML = '';

    if (members.length === 0) {
        membersList.innerHTML = '<p>В группе нет участников</p>';
        return;
    }

    members.forEach(member => {
        const memberItem = document.createElement('div');
        memberItem.className = 'member-item';
        memberItem.innerHTML = `
            <span>${member.full_name}</span>
            <button class="btn btn-sm btn-danger" onclick="removeMember(${member.id})">Удалить</button>
        `;
        membersList.appendChild(memberItem);
    });
}

async function handleGroupSubmit(e) {
    e.preventDefault();
    
    const groupData = {
        name: document.getElementById('groupName').value,
        description: document.getElementById('groupDescription').value
    };

    const url = currentGroupId ? 
        `${API_BASE}/groups/${currentGroupId}` : 
        `${API_BASE}/groups`;
    
    const method = currentGroupId ? 'PUT' : 'POST';

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken
            },
            body: JSON.stringify(groupData)
        });

        if (response.ok) {
            const message = currentGroupId ? 'Группа обновлена!' : 'Группа создана!';
            showNotification(message, 'success');
            closeModal();
            loadGroups();
        } else {
            showNotification('Ошибка сохранения группы', 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

async function addMembersToGroup() {
    if (!currentGroupId) return;
    
    const selectedOptions = Array.from(staffSelect.selectedOptions).map(option => option.value);
    const staffIds = Array.from(staffSelect.selectedOptions).map(option => parseInt(option.value)).filter(id => !isNaN(id));
    
    if (staffIds.length === 0) return;

    try {
        const response = await fetch(`${API_BASE}/groups/${currentGroupId}/members`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken
            },
            body: JSON.stringify({ staff_ids: staffIds })
        });

        if (response.ok) {
            showNotification('Сотрудники добавлены в группу', 'success');
            loadGroupMembers(currentGroupId);
            loadGroups(); // Refresh group count
        } else {
            showNotification('Ошибка добавления сотрудников', 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

async function removeMember(staffId) {
    if (!currentGroupId || !confirm('Удалить сотрудника из группы?')) return;

    try {
        const response = await fetch(`${API_BASE}/groups/${currentGroupId}/members/${staffId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': authToken
            }
        });

        if (response.ok) {
            showNotification('Сотрудник удален из группы', 'success');
            loadGroupMembers(currentGroupId);
            loadGroups();
        } else {
            showNotification('Ошибка удаления сотрудника', 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

function editGroup(id) {
    const group = allGroups.find(g => g.id === id);
    if (!group) return;

    openModal('edit', group);
}

async function deleteGroup(id) {
    if (!confirm('Вы уверены, что хотите удалить эту группу?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/groups/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': authToken
            }
        });

        if (response.ok) {
            showNotification('Группа удалена!', 'success');
            
            if (currentGroupId === id) {
                currentGroupId = null;
                membersList.innerHTML = '<p>Выберите группу для просмотра участников</p>';
            }
            
            loadGroups();
        } else {
            showNotification('Ошибка удаления группы', 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

// Search Functionality
function handleSearch(e) {
    const searchTerm = e.target.value.toLowerCase();
    
    if (searchTerm === '') {
        displayGroups(allGroups);
        return;
    }

    const filteredGroups = allGroups.filter(group => 
        group.name.toLowerCase().includes(searchTerm) ||
        (group.description && group.description.toLowerCase().includes(searchTerm))
    );
    
    displayGroups(filteredGroups);
}

// Close modal when clicking outside
window.addEventListener('click', (e) => {
    if (e.target === groupModal) {
        closeModal();
    }
});

// Notification function
function showNotification(message, type) {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    document.body.appendChild(notification);
    
    setTimeout(() => notification.classList.add('show'), 100);
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => document.body.removeChild(notification), 300);
    }, 3000);
}

// Global functions for inline event handlers
window.editGroup = editGroup;
window.deleteGroup = deleteGroup;
window.removeMember = removeMember;