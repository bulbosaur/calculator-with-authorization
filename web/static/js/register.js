document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('registerForm').addEventListener('submit', async function(e) {
      e.preventDefault();
      
      const login = document.getElementById('login').value;
      const password = document.getElementById('password').value;
      const errorDiv = document.getElementById('errorMessage');
  
      try {
        const response = await fetch('/api/v1/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ login, password })
        });
  
        const data = await response.json();
        
        if (!response.ok) {
          errorDiv.style.display = 'block';
          errorDiv.textContent = data.error_message || data.error || 'Ошибка регистрации';
          return;
        }
  
        errorDiv.style.display = 'none';
        alert('Регистрация успешна! Теперь войдите в систему.');
        window.location.href = '/login';
        
      } catch (error) {
        errorDiv.style.display = 'block';
        errorDiv.textContent = 'Ошибка сети: ' + error.message;
      }
    });
  });