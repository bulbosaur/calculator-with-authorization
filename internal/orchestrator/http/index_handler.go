package orchestrator

import "net/http"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Web-calculator-golang API</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      padding: 2rem;
      background: #f9f9f9;
    }
    .container {
      max-width: 600px;
      margin: auto;
      background: #fff;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 0 10px rgba(0,0,0,0.1);
    }
    input, button {
      margin: 8px 0;
      padding: 8px;
      font-size: 1rem;
    }
    .result {
      margin: 10px 0;
      padding: 10px;
      border: 1px solid #ccc;
      background: #eee;
      border-radius: 4px;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Web-calculator-golang API</h1>
    
    
    <form id="calculateForm">
      <label for="expression">Введите выражение:</label><br>
      <input type="text" id="expression" name="expression" placeholder="2 * 2" required><br>
      <button type="submit">Отправить выражение</button>
    </form>
    <div id="taskResult" class="result" style="display: none;"></div>
    
    <hr>
    
    
    <div id="finalResult" class="result" style="display: none;"></div>
  </div>
  
  <script>
    // Функция для автоматического опроса результата вычисления
    async function pollResult(taskId) {
      const resultDiv = document.getElementById('finalResult');
      // Обновляем сообщение о том, что происходит вычисление
      resultDiv.style.display = 'block';
      resultDiv.innerText = "Вычисляется...";
      try {
        const response = await fetch('http://localhost:8080/api/v1/expressions/' + encodeURIComponent(taskId));
        if (!response.ok) {
          resultDiv.innerText = 'Ошибка получения результата: ' + await response.text();
          return;
        }
        const data = await response.json();
        const expr = data.expression;
        if (expr.status === "done") {
          // Если вычисление завершено, выводим результат
          resultDiv.innerText = 'Результат вычисления: ' + expr.result;
        } else if (expr.status === "failed") {
          // Если вычисление завершилось с ошибкой, выводим текст ошибки
          resultDiv.innerText = 'Ошибка вычисления: ' + expr.ErrorMessage;
        } else {
          // Если вычисление ещё не завершено, повторяем запрос через 1 секунду
          setTimeout(() => pollResult(taskId), 1000);
        }
      } catch (error) {
        resultDiv.innerText = 'Ошибка запроса: ' + error;
      }
    }

    // Обработка отправки арифметического выражения
    document.getElementById('calculateForm').addEventListener('submit', async function(e) {
      e.preventDefault();
      // Сбрасываем предыдущий результат и сообщаем, что начинается вычисление
      const finalResultDiv = document.getElementById('finalResult');
      finalResultDiv.style.display = 'block';
      finalResultDiv.innerText = "Вычисляется...";
      
      try {
        const response = await fetch('http://localhost:8080/api/v1/calculate', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ expression: document.getElementById('expression').value })
        });
        const taskResultDiv = document.getElementById('taskResult');
        if (!response.ok) {
          const errorData = await response.json();
          taskResultDiv.style.display = 'block';
          taskResultDiv.innerText = 'Ошибка: ' + errorData.error_message;
          return;
        }
        const resJson = await response.json();
        taskResultDiv.style.display = 'block';
        taskResultDiv.innerText = 'Задача поставлена, ID задачи: ' + resJson.id;
        // Автоматический опрос результата, без необходимости нажимать кнопку
        pollResult(resJson.id);
      } catch (error) {
        const taskResultDiv = document.getElementById('taskResult');
        taskResultDiv.style.display = 'block';
        taskResultDiv.innerText = 'Ошибка запроса: ' + error;
      }
    });
  </script>
</body>
</html>`
	w.Write([]byte(html))
}
