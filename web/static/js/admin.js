// 管理后台 JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // 设置启动时间
    const startTimeElement = document.getElementById('startTime');
    if (startTimeElement) {
        startTimeElement.textContent = new Date().toLocaleString('zh-CN');
    }
    
    // 绑定表单事件
    const addTokenForm = document.getElementById('addTokenForm');
    if (addTokenForm) {
        addTokenForm.addEventListener('submit', handleAddToken);
    }
    
    // 绑定按钮事件
    const refreshBtn = document.querySelector('[onclick="refreshStats()"]');
    if (refreshBtn) {
        refreshBtn.addEventListener('click', refreshStats);
    }
    
    const viewLogsBtn = document.querySelector('[onclick="viewLogs()"]');
    if (viewLogsBtn) {
        viewLogsBtn.addEventListener('click', viewLogs);
    }
});

function handleAddToken(event) {
    event.preventDefault();
    
    const userid = document.getElementById('userid').value;
    const expire = document.getElementById('expire').value;
    const successMessage = document.getElementById('successMessage');
    const errorMessage = document.getElementById('errorMessage');
    const tokenDisplay = document.getElementById('tokenDisplay');
    const tokenValue = document.getElementById('tokenValue');
    
    // 隐藏之前的消息
    hideMessage('successMessage');
    hideMessage('errorMessage');
    tokenDisplay.style.display = 'none';
    
    fetch('/admin/add_token', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({userid: parseInt(userid), expire: parseInt(expire)})
    })
    .then(response => {
        if (response.status !== 200) {
            throw new Error('生成令牌失败');
        }
        return response.json();
    })
    .then(data => {
        if (data.code === 0) {
            tokenValue.textContent = data.access_token;
            tokenDisplay.style.display = 'block';
            showMessage('successMessage', '令牌生成成功！', 'success');
        } else {
            showMessage('errorMessage', data.message || '令牌生成失败，请重试', 'error');
        }
    })
    .catch(error => {
        showMessage('errorMessage', '网络错误，请稍后重试', 'error');
        console.error('Error:', error);
    });
}

function copyToken() {
    const tokenValue = document.getElementById('tokenValue').textContent;
    copyToClipboard(tokenValue);
}

function logout() {
    if (confirm('确定要退出登录吗？')) {
        // 清除 cookie 并跳转到登录页
        document.cookie = 'admin_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
        window.location.href = '/admin/login';
    }
}

function refreshStats() {
    showMessage('successMessage', '统计信息已刷新！', 'success');
    // 这里可以添加实际的刷新逻辑
    setTimeout(() => {
        hideMessage('successMessage');
    }, 2000);
}

function viewLogs() {
    alert('日志查看功能待实现');
    // 这里可以添加查看日志的逻辑
}
