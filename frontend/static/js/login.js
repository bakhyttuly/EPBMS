document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;

            try {
                const response = await fetch('/api/v1/auth/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });

                const result = await response.json();

                if (result.success) {
                    // Save token and user info to localStorage
                    localStorage.setItem('token', result.data.token);
                    localStorage.setItem('user', JSON.stringify(result.data.user));
                    
                    // Redirect based on role
                    const role = result.data.user.role;
                    if (role === 'admin') {
                        window.location.href = '/dashboard-page';
                    } else if (role === 'performer') {
                        window.location.href = '/my-schedule-page';
                    } else {
                        window.location.href = '/performers-page';
                    }
                } else {
                    alert(result.error || 'Login failed');
                }
            } catch (error) {
                console.error('Error:', error);
                alert('An error occurred during login');
            }
        });
    }
});
