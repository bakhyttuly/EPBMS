document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            let role = document.getElementById('role').value;
            // Map 'organizer' to 'client' for the new backend
            if (role === 'organizer') {
                role = 'client';
            }

            const payload = {
                full_name: document.getElementById('full_name').value,
                email: document.getElementById('email').value,
                password: document.getElementById('password').value,
                role: role
            };

            try {
                const response = await fetch('/api/v1/auth/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                const result = await response.json();

                if (result.success) {
                    alert('Registration successful! Please login.');
                    window.location.href = '/';
                } else {
                    alert(result.error || 'Registration failed');
                }
            } catch (error) {
                console.error('Error:', error);
                alert('An error occurred during registration');
            }
        });
    }
});
    
    // Update the dropdown to show "Client" instead of "Organizer"
    const roleSelect = document.getElementById('role');
    if (roleSelect) {
        for (let option of roleSelect.options) {
            if (option.value === 'organizer') {
                option.value = 'client';
                option.text = 'Client';
            }
        }
    }
});
