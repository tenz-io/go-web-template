// 管理后台 JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // 设置启动时间
    const startTimeElement = document.getElementById('startTime');
    if (startTimeElement) {
        startTimeElement.textContent = new Date().toLocaleString('zh-CN');
    }
    
    // 初始化选项卡
    initTabs();
    
    // 初始化表单
    initChangePasswordForm();
    initAddUserForm();
    initUserManagement();
});

// 选项卡功能
function initTabs() {
    const tabBtns = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');
    
    tabBtns.forEach(btn => {
        btn.addEventListener('click', function() {
            const targetTab = this.getAttribute('data-tab');
            
            // 移除所有活动状态
            tabBtns.forEach(b => b.classList.remove('active'));
            tabContents.forEach(c => c.classList.remove('active'));
            
            // 激活当前选项卡
            this.classList.add('active');
            document.getElementById(targetTab).classList.add('active');
        });
    });
}

// 修改密码表单
function initChangePasswordForm() {
    const form = document.getElementById('changePasswordForm');
    if (!form) return;
    
    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const formData = new FormData(form);
        const data = {
            old_password: formData.get('oldPassword'),
            new_password: formData.get('newPassword')
        };
        
        // 验证密码确认
        if (data.new_password !== formData.get('confirmPassword')) {
            showMessage('changePasswordMessage', '新密码和确认密码不匹配', 'error');
            return;
        }
        
        try {
            const response = await fetch('/admin/change_password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });
            
            const result = await response.json();
            
            if (result.code === 0) {
                showMessage('changePasswordMessage', '密码修改成功', 'success');
                form.reset();
            } else {
                showMessage('changePasswordMessage', result.message || '密码修改失败', 'error');
            }
        } catch (error) {
            console.error('Error:', error);
            showMessage('changePasswordMessage', '网络错误，请重试', 'error');
        }
    });
}

// 添加用户表单
function initAddUserForm() {
    const form = document.getElementById('addUserForm');
    if (!form) return;
    
    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const formData = new FormData(form);
        const data = {
            username: formData.get('username'),
            password: formData.get('password'),
            email: formData.get('email'),
            role: formData.get('role')
        };
        
        try {
            const response = await fetch('/admin/add_user', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });
            
            const result = await response.json();
            
            if (result.code === 0) {
                showMessage('addUserMessage', '用户添加成功', 'success');
                form.reset();
                // 刷新用户列表
                loadUsers();
            } else {
                showMessage('addUserMessage', result.message || '用户添加失败', 'error');
            }
        } catch (error) {
            console.error('Error:', error);
            showMessage('addUserMessage', '网络错误，请重试', 'error');
        }
    });
}

// 用户管理功能
function initUserManagement() {
    const refreshBtn = document.getElementById('refreshUsersBtn');
    if (refreshBtn) {
        refreshBtn.addEventListener('click', loadUsers);
    }
    
    // 初始加载用户列表
    loadUsers();
}

// 加载用户列表
async function loadUsers() {
    try {
        const response = await fetch('/admin/users');
        const result = await response.json();
        
        if (result.code === 0) {
            displayUsers(result.data || []);
        } else {
            console.error('Failed to load users:', result.message);
        }
    } catch (error) {
        console.error('Error loading users:', error);
    }
}

// 显示用户列表
function displayUsers(users) {
    const tbody = document.getElementById('userTableBody');
    if (!tbody) return;
    
    tbody.innerHTML = '';
    
        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${user.id}</td>
                <td>${user.username}</td>
                <td><span class="role-badge ${user.role}">${user.role === 'admin' ? '管理员' : '普通用户'}</span></td>
                <td>${user.email}</td>
                <td>${new Date(user.created_at).toLocaleString()}</td>
            `;
            tbody.appendChild(row);
        });
}

// 显示消息
function showMessage(elementId, message, type) {
    const element = document.getElementById(elementId);
    if (!element) return;
    
    element.textContent = message;
    element.className = `message ${type}`;
    element.style.display = 'block';
    
    setTimeout(() => {
        element.style.display = 'none';
    }, 3000);
}

// 退出登录
function logout() {
    if (confirm('确定要退出登录吗？')) {
        window.location.href = '/auth/logout';
    }
}

// 刷新统计
function refreshStats() {
    alert('统计信息已刷新！');
}

// 查看日志
function viewLogs() {
    alert('日志功能开发中...');
}