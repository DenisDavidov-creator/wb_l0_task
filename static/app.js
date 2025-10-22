// Находим на странице наши элементы по их 'id'
const searchButton = document.getElementById('searchButton');
const orderUidInput = document.getElementById('orderUidInput');
const orderDataElement = document.getElementById('orderData');

// "Вешаем" на кнопку "слушатель событий":
// когда произойдет событие 'click', выполнится указанная функция.
searchButton.addEventListener('click', async () => {
    
    // 1. Получаем ID заказа из поля ввода
    const orderUID = orderUidInput.value;

    // Проверяем, что ID не пустой
    if (!orderUID) {
        orderDataElement.textContent = 'Пожалуйста, введите Order UID.';
        orderDataElement.className = 'error';
        return;
    }

    try {
        // 2. Делаем GET-запрос на наш Go-сервер
        // `await` заставляет JavaScript "подождать", пока запрос не завершится
        const response = await fetch(`http://localhost:8080/api/orders/${orderUID}`);

        // 3. Проверяем статус ответа
        if (response.status === 404) {
            orderDataElement.textContent = `Заказ с UID "${orderUID}" не найден.`;
            orderDataElement.className = 'error';
            return;
        }
        if (!response.ok) {
            // Для всех других ошибок (500, 400 и т.д.)
            throw new Error(`Ошибка сервера: ${response.status}`);
        }

        // 4. Получаем тело ответа и парсим его как JSON
        const data = await response.json();

        // 5. Красиво форматируем JSON и отображаем его на странице
        // JSON.stringify(data, null, 2) превращает объект в красиво отформатированную строку
        orderDataElement.textContent = JSON.stringify(data, null, 2);
        orderDataElement.className = ''; // Убираем класс ошибки, если он был

    } catch (error) {
        // Если произошла сетевая ошибка (сервер не доступен и т.д.)
        orderDataElement.textContent = `Произошла ошибка: ${error.message}`;
        orderDataElement.className = 'error';
    }
});