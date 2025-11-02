/**
 * 用户页面JavaScript
 */

// 用户页面功能
const UserPage = {
    currentTab: 'dashboard',
    tokens: [],

    /**
     * 初始化页面
     */
    init: function() {
        this.bindEvents();
        this.loadTokens();
        this.updateDashboard();
    },

    /**
     * 绑定事件
     */
    bindEvents: function() {
        // 标签页切换
        $('.nav-link[data-tab]').click((e) => {
            e.preventDefault();
            const tab = $(e.currentTarget).data('tab');
            this.showTab(tab);
        });

        // 表单提交
        $('#changePasswordForm').submit(this.handleChangePassword.bind(this));
        $('#generateTokenForm').submit(this.handleGenerateToken.bind(this));
    },

    /**
     * 显示指定标签页
     * @param {string} tabName - 标签页名称
     */
    showTab: function(tabName) {
        // 隐藏所有标签页内容
        $('.tab-content').hide();
        
        // 显示指定标签页
        $('#' + tabName).show();
        
        // 更新导航状态
        $('.nav-link').removeClass('active');
        $(`.nav-link[data-tab="${tabName}"]`).addClass('active');
        
        this.currentTab = tabName;
        
        // 根据标签页加载相应数据
        switch(tabName) {
            case 'dashboard':
                this.updateDashboard();
                break;
            case 'tokens':
                this.loadTokens();
                break;
        }
    },

    /**
     * 加载仪表盘数据
     */
    updateDashboard: function() {
        const now = new Date();
        const total = this.tokens.length;
        const active = this.tokens.filter(token => {
            if (!token.expires_at) return true;
            return new Date(token.expires_at) > now;
        }).length;
        const expired = total - active;

        $('#totalTokens').text(total);
        $('#activeTokens').text(active);
        $('#expiredTokens').text(expired);
        $('#revokedTokens').text(0);
        $('#tokenCount').text(`${active} 个活跃Token`);
    },

    /**
     * 加载Token列表
     */
    loadTokens: async function() {
        try {
            const response = await API.get('/user/api_tokens');

            if (response.code === 0) {
                this.tokens = Array.isArray(response.data?.tokens) ? response.data.tokens : [];
                this.renderTokensTable(this.tokens);
                this.updateDashboard();
            } else {
                Utils.showAlert('加载Token列表失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('加载Token列表失败:', error);
            Utils.showAlert('加载Token列表失败', 'danger');
        }
    },

    /**
     * 渲染Token表格
     * @param {Array} tokens - Token列表
     */
    renderTokensTable: function(tokens) {
        const tbody = $('#tokensTableBody');
        
        if (tokens.length === 0) {
            tbody.html('<tr><td colspan="5" class="text-center text-muted">暂无Token</td></tr>');
            return;
        }
        
        const html = tokens.map(token => `
            <tr>
                <td>${token.name}</td>
                <td>${Utils.formatDate(token.created_at)}</td>
                <td>${token.expires_at ? Utils.formatDate(token.expires_at) : '永久有效'}</td>
                <td>
                    <span class="badge ${this.isTokenActive(token) ? 'bg-success' : 'bg-secondary'}">
                        ${this.isTokenActive(token) ? '活跃' : '已过期'}
                    </span>
                </td>
                <td>
                    <div class="btn-group btn-group-sm">
                        <button class="btn btn-outline-danger" onclick="UserPage.deleteToken(${token.id})">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');
        
        tbody.html(html);
    },

    /**
     * 显示生成Token模态框
     */
    showGenerateTokenModal: function() {
        $('#generateTokenModal').modal('show');
        $('#generateTokenForm')[0].reset();
    },

    /**
     * 生成Token
     */
    generateToken: async function() {
        const form = $('#generateTokenForm');
        const formData = new FormData(form[0]);
        
        const data = {
            name: formData.get('tokenName'),
            expire: parseInt(formData.get('expire'), 10)
        };

        // 验证输入
        if (!data.name) {
            Utils.showAlert('请输入Token名称', 'warning');
            return;
        }

        if (!Number.isFinite(data.expire) || data.expire <= 0) {
            Utils.showAlert('请选择有效的过期时间', 'warning');
            return;
        }

        try {
            const submitSelector = '#generateTokenForm button[type="submit"]';
            Utils.showLoading(true, submitSelector);
            
            const response = await API.post('/user/api_tokens', data);

            if (response.code === 0) {
                // 显示生成的Token
                this.showGeneratedToken(response.token);
                $('#generateTokenModal').modal('hide');
                form[0].reset();
                this.loadTokens();
            } else {
                Utils.showAlert('生成Token失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('生成Token失败:', error);
            Utils.showAlert('生成Token失败', 'danger');
        } finally {
            const submitSelector = '#generateTokenForm button[type="submit"]';
            Utils.showLoading(false, submitSelector);
        }
    },

    /**
     * 显示生成的Token
     * @param {string} token - 生成的Token
     */
    showGeneratedToken: function(token) {
        $('#generatedToken').val(token);
        $('#tokenDisplayModal').modal('show');
    },

    /**
     * 复制Token
     */
    copyToken: function() {
        const token = $('#generatedToken').val();
        Utils.copyToClipboard(token);
    },

    /**
     * 下载Token
     */
    downloadToken: function() {
        const token = $('#generatedToken').val();
        const blob = new Blob([token], { type: 'text/plain' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'api-token.txt';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
    },

    /**
     * 查看Token
     * @param {number} tokenId - Token ID
     */
    // 判断 token 是否有效
    isTokenActive: function(token) {
        if (!token.expires_at) {
            return true;
        }
        return new Date(token.expires_at) > new Date();
    },

    deleteToken: async function(tokenId) {
        if (!confirm('确定要删除这个Token吗？')) {
            return;
        }

        try {
            const response = await API.delete(`/user/api_tokens/${tokenId}`);
            
            if (response.code === 0) {
                Utils.showAlert('Token删除成功', 'success');
                this.loadTokens();
            } else {
                Utils.showAlert('删除Token失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('删除Token失败:', error);
            Utils.showAlert('删除Token失败', 'danger');
        }
    },

    /**
     * 加载个人设置
     */
    handleChangePassword: async function(e) {
        e.preventDefault();
        
        const form = $('#changePasswordForm');
        const formData = new FormData(form[0]);
        
        const data = {
            old_password: formData.get('old_password'),
            new_password: formData.get('new_password')
        };

        // 验证密码确认
        const confirmPassword = $('#confirmPassword').val();
        if (data.new_password !== confirmPassword) {
            Utils.showAlert('两次输入的密码不一致', 'warning');
            return;
        }

        try {
            const submitSelector = '#changePasswordForm button[type="submit"]';
            Utils.showLoading(true, submitSelector);
            
            const response = await API.post('/user/change_password', data);
            
            if (response.code === 0) {
                Utils.showAlert('密码修改成功', 'success');
                form[0].reset();
            } else {
                Utils.showAlert('密码修改失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('密码修改失败:', error);
            Utils.showAlert('密码修改失败', 'danger');
        } finally {
            const submitSelector = '#changePasswordForm button[type="submit"]';
            Utils.showLoading(false, submitSelector);
        }
    },

    /**
     * 处理Token生成
     */
    handleGenerateToken: function(e) {
        e.preventDefault();
        this.generateToken();
    },

    /**
     * 获取Token数量
     */
    getTokenCount: function() {
        const now = new Date();
        return this.tokens.filter(token => {
            if (!token.expires_at) return true;
            return new Date(token.expires_at) > now;
        }).length;
    },

    /**
     * 登出
     */
    logout: function() {
        if (confirm('确定要退出登录吗？')) {
            API.post('/logout', {}).finally(() => {
                window.location.href = '/login';
            });
        }
    }
};

// 页面加载完成后初始化
$(document).ready(function() {
    UserPage.init();
});

// 导出到全局
window.UserPage = UserPage;
