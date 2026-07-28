/**
 * GoFood - Main JavaScript
 * Handles: Cart interactions, UI animations, Toast notifications
 */

// =============================================
// CART MANAGEMENT
// =============================================

/**
 * Add item to cart via API
 * @param {number} menuId - Menu ID
 * @param {string} menuName - Menu name for notification
 * @param {number} qty - Quantity (default: 1)
 */
async function addToCart(menuId, menuName, qty = 1) {
    try {
        const response = await fetch('/api/cart/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ menu_id: menuId, quantity: qty })
        });

        const data = await response.json();

        if (data.success) {
            // Update cart count in navbar
            updateCartCount(data.cart_count);
            // Animate cart button
            animateCartBtn();
            // Show success toast
            showToast(`🛒 ${menuName} ditambahkan ke keranjang!`, 'success');
        } else {
            showToast(data.message || 'Gagal menambahkan ke keranjang', 'error');
        }
    } catch (err) {
        showToast('Koneksi bermasalah. Coba lagi.', 'error');
        console.error('Add to cart error:', err);
    }
}

/**
 * Update cart item quantity
 * @param {number} menuId - Menu ID
 * @param {number} newQty - New quantity
 */
async function updateCartItem(menuId, newQty) {
    try {
        const response = await fetch('/api/cart/update', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ menu_id: menuId, quantity: newQty })
        });

        const data = await response.json();

        if (data.success) {
            if (newQty <= 0) {
                // Remove item from DOM
                const itemEl = document.getElementById(`cart-item-${menuId}`);
                if (itemEl) {
                    itemEl.style.animation = 'fadeOut 0.3s ease forwards';
                    setTimeout(() => {
                        itemEl.remove();
                        checkEmptyCart();
                    }, 300);
                }
                showToast('Item dihapus dari keranjang', 'info');
            } else {
                // Update qty display
                const qtyEl = document.getElementById(`cart-qty-${menuId}`);
                if (qtyEl) qtyEl.textContent = newQty;

                // Update subtotal
                const subtotalEl = document.getElementById(`subtotal-${menuId}`);
                if (subtotalEl && data.cart) {
                    const item = data.cart.items?.find(i => i.menu_id === menuId);
                    if (item) {
                        subtotalEl.textContent = formatRupiah(item.sub_total);
                    }
                }
            }

            // Update total
            updateCartTotal(data.cart_total);
            updateCartCount(data.cart_count);
        }
    } catch (err) {
        console.error('Update cart error:', err);
        showToast('Gagal mengupdate keranjang', 'error');
    }
}

/**
 * Remove item from cart
 * @param {number} menuId - Menu ID to remove
 */
async function removeFromCart(menuId) {
    try {
        const response = await fetch(`/api/cart/remove/${menuId}`, {
            method: 'DELETE'
        });

        const data = await response.json();

        if (data.success) {
            const itemEl = document.getElementById(`cart-item-${menuId}`);
            if (itemEl) {
                itemEl.style.opacity = '0';
                itemEl.style.transform = 'translateX(-20px)';
                itemEl.style.transition = 'all 0.3s ease';
                setTimeout(() => {
                    itemEl.remove();
                    checkEmptyCart();
                }, 300);
            }
            updateCartTotal(data.cart_total);
            updateCartCount(data.cart_count);
            showToast('Item dihapus dari keranjang', 'info');
        }
    } catch (err) {
        console.error('Remove from cart error:', err);
    }
}

/**
 * Update cart count badge in navbar
 */
function updateCartCount(count) {
    const cartCountEl = document.getElementById('cart-count');
    if (cartCountEl) {
        cartCountEl.textContent = count || 0;
        cartCountEl.classList.add('bounce');
        setTimeout(() => cartCountEl.classList.remove('bounce'), 500);
    }
}

/**
 * Update cart total display
 */
function updateCartTotal(total) {
    const subtotalEl = document.getElementById('summary-subtotal');
    const totalEl = document.getElementById('summary-total');

    if (subtotalEl) subtotalEl.textContent = formatRupiah(total);
    if (totalEl) {
        const delivery = 5000;
        totalEl.textContent = formatRupiah(total + delivery);
    }
}

/**
 * Check if cart is empty and redirect if needed
 */
function checkEmptyCart() {
    const items = document.querySelectorAll('.cart-item');
    if (items.length === 0) {
        setTimeout(() => location.reload(), 500);
    }
}

/**
 * Animate cart button
 */
function animateCartBtn() {
    const btn = document.getElementById('cart-btn');
    if (btn) {
        btn.style.transform = 'scale(1.15)';
        setTimeout(() => btn.style.transform = '', 300);
    }
}

// =============================================
// CART PAGE - FETCH CART DATA ON LOAD
// =============================================
async function fetchCartCount() {
    try {
        const response = await fetch('/api/cart');
        const data = await response.json();
        if (data.success) {
            updateCartCount(data.cart?.item_count || 0);
        }
    } catch (err) {
        // Silent fail
    }
}

