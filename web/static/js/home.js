/**
 * 统一Home页面JavaScript
 */

// Home页面功能
const HomePage = {
    currentTab: 'dashboard',
    userRole: 'user', // 从token中获取
    userInfo: null,
    
    /**
     * 初始化页面
     */
    init: function() {
        this.bindEvents();
        this.loadUserInfo();
        this.setupMenu();
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
        $('#addUserForm').submit(this.handleAddUser.bind(this));
        $('#generateTokenForm').submit(this.handleGenerateToken.bind(this));
        $('#changePasswordForm').submit(this.handleChangePassword.bind(this));
        
        // 搜索功能
        $('#searchInput').on('input', this.debounce(this.searchUsers.bind(this), 500));
        $('#roleFilter').change(this.filterUsers.bind(this));
        
        // 点击外部关闭下拉菜单
        $(document).click((e) => {
            if (!$(e.target).closest('.user-dropdown').length) {
                $('#userDropdownMenu').removeClass('show');
            }
        });
    },

    /**
     * 加载用户信息
     */
    loadUserInfo: async function() {
        try {
            const token = Utils.getToken();
            if (!token) {
                window.location.href = '/';
                return;
            }
            
            // 从token中解析用户信息（这里简化处理）
            this.userInfo = {
                username: '当前用户',
                role: 'user' // 实际应该从JWT token中解析
            };
            
            this.userRole = this.userInfo.role;
            $('#userDisplayName').text(this.userInfo.username);
            
        } catch (error) {
            console.error('加载用户信息失败:', error);
            Utils.showAlert('加载用户信息失败', 'danger');
        }
    },

    /**
     * 设置菜单显示
     */
    setupMenu: function() {
        if (this.userRole === 'admin') {
            $('#adminMenu').show();
            $('#userMenu').hide();
        } else {
            $('#adminMenu').hide();
            $('#userMenu').show();
        }
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
            case 'users':
                this.loadUsers();
                break;
            case 'tokens':
                this.loadTokens();
                break;
        }
    },

    /**
     * 加载仪表盘数据
     */
    loadDashboard: async function() {
        try {
            // 模拟数据，实际应该从API获取
            $('#totalUsers').text('12');
            $('#activeUsers').text('8');
            $('#tokenCount').text('5');
            $('#todayLogins').text('3');
        } catch (error) {
            console.error('加载仪表盘数据失败:', error);
            Utils.showAlert('加载数据失败', 'danger');
        }
    },

    /**
     * 加载用户列表
     */
    loadUsers: async function() {
        try {
            const response = await API.get('/admin/users');
            
            if (response.code === 0) {
                const users = response.data && Array.isArray(response.data.users) ? response.data.users : [];
                this.renderUsersTable(users);

                if (response.data && typeof response.data.total === 'number') {
                    $('#totalUsers').text(response.data.total);
                }
            } else {
                Utils.showAlert('加载用户列表失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('加载用户列表失败:', error);
            Utils.showAlert('加载用户列表失败', 'danger');
        }
    },

    /**
     * 渲染用户表格
     * @param {Array} users - 用户列表
     */
    renderUsersTable: function(users) {
        const tbody = $('#usersTableBody');
        
        if (users.length === 0) {
            tbody.html('<tr><td colspan="6" class="text-center text-muted">暂无用户数据</td></tr>');
            return;
        }
        
        const html = users.map(user => `
            <tr>
                <td>${user.id}</td>
                <td>${user.username}</td>
                <td>
                    <span class="badge ${user.role === 'admin' ? 'bg-danger' : 'bg-primary'}">
                        ${user.role === 'admin' ? '管理员' : '普通用户'}
                    </span>
                </td>
                <td>
                    <span class="badge bg-success">活跃</span>
                </td>
                <td>${Utils.formatDate(user.created_at)}</td>
                <td>
                    <div class="btn-group btn-group-sm">
                        <button class="btn btn-outline-primary" onclick="HomePage.editUser(${user.id})">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="btn btn-outline-danger" onclick="HomePage.deleteUser(${user.id})">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');
        
        tbody.html(html);
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
                        <button class="btn btn-outline-primary" onclick="HomePage.viewToken(${token.id})">
                            <i class="fas fa-eye"></i>
                        </button>
                        <button class="btn btn-outline-danger" onclick="HomePage.deleteToken(${token.id})">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');
        
        tbody.html(html);
    },

    /**
     * 切换用户下拉菜单
     */
    toggleUserMenu: function() {
        $('#userDropdownMenu').toggleClass('show');
    },

    /**
     * 显示用户信息
     */
    showUserInfo: function() {
        $('#userDropdownMenu').removeClass('show');
        Utils.showAlert('用户信息功能开发中...', 'info');
    },

    /**
     * 显示修改密码模态框
     */
    showChangePassword: function() {
        $('#userDropdownMenu').removeClass('show');
        $('#changePasswordModal').modal('show');
        $('#changePasswordForm')[0].reset();
    },

    /**
     * 显示添加用户模态框
     */
    showAddUserModal: function() {
        $('#addUserModal').modal('show');
        $('#addUserForm')[0].reset();
    },

    /**
     * 显示生成Token模态框
     */
    showGenerateTokenModal: function() {
        $('#generateTokenModal').modal('show');
        $('#generateTokenForm')[0].reset();
    },

    /**
     * 添加用户
     */
    addUser: async function() {
        const form = $('#addUserForm');
        const formData = new FormData(form[0]);
        
        const data = {
            username: formData.get('username'),
            password: formData.get('password'),
            role: formData.get('role')
        };

        if (!data.username || !data.password) {
            Utils.showAlert('请填写完整信息', 'warning');
            return;
        }

        const submitSelector = '#addUserForm button[type="submit"]';

        try {
            Utils.showLoading(true, submitSelector);
            
            const response = await API.post('/admin/add_user', data);
            
            if (response.code === 0) {
                Utils.showAlert('用户添加成功', 'success');
                $('#addUserModal').modal('hide');
                this.loadUsers();
            } else {
                Utils.showAlert('添加用户失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('添加用户失败:', error);
            Utils.showAlert('添加用户失败', 'danger');
        } finally {
            Utils.showLoading(false, submitSelector);
        }
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

        if (!data.name) {
            Utils.showAlert('请输入Token名称', 'warning');
            return;
        }

        try {
            Utils.showLoading(true);
            
            const response = await API.post('/user/api_tokens', data);
            
            if (response.code === 0) {
                Utils.showAlert('Token生成成功', 'success');
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
     * 修改密码
     */
    changePassword: async function() {
        const form = $('#changePasswordForm');
        const formData = new FormData(form[0]);
        
        const data = {
            old_password: formData.get('old_password'),
            new_password: formData.get('new_password')
        };

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
                $('#changePasswordModal').modal('hide');
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
     * 删除用户
     */
    deleteUser: async function(userId) {
        if (!confirm('确定要删除这个用户吗？此操作不可恢复！')) {
            return;
        }

        try {
            const response = await API.delete('/admin/delete_user', { user_id: userId });
            
            if (response.code === 0) {
                Utils.showAlert('用户删除成功', 'success');
                this.loadUsers();
            } else {
                Utils.showAlert('删除用户失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('删除用户失败:', error);
            Utils.showAlert('删除用户失败', 'danger');
        }
    },

    /**
     * 删除Token
     */
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
     * 搜索用户
     */
    searchUsers: function() {
        const keyword = $('#searchInput').val();
        console.log('搜索用户:', keyword);
    },

    /**
     * 筛选用户
     */
    filterUsers: function() {
        const role = $('#roleFilter').val();
        console.log('筛选角色:', role);
    },

    /**
     * 刷新用户列表
     */
    refreshUsers: function() {
        this.loadUsers();
    },

    /**
     * 编辑用户
     */
    editUser: function(userId) {
        Utils.showAlert('编辑功能开发中...', 'info');
    },

    /**
     * 查看Token
     */
    viewToken: function(tokenId) {
        Utils.showAlert('查看Token功能开发中...', 'info');
    },

    /**
     * 处理表单提交
     */
    handleAddUser: function(e) {
        e.preventDefault();
        this.addUser();
    },

    handleGenerateToken: function(e) {
        e.preventDefault();
        this.generateToken();
    },

    handleChangePassword: function(e) {
        e.preventDefault();
        this.changePassword();
    },

    /**
     * 登出
     */
    logout: function() {
        if (confirm('确定要退出登录吗？')) {
            Utils.clearToken();
            window.location.href = '/';
        }
    },

    /**
     * 防抖函数
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
    }
};

// 页面加载完成后初始化
$(document).ready(function() {
    HomePage.init();
});

// 导出到全局
window.HomePage = HomePage;
