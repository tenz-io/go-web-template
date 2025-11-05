/**
 * 统一的控制台脚本
 */

const HomePage = {
    currentTab: 'dashboard',
    userRole: 'user',
    userInfo: {},

    init() {
        this.loadUserInfo();
        this.bindEvents();
        this.setupMenu();
        this.showTab(this.currentTab);
    },

    loadUserInfo() {
        const dataset = document.body.dataset || {};
        const role = (dataset.role || '').toLowerCase();
        const isAdmin = dataset.isAdmin === 'true' || role === 'admin';
        const displayName = dataset.displayName || dataset.username || '用户';

        this.userInfo = {
            username: dataset.username || '',
            displayName,
            appName: dataset.appName || '',
            role: isAdmin ? 'admin' : (role || 'user'),
            isAdmin,
        };
        this.userRole = this.userInfo.role;

        $('#userDisplayName').text(displayName);
        $('#userRoleBadge').text(isAdmin ? '管理员' : '普通用户');

        if (this.userInfo.appName) {
            document.title = `${this.userInfo.appName} - 控制台`;
        }
    },

    bindEvents() {
        $('#sidebarMenu').on('click', '.nav-link[data-tab]', (e) => {
            e.preventDefault();
            const tab = $(e.currentTarget).data('tab');
            this.showTab(tab);
        });

        const $addUserForm = $('#addUserForm');
        if ($addUserForm.length) {
            $addUserForm.on('submit', this.handleAddUser.bind(this));
        }

        const $generateTokenForm = $('#generateTokenForm');
        if ($generateTokenForm.length) {
            $generateTokenForm.on('submit', this.handleGenerateToken.bind(this));
        }

        const $changePasswordForm = $('#changePasswordForm');
        if ($changePasswordForm.length) {
            $changePasswordForm.on('submit', this.handleChangePassword.bind(this));
        }

        $(document).on('click', (e) => {
            if (!$(e.target).closest('.user-dropdown').length) {
                $('#userDropdownMenu').removeClass('show');
            }
        });
    },

    setupMenu() {
        $('#sidebarMenu [data-visible-for]').each((_, element) => {
            const $element = $(element);
            const visibleFor = ($element.data('visible-for') || 'all').toString();
            const canShow = this.canDisplay(visibleFor);
            $element.toggle(canShow);
            if (!canShow) {
                $element.removeClass('active');
            }
        });

        $('[data-role="admin"]').toggle(this.userRole === 'admin');

        if (this.userRole !== 'admin' && this.currentTab === 'users') {
            this.currentTab = 'dashboard';
        }

        const $active = $('#sidebarMenu .nav-link.active:visible');
        if (!$active.length) {
            const fallback = $('#sidebarMenu .nav-link[data-tab]:visible').first().data('tab') || 'dashboard';
            this.currentTab = fallback;
        } else {
            this.currentTab = $active.data('tab') || this.currentTab;
        }
    },

    canDisplay(visibleFor) {
        if (!visibleFor || visibleFor === 'all') {
            return true;
        }

        const role = this.userRole;
        const tokens = visibleFor.split(',').map((item) => item.trim().toLowerCase());

        if (tokens.includes('admin')) {
            return role === 'admin';
        }

        if (tokens.includes('user')) {
            return role !== 'admin';
        }

        return tokens.includes(role);
    },

    showTab(tabName) {
        const $targetLink = $(`.nav-link[data-tab="${tabName}"]`);
        if (!$targetLink.length || !$targetLink.is(':visible')) {
            const fallback = $('#sidebarMenu .nav-link[data-tab]:visible').first().data('tab') || 'dashboard';
            if (fallback && fallback !== tabName) {
                this.showTab(fallback);
            }
            return;
        }

        $('.tab-content').hide();
        $(`#${tabName}`).show();
        $('#sidebarMenu .nav-link').removeClass('active');
        $targetLink.addClass('active');
        this.currentTab = tabName;
        this.updatePageTitle(tabName);

        switch (tabName) {
            case 'dashboard':
                this.loadDashboard();
                break;
            case 'users':
                this.loadUsers();
                break;
            case 'tokens':
                this.loadTokens();
                break;
            default:
                break;
        }
    },

    updatePageTitle(tabName) {
        const $nav = $(`.nav-link[data-tab="${tabName}"]`);
        if ($nav.length) {
            const text = $.trim($nav.text());
            $('#pageTitle').text(text || '仪表盘');
        }
    },

    async loadDashboard() {
        try {
            $('#totalUsers').text('12');
            $('#activeUsers').text('8');
            $('#tokenCount').text('5');
            $('#todayLogins').text('3');
        } catch (error) {
            console.error('加载仪表盘数据失败:', error);
            Utils.showAlert('加载仪表盘数据失败', 'danger');
        }
    },

    async loadUsers() {
        if (this.userRole !== 'admin') {
            return;
        }

        try {
            const response = await API.get('/admin/users');
            if (response.code === 0) {
                const users = response.data && Array.isArray(response.data.users) ? response.data.users : [];
                this.renderUsersTable(users);
                if (response.data && typeof response.data.total === 'number') {
                    $('#totalUsers').text(response.data.total);
                }
            } else {
                Utils.showAlert(`加载用户列表失败: ${response.message}`, 'danger');
            }
        } catch (error) {
            console.error('加载用户列表失败:', error);
            Utils.showAlert('加载用户列表失败', 'danger');
        }
    },

    renderUsersTable(users) {
        const $tbody = $('#usersTableBody');
        if (!users.length) {
            $tbody.html('<tr><td colspan="6" class="text-center text-muted">暂无用户数据</td></tr>');
            return;
        }

        const html = users.map((user) => {
            const roleLabel = user.role === 'admin' ? '管理员' : '普通用户';
            const roleClass = user.role === 'admin' ? 'role-badge admin' : 'role-badge user';

            return `
                <tr>
                    <td>${user.id}</td>
                    <td>${user.username}</td>
                    <td><span class="${roleClass}">${roleLabel}</span></td>
                    <td><span class="badge bg-success">活跃</span></td>
                    <td>${Utils.formatDate(user.created_at)}</td>
                    <td>
                        <button class="btn btn-sm btn-outline-danger" onclick="HomePage.deleteUser(${user.id})">
                            <i class="fas fa-trash"></i>
                        </button>
                    </td>
                </tr>
            `;
        }).join('');

        $tbody.html(html);
    },

    loadTokens() {
        this.renderTokensTable([]);
    },

    renderTokensTable(tokens) {
        const $tbody = $('#tokensTableBody');
        if (!tokens.length) {
            $tbody.html('<tr><td colspan="5" class="text-center text-muted">生成的 Token 仅会在创建时显示，请在上方保存。</td></tr>');
            return;
        }

        const rows = tokens.map((token) => `
            <tr>
                <td>${token.value || '-'}</td>
                <td>${token.role || '-'}</td>
                <td>${token.expire_at ? Utils.formatDate(token.expire_at) : '-'}</td>
                <td>${token.created_at ? Utils.formatDate(token.created_at) : '-'}</td>
                <td>
                    <button class="btn btn-sm btn-outline-danger" onclick="HomePage.deleteToken(${token.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            </tr>
        `).join('');

        $tbody.html(rows);
    },

    toggleUserMenu() {
        $('#userDropdownMenu').toggleClass('show');
    },

    showUserInfo() {
        $('#userDropdownMenu').removeClass('show');
        Utils.showAlert('用户信息功能开发中...', 'info');
    },

    showChangePassword() {
        $('#userDropdownMenu').removeClass('show');
        $('#changePasswordForm')[0].reset();
        $('#changePasswordModal').modal('show');
    },

    showAddUserModal() {
        if (this.userRole !== 'admin') {
            Utils.showAlert('只有管理员可以添加用户', 'warning');
            return;
        }
        $('#addUserForm')[0].reset();
        $('#addUserModal').modal('show');
    },

    async addUser() {
        if (this.userRole !== 'admin') {
            Utils.showAlert('只有管理员可以添加用户', 'warning');
            return;
        }

        const form = $('#addUserForm')[0];
        const formData = new FormData(form);
        const data = {
            username: formData.get('username'),
            password: formData.get('password'),
            role: formData.get('role'),
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
                Utils.showAlert(`添加用户失败: ${response.message}`, 'danger');
            }
        } catch (error) {
            console.error('添加用户失败:', error);
            Utils.showAlert('添加用户失败', 'danger');
        } finally {
            Utils.showLoading(false, submitSelector);
        }
    },

    handleAddUser(e) {
        e.preventDefault();
        this.addUser();
    },

    async generateToken() {
        const form = $('#generateTokenForm')[0];
        const formData = new FormData(form);
        const expireHours = parseInt(formData.get('expire_hours'), 10);

        if (!Number.isFinite(expireHours) || expireHours <= 0) {
            Utils.showAlert('请选择有效的过期时间', 'warning');
            return;
        }

        const submitSelector = '#generateTokenForm button[type="submit"]';

        try {
            Utils.showLoading(true, submitSelector);
            const response = await API.post('/user/generate_token', { expire_hours: expireHours });

            if (response.code === 0 && response.token) {
                $('#generatedToken').val(response.token);
                $('#tokenResultPlaceholder').addClass('d-none');
                $('#tokenResultRow').removeClass('d-none');
                Utils.showAlert('Token 生成成功，请及时保存', 'success');
                form.reset();
            } else {
                Utils.showAlert(`生成 Token 失败: ${response.message || '未知错误'}`, 'danger');
            }
        } catch (error) {
            console.error('生成Token失败:', error);
            Utils.showAlert('生成 Token 失败', 'danger');
        } finally {
            Utils.showLoading(false, submitSelector);
        }
    },

    handleGenerateToken(e) {
        e.preventDefault();
        this.generateToken();
    },

    async changePassword() {
        const form = $('#changePasswordForm')[0];
        const formData = new FormData(form);
        const data = {
            old_password: formData.get('old_password'),
            new_password: formData.get('new_password'),
        };
        const confirmPassword = $('#confirmPassword').val();

        if (data.new_password !== confirmPassword) {
            Utils.showAlert('两次输入的新密码不一致', 'warning');
            return;
        }

        const submitSelector = '#changePasswordForm button[type="submit"]';

        try {
            Utils.showLoading(true, submitSelector);
            const response = await API.post('/auth/change_password', data);
            if (response.code === 0) {
                Utils.showAlert('密码修改成功', 'success');
                $('#changePasswordModal').modal('hide');
                form.reset();
            } else {
                Utils.showAlert(`密码修改失败: ${response.message}`, 'danger');
            }
        } catch (error) {
            console.error('密码修改失败:', error);
            Utils.showAlert('密码修改失败', 'danger');
        } finally {
            Utils.showLoading(false, submitSelector);
        }
    },

    handleChangePassword(e) {
        e.preventDefault();
        this.changePassword();
    },

    async deleteUser(userId) {
        if (this.userRole !== 'admin') {
            Utils.showAlert('只有管理员可以删除用户', 'warning');
            return;
        }

        if (!confirm('确定要删除这个用户吗？此操作不可恢复！')) {
            return;
        }

        try {
            const response = await API.delete('/admin/delete_user', { user_id: userId });
            if (response.code === 0) {
                Utils.showAlert('用户删除成功', 'success');
                this.loadUsers();
            } else {
                Utils.showAlert(`删除用户失败: ${response.message}`, 'danger');
            }
        } catch (error) {
            console.error('删除用户失败:', error);
            Utils.showAlert('删除用户失败', 'danger');
        }
    },

    deleteToken() {
        Utils.showAlert('当前版本仅支持临时 Token，无需删除。', 'info');
    },

    copyToken() {
        const token = $('#generatedToken').val();
        if (!token) {
            Utils.showAlert('暂无可复制的 Token', 'warning');
            return;
        }
        Utils.copyToClipboard(token);
    },

    logout() {
        $('#userDropdownMenu').removeClass('show');
        if (!confirm('确定要退出登录吗？')) {
            return;
        }
        Utils.clearToken();
        API.post('/logout', {}).finally(() => {
            window.location.href = '/login';
        });
    },
};

$(document).ready(() => {
    HomePage.init();
});

window.HomePage = HomePage;
