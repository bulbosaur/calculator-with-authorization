document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('loginForm').addEventListener('submit', async function(e) {
      e.preventDefault();
      
      const login = document.getElementById('login').value;
      const password = document.getElementById('password').value;
      const errorDiv = document.getElementById('errorMessage');
  
      try {
        const response = await fetch('/api/v1/login', {
          method: 'POST',
          body: JSON.stringify({ login, password }), 
        });
  
        const data = await response.json();
        
        if (!response.ok) {
          errorDiv.style.display = 'block';
          errorDiv.textContent = data.error_message || data.error || 'Ошибка авторизации';
          return;
        }
  
        localStorage.setItem('token', data.token);
        window.location.href = '/calculator?token=' + data.token;
        errorDiv.style.display = 'none';
        
      } catch (error) {
        errorDiv.style.display = 'block';
        errorDiv.textContent = 'Ошибка сети: ' + error.message;
      }
    });
  });