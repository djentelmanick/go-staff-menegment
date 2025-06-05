// Global Variables
let authToken = '';

// DOM Elements
const loginForm = document.getElementById('loginForm');

// Event Listeners
document.addEventListener('DOMContentLoaded', checkAuth);
loginForm.addEventListener('submit', handleLogin);

// Auth Functions
async function checkAuth() {
    const savedToken = localStorage.getItem('authToken');
    if (savedToken) {
        // Если есть сохраненный токен, перенаправляем на dashboard
        window.location.href = '/dashboard';
    }
}

async function handleLogin(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const loginData = {
        login: formData.get('login'),
        password: formData.get('password')
    };

    try {
        const response = await fetch('/api/login', {
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
            window.location.href = '/dashboard';
        } else {
            showNotification(result.message, 'error');
        }
    } catch (error) {
        showNotification('Ошибка подключения к серверу', 'error');
    }
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