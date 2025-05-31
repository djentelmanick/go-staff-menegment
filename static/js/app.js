// Global Variables
let authToken = '';
let currentStaffId = null;
let allStaff = [];
const API_BASE = '/api';

// DOM Elements
const loginSection = document.getElementById('loginSection');
const mainPanel = document.getElementById('mainPanel');
const loginForm = document.getElementById('loginForm');
const staffForm = document.getElementById('staffForm');
const searchInput = document.getElementById('searchInput');
const staffTableBody = document.getElementById('staffTableBody');
const modal = document.getElementById('staffModal');
const closeModalBtn = document.getElementById('closeModalBtn');
const cancelModalBtn = document.getElementById('cancelModalBtn');
const addStaffBtn = document.getElementById('addStaffBtn');
const logoutBtn = document.getElementById('logoutBtn');

// Event Listeners
document.addEventListener('DOMContentLoaded', initApp);
loginForm.addEventListener('submit', handleLogin);
staffForm.addEventListener('submit', handleStaffSubmit);
searchInput.addEventListener('input', handleSearch);
closeModalBtn.addEventListener('click', closeModal);
cancelModalBtn.addEventListener('click', closeModal);
addStaffBtn.addEventListener('click', () => openModal('add'));
logoutBtn.addEventListener('click', logout);

// Initialize Application
function initApp() {
    const savedToken = localStorage.getItem('authToken');
    if (savedToken) {
        authToken = savedToken;
        showMainPanel();
        loadStaff();
    }
}

// Auth Functions
async function handleLogin(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const loginData = {
        login: formData.get('login'),
        password: formData.get('password')
    };

    try {
        const response = await fetch(`${API_BASE}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData)
        });

        const result = await response.json();

        if (result.success) {
            authToken = result.token;
            localStorage.setItem('authToken', authToken);
            showNotification('Успешная авторизация!', 'success');
            showMainPanel();
            loadStaff();
        } else {
            showNotification(result.message, 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

function logout() {
    authToken = '';
    localStorage.removeItem('authToken');
    showLoginSection();
    loginForm.reset();
}

// UI Functions
function showMainPanel() {
    loginSection.style.display = 'none';
    mainPanel.style.display = 'block';
}

function showLoginSection() {
    loginSection.style.display = 'block';
    mainPanel.style.display = 'none';
}

function openModal(action) {
    currentStaffId = null;
    document.getElementById('modalTitle').textContent = action === 'edit' ? 'Редактировать сотрудника' : 'Добавить сотрудника';
    staffForm.reset();
    modal.style.display = 'block';
}

function closeModal() {
    modal.style.display = 'none';
    currentStaffId = null;
}

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

// Staff CRUD Operations
async function loadStaff() {
    try {
        const response = await fetch(`${API_BASE}/staff`, {
            headers: {
                'Authorization': authToken
            }
        });

        if (response.status === 401) {
            logout();
            return;
        }

        const staff = await response.json();
        allStaff = staff || [];
        displayStaff(allStaff);
    } catch (error) {
        showNotification('Ошибка загрузки данных', 'error');
    }
}

function displayStaff(staff) {
    staffTableBody.innerHTML = '';

    staff.forEach(person => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${person.id}</td>
            <td>${person.full_name}</td>
            <td>${person.phone || '-'}</td>
            <td>${person.email || '-'}</td>
            <td>${person.address || '-'}</td>
            <td class="action-btns">
                <button class="btn btn-sm btn-secondary" onclick="editStaff(${person.id})">Редактировать</button>
                <button class="btn btn-sm btn-danger" onclick="deleteStaff(${person.id})">Удалить</button>
            </td>
        `;
        staffTableBody.appendChild(row);
    });
}

async function handleStaffSubmit(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const staffData = {
        full_name: formData.get('fullName'),
        phone: formData.get('phone'),
        email: formData.get('email'),
        address: formData.get('address')
    };

    const url = currentStaffId ? 
        `${API_BASE}/staff/${currentStaffId}` : 
        `${API_BASE}/staff`;
    
    const method = currentStaffId ? 'PUT' : 'POST';

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken
            },
            body: JSON.stringify(staffData)
        });

        if (response.ok) {
            const message = currentStaffId ? 'Сотрудник обновлен!' : 'Сотрудник добавлен!';
            showNotification(message, 'success');
            closeModal();
            loadStaff();
        } else {
            showNotification('Ошибка сохранения данных', 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

function editStaff(id) {
    const staff = allStaff.find(s => s.id === id);
    if (!staff) return;

    currentStaffId = id;
    document.getElementById('modalTitle').textContent = 'Редактировать сотрудника';
    document.getElementById('fullName').value = staff.full_name;
    document.getElementById('phone').value = staff.phone || '';
    document.getElementById('email').value = staff.email || '';
    document.getElementById('address').value = staff.address || '';
    
    openModal('edit');
}

async function deleteStaff(id) {
    if (!confirm('Вы уверены, что хотите удалить этого сотрудника?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/staff/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': authToken
            }
        });

        if (response.ok) {
            showNotification('Сотрудник удален!', 'success');
            loadStaff();
        } else {
            showNotification('Ошибка удаления сотрудника', 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

// Search Functionality
function handleSearch(e) {
    const searchTerm = e.target.value.toLowerCase();
    const filteredStaff = allStaff.filter(staff => 
        staff.full_name.toLowerCase().includes(searchTerm) ||
        (staff.phone && staff.phone.toLowerCase().includes(searchTerm)) ||
        (staff.email && staff.email.toLowerCase().includes(searchTerm)) ||
        (staff.address && staff.address.toLowerCase().includes(searchTerm))
    );
    displayStaff(filteredStaff);
}

// Close modal when clicking outside
window.addEventListener('click', (e) => {
    if (e.target === modal) {
        closeModal();
    }
});

// Global functions for inline event handlers
window.editStaff = editStaff;
window.deleteStaff = deleteStaff;