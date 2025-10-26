/**
 * 登录页面JavaScript
 */

// 登录页面功能
const LoginPage = {
    /**
     * 初始化页面
     */
    init: function() {
        this.bindEvents();
        this.checkAuth();
    },

    /**
     * 绑定事件
     */
    bindEvents: function() {
        // 登录表单提交
        $('#loginForm').submit(this.handleLogin.bind(this));
        
        // 输入框回车事件
        $('#username, #password').keypress((e) => {
            if (e.which === 13) {
                $('#loginForm').submit();
            }
        });
        
        // 表单验证
        $('#loginForm input').on('blur', this.validateField.bind(this));
    },

    /**
     * 检查认证状态
     */
    checkAuth: function() {
        // 如果已经有token，直接跳转
        if (Utils.getToken()) {
            this.redirectToHome();
        }
    },

    /**
     * 处理登录
     * @param {Event} e - 事件对象
     */
    handleLogin: async function(e) {
        e.preventDefault();
        
        const formData = new FormData(e.target);
        const data = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        // 验证输入
        if (!this.validateForm(data)) {
            return;
        }

        // 显示加载状态
        Utils.showLoading(true);
        this.hideAlert();

        try {
            const result = await API.post('/auth/login', data);
            
            if (result.code === 0) {
                this.handleLoginSuccess(result);
            } else {
                this.handleLoginError(result.message);
            }
        } catch (error) {
            console.error('Login error:', error);
            this.handleLoginError('登录失败，请检查网络连接');
        } finally {
            Utils.showLoading(false);
        }
    },

    /**
     * 处理登录成功
     * @param {Object} result - 登录结果
     */
    handleLoginSuccess: function(result) {
        Utils.showAlert('登录成功，正在跳转...', 'success');
        
        // 保存token
        if (result.token) {
            Utils.setToken(result.token);
        }
        
        // 根据角色跳转
        setTimeout(() => {
            this.redirectToHome(result.redirect);
        }, 1000);
    },

    /**
     * 处理登录错误
     * @param {string} message - 错误信息
     */
    handleLoginError: function(message) {
        Utils.showAlert('登录失败: ' + message, 'danger');
    },

    /**
     * 验证表单
     * @param {Object} data - 表单数据
     */
    validateForm: function(data) {
        let isValid = true;
        
        if (!data.username || data.username.trim() === '') {
            this.showFieldError('username', '请输入用户名');
            isValid = false;
        } else {
            this.clearFieldError('username');
        }
        
        if (!data.password || data.password.trim() === '') {
            this.showFieldError('password', '请输入密码');
            isValid = false;
        } else {
            this.clearFieldError('password');
        }
        
        return isValid;
    },

    /**
     * 验证单个字段
     * @param {Event} e - 事件对象
     */
    validateField: function(e) {
        const field = e.target;
        const value = field.value.trim();
        const fieldName = field.name;
        
        if (fieldName === 'username') {
            if (value === '') {
                this.showFieldError(fieldName, '请输入用户名');
            } else {
                this.clearFieldError(fieldName);
            }
        } else if (fieldName === 'password') {
            if (value === '') {
                this.showFieldError(fieldName, '请输入密码');
            } else {
                this.clearFieldError(fieldName);
            }
        }
    },

    /**
     * 显示字段错误
     * @param {string} fieldName - 字段名
     * @param {string} message - 错误信息
     */
    showFieldError: function(fieldName, message) {
        const field = $(`[name="${fieldName}"]`);
        field.addClass('is-invalid');
        
        // 移除旧的错误信息
        field.siblings('.invalid-feedback').remove();
        
        // 添加新的错误信息
        field.after(`<div class="invalid-feedback">${message}</div>`);
    },

    /**
     * 清除字段错误
     * @param {string} fieldName - 字段名
     */
    clearFieldError: function(fieldName) {
        const field = $(`[name="${fieldName}"]`);
        field.removeClass('is-invalid');
        field.siblings('.invalid-feedback').remove();
    },

    /**
     * 隐藏提示信息
     */
    hideAlert: function() {
        $('.alert').remove();
    },

    /**
     * 跳转到首页
     * @param {string} redirect - 重定向地址
     */
    redirectToHome: function(redirect) {
        if (redirect) {
            window.location.href = redirect;
        } else {
            window.location.href = '/';
        }
    }
};

// 页面加载完成后初始化
$(document).ready(function() {
    LoginPage.init();
});

// 导出到全局
window.LoginPage = LoginPage;