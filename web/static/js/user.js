/**
 * 用户页面JavaScript
 */

// 用户页面功能
const UserPage = {
    currentTab: 'dashboard',
    
    /**
     * 初始化页面
     */
    init: function() {
        this.bindEvents();
        this.loadUserInfo();
        this.loadDashboard();
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
                this.loadDashboard();
                break;
            case 'tokens':
                this.loadTokens();
                break;
            case 'profile':
                this.loadProfile();
                break;
        }
    },

    /**
     * 加载用户信息
     */
    loadUserInfo: async function() {
        try {
            // 从token中解析用户信息，或者从API获取
            const token = Utils.getToken();
            if (token) {
                // 这里应该解析JWT token获取用户信息
                $('#userInfo').text('当前用户');
                $('#userRole').text('普通用户');
                $('#userCreated').text(Utils.formatDate(new Date()));
            }
        } catch (error) {
            console.error('加载用户信息失败:', error);
        }
    },

    /**
     * 加载仪表盘数据
     */
    loadDashboard: async function() {
        try {
            // 加载token统计
            const tokenCount = await this.getTokenCount();
            $('#tokenCount').text(`${tokenCount} 个活跃Token`);
        } catch (error) {
            console.error('加载仪表盘数据失败:', error);
        }
    },

    /**
     * 加载Token列表
     */
    loadTokens: async function() {
        try {
            // 模拟数据，实际应该从API获取
            const tokens = [
                {
                    id: 1,
                    name: 'API访问Token',
                    created_at: new Date(),
                    expires_at: new Date(Date.now() + 24 * 60 * 60 * 1000),
                    status: 'active'
                }
            ];
            
            this.renderTokensTable(tokens);
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
                <td>${Utils.formatDate(token.expires_at)}</td>
                <td>
                    <span class="badge ${token.status === 'active' ? 'bg-success' : 'bg-secondary'}">
                        ${token.status === 'active' ? '活跃' : '已过期'}
                    </span>
                </td>
                <td>
                    <div class="btn-group btn-group-sm">
                        <button class="btn btn-outline-primary" onclick="UserPage.viewToken(${token.id})">
                            <i class="fas fa-eye"></i>
                        </button>
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
            expire: parseInt(formData.get('expire'))
        };

        // 验证输入
        if (!data.name) {
            Utils.showAlert('请输入Token名称', 'warning');
            return;
        }

        try {
            Utils.showLoading(true);
            
            const response = await API.post('/user/generate_token', data);
            
            if (response.code === 0) {
                // 显示生成的Token
                this.showGeneratedToken(response.token);
                $('#generateTokenModal').modal('hide');
                this.loadTokens();
            } else {
                Utils.showAlert('生成Token失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('生成Token失败:', error);
            Utils.showAlert('生成Token失败', 'danger');
        } finally {
            Utils.showLoading(false);
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
    viewToken: function(tokenId) {
        Utils.showAlert('查看Token功能开发中...', 'info');
    },

    /**
     * 删除Token
     * @param {number} tokenId - Token ID
     */
    deleteToken: async function(tokenId) {
        if (!confirm('确定要删除这个Token吗？')) {
            return;
        }

        try {
            const response = await API.delete('/user/delete_token', { token_id: tokenId });
            
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
    loadProfile: function() {
        // 加载个人设置数据
        console.log('加载个人设置');
    },

    /**
     * 处理密码修改
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
            Utils.showLoading(true);
            
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
            Utils.showLoading(false);
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
    getTokenCount: async function() {
        try {
            // 模拟API调用
            return 1;
        } catch (error) {
            console.error('获取Token数量失败:', error);
            return 0;
        }
    },

    /**
     * 登出
     */
    logout: function() {
        if (confirm('确定要退出登录吗？')) {
            // 清除token
            Utils.clearToken();
            // 跳转到登录页
            window.location.href = '/login';
        }
    }
};

// 页面加载完成后初始化
$(document).ready(function() {
    UserPage.init();
});

// 导出到全局
window.UserPage = UserPage;