// =============================================
// FORMATTING UTILITIES
// =============================================

/**
 * Format number to Indonesian Rupiah format
 */
function formatRupiah(num) {
    if (num === undefined || num === null) return 'Rp 0';
    return 'Rp ' + Math.floor(num).toLocaleString('id-ID');
}

// =============================================
// TOAST NOTIFICATIONS
// =============================================

/**
 * Show toast notification
 * @param {string} message - Message to display
 * @param {string} type - 'success' | 'error' | 'info'
 * @param {number} duration - Duration in ms (default: 3000)
 */
function showToast(message, type = 'info', duration = 3000) {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const icons = {
        success: '✅',
        error: '❌',
        info: 'ℹ️'
    };

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `
        <span class="toast-icon">${icons[type] || '📢'}</span>
        <span>${message}</span>
    `;

    container.appendChild(toast);

    // Auto remove
    setTimeout(() => {
        toast.classList.add('hiding');
        setTimeout(() => toast.remove(), 300);
    }, duration);

    // Click to dismiss
    toast.addEventListener('click', () => {
        toast.classList.add('hiding');
        setTimeout(() => toast.remove(), 300);
    });
}

// =============================================
// NAVBAR
// =============================================

// Navbar scroll effect
const navbar = document.getElementById('navbar');
if (navbar) {
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            navbar.classList.add('scrolled');
        } else {
            navbar.classList.remove('scrolled');
        }
    });
}

// Mobile hamburger menu
const hamburger = document.getElementById('hamburger');
const navLinks = document.getElementById('nav-links');
if (hamburger && navLinks) {
    hamburger.addEventListener('click', () => {
        navLinks.classList.toggle('mobile-open');
        const isOpen = navLinks.classList.contains('mobile-open');
        hamburger.setAttribute('aria-expanded', isOpen);

        // Animate hamburger
        const spans = hamburger.querySelectorAll('span');
        if (isOpen) {
            spans[0].style.transform = 'rotate(45deg) translate(5px, 5px)';
            spans[1].style.opacity = '0';
            spans[2].style.transform = 'rotate(-45deg) translate(5px, -5px)';
        } else {
            spans[0].style.transform = '';
            spans[1].style.opacity = '';
            spans[2].style.transform = '';
        }
    });

    // Close on outside click
    document.addEventListener('click', (e) => {
        if (!hamburger.contains(e.target) && !navLinks.contains(e.target)) {
            navLinks.classList.remove('mobile-open');
        }
    });
}

// =============================================
// HERO PARTICLES
// =============================================
function createParticles() {
    const container = document.getElementById('particles');
    if (!container) return;

    const colors = ['rgba(255, 107, 53, 0.4)', 'rgba(255, 215, 0, 0.3)', 'rgba(255, 140, 90, 0.3)'];
    const emojis = ['🌟', '✨', '🔥'];

    for (let i = 0; i < 15; i++) {
        const particle = document.createElement('div');
        particle.className = 'particle';

        const size = Math.random() * 8 + 4;
        const isEmoji = Math.random() > 0.7;

        if (isEmoji) {
            particle.textContent = emojis[Math.floor(Math.random() * emojis.length)];
            particle.style.fontSize = `${size * 2}px`;
        } else {
            particle.style.width = `${size}px`;
            particle.style.height = `${size}px`;
            particle.style.background = colors[Math.floor(Math.random() * colors.length)];
        }

        particle.style.left = `${Math.random() * 100}%`;
        particle.style.top = `${Math.random() * 100}%`;
        particle.style.animationDelay = `${Math.random() * 6}s`;
        particle.style.animationDuration = `${Math.random() * 4 + 5}s`;

        container.appendChild(particle);
    }
}

// =============================================
// SCROLL ANIMATIONS
// =============================================
function initScrollAnimations() {
    const animateEls = document.querySelectorAll(
        '.menu-card, .category-card, .step-card, .promo-card, .cart-item, .form-card, .summary-card'
    );

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.animation = 'fadeInUp 0.5s ease forwards';
                observer.unobserve(entry.target);
            }
        });
    }, { threshold: 0.1 });

    animateEls.forEach((el, index) => {
        el.style.opacity = '0';
        el.style.animationDelay = `${index * 0.05}s`;
        observer.observe(el);
    });
}

// =============================================
// INIT ON DOM LOADED
// =============================================
document.addEventListener('DOMContentLoaded', () => {
    // Fetch cart count on page load
    fetchCartCount();

    // Create particles
    createParticles();

    // Init scroll animations
    setTimeout(initScrollAnimations, 100);

    // Smooth scroll for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(link => {
        link.addEventListener('click', (e) => {
            const target = document.querySelector(link.getAttribute('href'));
            if (target) {
                e.preventDefault();
                target.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }
        });
    });

    // Close mobile nav on link click
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', () => {
            const navLinks = document.getElementById('nav-links');
            if (navLinks) navLinks.classList.remove('mobile-open');
        });
    });
});
