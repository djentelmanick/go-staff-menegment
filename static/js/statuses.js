// Global Variables
let authToken = localStorage.getItem('authToken');
let allEmployees = [];
const API_BASE = '/api';

// DOM Elements
const searchInput = document.getElementById('searchInput');
const statusTableBody = document.getElementById('statusTableBody');

// Event Listeners
document.addEventListener('DOMContentLoaded', initStatuses);
searchInput.addEventListener('input', handleSearch);

// Initialize Statuses Page
function initStatuses() {
    if (!authToken) {
        window.location.href = '/login';
        return;
    }
    
    loadEmployees();
    
    // Add event listeners to all status selects
    document.addEventListener('change', function(e) {
        if (e.target.classList.contains('status-select')) {
            updateEmployeeStatus(e.target);
        }
    });
}

// Employee Status Operations
async function loadEmployees() {
    try {
        const response = await fetch(`${API_BASE}/employees/statuses`, {
            headers: {
                'Authorization': authToken
            }
        });

        if (response.status === 401) {
            window.location.href = '/login';
            return;
        }

        const employees = await response.json();
        allEmployees = employees || [];
        displayEmployees(allEmployees);
    } catch (error) {
        showNotification('Ошибка загрузки данных', 'error');
    }
}

function displayEmployees(employees) {
    statusTableBody.innerHTML = '';

    const sortedEmployees = [...employees].sort((a, b) => {
        if (a.status < b.status) return -1;
        if (a.status > b.status) return 1;
        
        return a.full_name.localeCompare(b.full_name);
    });

    sortedEmployees.forEach(employee => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${employee.id}</td>
            <td>${employee.full_name}</td>
            <td>
                <select class="status-select" data-employee-id="${employee.id}">
                    <option value="">-</option>
                    <option value="working" ${employee.status === 'working' ? 'selected' : ''}>На работе</option>
                    <option value="vacation" ${employee.status === 'vacation' ? 'selected' : ''}>В отпуске</option>
                    <option value="consideration" ${employee.status === 'consideration' ? 'selected' : ''}>На рассмотрении</option>
                </select>
            </td>
        `;
        statusTableBody.appendChild(row);
    });
}

async function updateEmployeeStatus(selectElement) {
    const employeeId = selectElement.dataset.employeeId;
    const newStatus = selectElement.value;
    
    if (!employeeId) return;

    try {
        const response = await fetch(`${API_BASE}/employees/${employeeId}/status`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken
            },
            body: JSON.stringify({ status: newStatus })
        });

        if (response.ok) {
            showNotification('Статус сотрудника обновлен!', 'success');
            // Update local data
            const employee = allEmployees.find(e => e.id == employeeId);
            if (employee) {
                employee.status = newStatus;
            }
        } else {
            showNotification('Ошибка обновления статуса', 'error');
            // Reset to previous value
            selectElement.value = allEmployees.find(e => e.id == employeeId)?.status || '';
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
        // Reset to previous value
        selectElement.value = allEmployees.find(e => e.id == employeeId)?.status || '';
    }
}

// Search Functionality
function handleSearch(e) {
    const searchTerm = e.target.value.toLowerCase();
    
    if (searchTerm === '') {
        displayEmployees(allEmployees);
        return;
    }

    const filteredEmployees = allEmployees.filter(employee => 
        employee.full_name.toLowerCase().includes(searchTerm)
    );
    
    displayEmployees(filteredEmployees);
}

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
