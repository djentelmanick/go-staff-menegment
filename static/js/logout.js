(function() {
    const logoutBtn = document.getElementById('logoutBtn');
    
    if (!logoutBtn) return;

    logoutBtn.addEventListener('click', function() {
        localStorage.removeItem('authToken');
        window.location.href = '/login';
    });
})();