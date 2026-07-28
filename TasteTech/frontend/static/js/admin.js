/**
 * GoFood - Admin Panel Helper JS
 */

// Sidebar toggle for mobile responsive
function toggleSidebar() {
    const sidebar = document.getElementById('admin-sidebar');
    const overlay = document.getElementById('sidebar-overlay');
    
    if (sidebar && overlay) {
        sidebar.classList.toggle('open');
        overlay.classList.toggle('open');
    }
}

// Update clock in admin header
function updateClock() {
    const timeEl = document.getElementById('header-time');
    if (!timeEl) return;
    
    const now = new Date();
    const options = { 
        weekday: 'long', 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    };
    timeEl.textContent = now.toLocaleDateString('id-ID', options);
}

// Show admin toast notification
function showToast(message, type = 'info', duration = 3000) {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    
    const icons = {
        success: '✅',
        error: '❌',
        info: 'ℹ️'
    };
    
    toast.innerHTML = `
        <span>${icons[type] || 'ℹ️'}</span>
        <span>${message}</span>
    `;

    container.appendChild(toast);

    // Auto remove
    setTimeout(() => {
        toast.classList.add('hiding');
        setTimeout(() => toast.remove(), 300);
    }, duration);
}

// Trigger success message from URL params if exists
function checkURLParams() {
    const urlParams = new URLSearchParams(window.location.search);
    const successMsg = urlParams.get('success');
    const errorMsg = urlParams.get('error');
    
    if (successMsg) {
        showToast(successMsg, 'success');
        // clean URL parameter without reloading
        window.history.replaceState({}, document.title, window.location.pathname);
    }
    if (errorMsg) {
        showToast(errorMsg, 'error');
        window.history.replaceState({}, document.title, window.location.pathname);
    }
}

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', () => {
    updateClock();
    setInterval(updateClock, 60000);
    checkURLParams();
});
