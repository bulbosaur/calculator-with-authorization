async function pollResult(taskId) {
    if (!localStorage.getItem('token')) {
      alert('Требуется авторизация!');
      window.location.href = '/login';
      return;
    }

    const resultDiv = document.getElementById('finalResult');
    resultDiv.style.display = 'block';
    resultDiv.innerText = "Вычисляется...";
    try {
      const response = await fetch('/api/v1/expressions/' + encodeURIComponent(taskId), {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      if (!response.ok) {
        resultDiv.innerText = 'Ошибка получения результата: ' + await response.text() + 'id: ' + taskId;
        return;
      }
      const data = await response.json();
      const expr = data.expression;
      if (expr.status === "done") {
        resultDiv.innerText = 'Результат вычисления: ' + expr.result;
      } else if (expr.status === "failed") {
        resultDiv.innerText = 'Ошибка вычисления: ' + expr.ErrorMessage;
      } else {
        setTimeout(() => pollResult(taskId), 1000);
      }
    } catch (error) {
      resultDiv.innerText = 'Ошибка запроса: ' + error;
    }
  }

  document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('calculateForm').addEventListener('submit', async function(e) {
      e.preventDefault();
      if (!localStorage.getItem('token')) {
        alert('Требуется авторизация!');
        window.location.href = '/login';
        return;
      }
      const finalResultDiv = document.getElementById('finalResult');
      finalResultDiv.style.display = 'block';
      finalResultDiv.innerText = "Вычисляется...";
      
      try {
        const response = await fetch('/api/v1/calculate', {
          method: 'POST',
          headers: { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: JSON.stringify({ expression: document.getElementById('expression').value })
        });
        
        const taskResultDiv = document.getElementById('taskResult');
        if (!response.ok) {
          const errorData = await response.json();
          taskResultDiv.style.display = 'block';
          taskResultDiv.innerText = 'Ошибка: ' + (errorData.error_message || errorData.error);
          return;
        }
        
        const resJson = await response.json();
        taskResultDiv.style.display = 'block';
        taskResultDiv.innerText = 'Задача поставлена, ID задачи: ' + resJson.id;
        pollResult(resJson.id);
      } catch (error) {
        const taskResultDiv = document.getElementById('taskResult');
        taskResultDiv.style.display = 'block';
        taskResultDiv.innerText = 'Ошибка запроса: ' + error;
      }
    });
  });