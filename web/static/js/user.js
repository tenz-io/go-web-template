/**
 * 用户页面脚本
 */

const UserPage = {
    init() {
        this.bindEvents();
        this.showTab('dashboard');
    },

    bindEvents() {
        // 标签切换
        $('.nav-link[data-tab]').on('click', (e) => {
            e.preventDefault();
            const tab = $(e.currentTarget).data('tab');
            this.showTab(tab);
        });

        // 修改密码
        $('#changePasswordForm').on('submit', this.handleChangePassword.bind(this));

        // 生成 Token
        $('#generateTokenForm').on('submit', this.handleGenerateToken.bind(this));
    },

    showTab(tab) {
        $('.nav-link[data-tab]').removeClass('active');
        $(`.nav-link[data-tab="${tab}"]`).addClass('active');

        $('.tab-pane').addClass('d-none');
        $(`#tab-${tab}`).removeClass('d-none');
    },

    async handleChangePassword(e) {
        e.preventDefault();

        const form = $('#changePasswordForm');
        const formData = new FormData(form[0]);
        const submitSelector = '#changePasswordForm button[type="submit"]';

        const data = {
            old_password: formData.get('old_password'),
            new_password: formData.get('new_password')
        };

        if (data.new_password !== $('#confirmPassword').val()) {
            Utils.showAlert('两次输入的新密码不一致', 'warning');
            return;
        }

        try {
            Utils.showLoading(true, submitSelector);
            const response = await API.post('/auth/change_password', data);

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
            Utils.showLoading(false, submitSelector);
        }
    },

    async handleGenerateToken(e) {
        e.preventDefault();

        const form = $('#generateTokenForm');
        const formData = new FormData(form[0]);
        const submitSelector = '#generateTokenForm button[type="submit"]';

        let expire = parseInt(formData.get('expire'), 10);
        if (!Number.isFinite(expire) || expire <= 0) {
            Utils.showAlert('请选择有效的过期时间', 'warning');
            return;
        }

        try {
            Utils.showLoading(true, submitSelector);
            const response = await API.post('/user/generate_token', { expire });

            if (response.code === 0 && response.token) {
                $('#generatedToken').val(response.token);
                $('#tokenResultRow').removeClass('d-none');
                Utils.showAlert('令牌生成成功，请及时保存', 'success');
                form[0].reset();
            } else {
                Utils.showAlert('生成令牌失败: ' + (response.message || '未知错误'), 'danger');
            }
        } catch (error) {
            console.error('生成令牌失败:', error);
            Utils.showAlert('生成令牌失败', 'danger');
        } finally {
            Utils.showLoading(false, submitSelector);
        }
    },

    copyToken() {
        const token = $('#generatedToken').val();
        if (!token) {
            Utils.showAlert('没有可复制的令牌', 'warning');
            return;
        }
        Utils.copyToClipboard(token);
    },

    logout() {
        if (!confirm('确定要退出登录吗？')) {
            return;
        }

        Utils.clearToken();
        API.post('/logout', {}).finally(() => {
            window.location.href = '/login';
        });
    }
};

$(document).ready(() => {
    UserPage.init();
});

window.UserPage = UserPage;
