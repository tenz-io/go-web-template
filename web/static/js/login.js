// 登录页面 JavaScript

document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const messageDiv = document.getElementById('message');

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        
        if (!username || !password) {
            showMessage('请填写用户名和密码', 'error');
            return;
        }

        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                })
            });

            const data = await response.json();
            
            if (data.code === 0) {
                showMessage('登录成功，正在跳转...', 'success');
                
                // 根据用户角色跳转到不同页面
                setTimeout(() => {
                    if (data.redirect) {
                        window.location.href = data.redirect;
                    } else {
                        window.location.href = '/';
                    }
                }, 1000);
            } else {
                showMessage(data.message || '登录失败', 'error');
            }
        } catch (error) {
            console.error('Login error:', error);
            showMessage('网络错误，请稍后重试', 'error');
        }
    });

    function showMessage(text, type) {
        messageDiv.textContent = text;
        messageDiv.className = `message ${type}`;
        messageDiv.style.display = 'block';
        
        // 3秒后隐藏消息
        setTimeout(() => {
            messageDiv.style.display = 'none';
        }, 3000);
    }
});