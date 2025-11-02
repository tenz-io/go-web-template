/**
 * 通用JavaScript工具函数
 */

// 全局配置
const AppConfig = {
    apiBaseUrl: '',
    tokenKey: 'user_token',
    cookieName: 'jwt_token'
};

// 工具函数
const Utils = {
    /**
     * 显示提示信息
     * @param {string} message - 提示信息
     * @param {string} type - 提示类型 (success, danger, warning, info)
     * @param {number} duration - 显示时长(毫秒)
     */
    showAlert: function(message, type = 'info', duration = 3000) {
        const alertHtml = `
            <div class="alert alert-${type} alert-dismissible fade show" role="alert">
                <i class="fas fa-${this.getAlertIcon(type)}"></i> ${message}
                <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
            </div>
        `;
        
        // 移除现有的提示
        $('.alert').remove();
        
        // 添加新提示到页面顶部
        $('.main-content, .login-body').first().prepend(alertHtml);
        
        // 自动消失
        if (duration > 0) {
            setTimeout(() => {
                $('.alert').fadeOut();
            }, duration);
        }
    },

    /**
     * 获取提示图标
     */
    getAlertIcon: function(type) {
        const icons = {
            success: 'check-circle',
            danger: 'exclamation-triangle',
            warning: 'exclamation-circle',
            info: 'info-circle'
        };
        return icons[type] || 'info-circle';
    },

    /**
     * 显示加载状态
     * @param {boolean} show - 是否显示加载
     * @param {string} selector - 按钮选择器
     */
    showLoading: function(show, selector = 'button[type="submit"]') {
        if (show) {
            $(selector).prop('disabled', true);
            $(selector).find('.btn-text').hide();
            $(selector).find('.loading').show();
        } else {
            $(selector).prop('disabled', false);
            $(selector).find('.btn-text').show();
            $(selector).find('.loading').hide();
        }
    },

    /**
     * 获取Token
     */
    getToken: function() {
        const token = localStorage.getItem(AppConfig.tokenKey);
        return token ? token.trim() : '';
    },

    /**
     * 设置Token
     */
    setToken: function(token) {
        localStorage.setItem(AppConfig.tokenKey, token);
    },

    /**
     * 清除Token
     */
    clearToken: function() {
        localStorage.removeItem(AppConfig.tokenKey);
    },

    /**
     * 复制文本到剪贴板
     * @param {string} text - 要复制的文本
     */
    copyToClipboard: function(text) {
        if (navigator.clipboard) {
            navigator.clipboard.writeText(text).then(() => {
                this.showAlert('已复制到剪贴板', 'success');
            }).catch(() => {
                this.showAlert('复制失败，请手动复制', 'warning');
            });
        } else {
            // 降级方案
            const textArea = document.createElement('textarea');
            textArea.value = text;
            document.body.appendChild(textArea);
            textArea.select();
            try {
                document.execCommand('copy');
                this.showAlert('已复制到剪贴板', 'success');
            } catch (err) {
                this.showAlert('复制失败，请手动复制', 'warning');
            }
            document.body.removeChild(textArea);
        }
    },

    /**
     * 格式化日期
     * @param {string|Date} date - 日期
     * @param {string} format - 格式
     */
    formatDate: function(date, format = 'YYYY-MM-DD HH:mm:ss') {
        const d = new Date(date);
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        const hours = String(d.getHours()).padStart(2, '0');
        const minutes = String(d.getMinutes()).padStart(2, '0');
        const seconds = String(d.getSeconds()).padStart(2, '0');
        
        return format
            .replace('YYYY', year)
            .replace('MM', month)
            .replace('DD', day)
            .replace('HH', hours)
            .replace('mm', minutes)
            .replace('ss', seconds);
    },

    /**
     * 防抖函数
     * @param {Function} func - 要防抖的函数
     * @param {number} wait - 等待时间
     */
    debounce: function(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },

    /**
     * 节流函数
     * @param {Function} func - 要节流的函数
     * @param {number} limit - 限制时间
     */
    throttle: function(func, limit) {
        let inThrottle;
        return function() {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    },

    /**
     * 验证密码强度
     * @param {string} password - 密码
     */
    validatePassword: function(password) {
        const minLength = 6;
        if (password.length < minLength) {
            return { valid: false, message: `密码长度至少${minLength}位` };
        }
        return { valid: true, message: '密码格式正确' };
    }
};

// API请求封装
const API = {
    /**
     * 发送GET请求
     * @param {string} url - 请求URL
     * @param {Object} headers - 请求头
     */
    get: async function(url, headers = {}) {
        const defaultHeaders = {
            'Content-Type': 'application/json'
        };
        
        if (Utils.getToken()) {
            defaultHeaders['Authorization'] = `Bearer ${Utils.getToken()}`;
        }
        
        const response = await fetch(url, {
            method: 'GET',
            headers: { ...defaultHeaders, ...headers },
            credentials: 'include'
        });
        
        return await response.json();
    },

    /**
     * 发送POST请求
     * @param {string} url - 请求URL
     * @param {Object} data - 请求数据
     * @param {Object} headers - 请求头
     */
    post: async function(url, data = {}, headers = {}) {
        const defaultHeaders = {
            'Content-Type': 'application/json'
        };
        
        if (Utils.getToken()) {
            defaultHeaders['Authorization'] = `Bearer ${Utils.getToken()}`;
        }
        
        const response = await fetch(url, {
            method: 'POST',
            headers: { ...defaultHeaders, ...headers },
            body: JSON.stringify(data),
            credentials: 'include'
        });
        
        return await response.json();
    },

    /**
     * 发送DELETE请求
     * @param {string} url - 请求URL
     * @param {Object} data - 请求数据
     * @param {Object} headers - 请求头
     */
    delete: async function(url, data = {}, headers = {}) {
        const defaultHeaders = {
            'Content-Type': 'application/json'
        };
        
        if (Utils.getToken()) {
            defaultHeaders['Authorization'] = `Bearer ${Utils.getToken()}`;
        }
        
        const response = await fetch(url, {
            method: 'DELETE',
            headers: { ...defaultHeaders, ...headers },
            body: JSON.stringify(data),
            credentials: 'include'
        });
        
        return await response.json();
    }
};

// 页面初始化
$(document).ready(function() {
    // 初始化工具提示
    $('[data-bs-toggle="tooltip"]').tooltip();
    
    // 初始化弹出框
    $('[data-bs-toggle="popover"]').popover();
    
    // 自动隐藏提示
    $('.alert').each(function() {
        const $alert = $(this);
        if (!$alert.hasClass('alert-permanent')) {
            setTimeout(() => {
                $alert.fadeOut();
            }, 5000);
        }
    });
    
    // 表单验证
    $('form').on('submit', function(e) {
        const $form = $(this);
        const $submitBtn = $form.find('button[type="submit"]');
        
        // 显示加载状态
        Utils.showLoading(true, $submitBtn);
        
        // 3秒后自动隐藏加载状态（防止请求失败时按钮一直处于加载状态）
        setTimeout(() => {
            Utils.showLoading(false, $submitBtn);
        }, 3000);
    });
});

// 全局错误处理
window.addEventListener('error', function(e) {
    console.error('Global error:', e.error);
    Utils.showAlert('发生未知错误，请刷新页面重试', 'danger');
});

// 未处理的Promise拒绝
window.addEventListener('unhandledrejection', function(e) {
    console.error('Unhandled promise rejection:', e.reason);
    Utils.showAlert('网络请求失败，请检查网络连接', 'danger');
});

// 导出到全局
window.Utils = Utils;
window.API = API;
window.AppConfig = AppConfig;
