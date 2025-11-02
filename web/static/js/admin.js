/**
 * 管理员页面JavaScript
 */

// 管理员页面功能
const AdminPage = {
    currentTab: 'dashboard',
    
    /**
     * 初始化页面
     */
    init: function() {
        this.bindEvents();
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
        $('#adminChangePasswordForm').submit(this.handleChangePassword.bind(this));
        
        // 搜索功能
        $('#searchInput').on('input', this.debounce(this.searchUsers.bind(this), 500));
        $('#roleFilter').change(this.filterUsers.bind(this));
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
            case 'settings':
                this.loadSettings();
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
            $('#adminUsers').text('2');
            $('#todayLogins').text('5');
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
                        <button class="btn btn-outline-primary" onclick="AdminPage.editUser(${user.id})">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="btn btn-outline-danger" onclick="AdminPage.deleteUser(${user.id})">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');
        
        tbody.html(html);
    },

    /**
     * 显示添加用户模态框
     */
    showAddUserModal: function() {
        $('#addUserModal').modal('show');
        $('#addUserForm')[0].reset();
    },

    /**
     * 显示修改密码模态框
     */
    showChangePasswordModal: function() {
        $('#adminChangePasswordModal').modal('show');
        $('#adminChangePasswordForm')[0].reset();
    },

    /**
     * 处理添加用户表单提交
     * @param {Event} e - 事件对象
     */
    handleAddUser: function(e) {
        e.preventDefault();
        this.addUser();
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

        // 验证输入
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
     * 删除用户
     * @param {number} userId - 用户ID
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
     * 编辑用户
     * @param {number} userId - 用户ID
     */
    editUser: function(userId) {
        // 实现编辑用户功能
        Utils.showAlert('编辑功能开发中...', 'info');
    },

    /**
     * 搜索用户
     */
    searchUsers: function() {
        const keyword = $('#searchInput').val();
        // 实现搜索功能
        console.log('搜索用户:', keyword);
    },

    /**
     * 筛选用户
     */
    filterUsers: function() {
        const role = $('#roleFilter').val();
        // 实现筛选功能
        console.log('筛选角色:', role);
    },

    /**
     * 加载设置
     */
    loadSettings: function() {
        // 加载系统设置
        console.log('加载系统设置');
    },


    /**
     * 处理管理员修改密码
     */
    handleChangePassword: async function(e) {
        e.preventDefault();

        const form = $('#adminChangePasswordForm');
        const formData = new FormData(form[0]);

        const data = {
            old_password: formData.get('old_password'),
            new_password: formData.get('new_password')
        };

        const confirmPassword = $('#adminConfirmPassword').val();
        if (data.new_password !== confirmPassword) {
            Utils.showAlert('两次输入的新密码不一致', 'warning');
            return;
        }

        const submitSelector = '#adminChangePasswordForm button[type="submit"]';

        try {
            Utils.showLoading(true, submitSelector);

            const response = await API.post('/admin/change_password', data);

            if (response.code === 0) {
                Utils.showAlert('密码修改成功', 'success');
                form[0].reset();
                $('#adminChangePasswordModal').modal('hide');
            } else {
                Utils.showAlert('密码修改失败: ' + response.message, 'danger');
            }
        } catch (error) {
            console.error('管理员修改密码失败:', error);
            Utils.showAlert('密码修改失败', 'danger');
        } finally {
            Utils.showLoading(false, submitSelector);
        }
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
    AdminPage.init();
});

// 导出到全局
window.AdminPage = AdminPage;
