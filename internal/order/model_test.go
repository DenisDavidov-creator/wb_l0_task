package order // Тестовый файл находится в том же пакете, что и тестируемый код

import (
	"testing" // Главный пакет для тестирования в Go
	"time"
)

// Название тестовой функции ОБЯЗАТЕЛЬНО должно начинаться с `Test`.
// Она принимает один аргумент: `t *testing.T`.
// `t` - это наш "помощник" для управления тестом.
func TestOrder_Validate(t *testing.T) {

	// --- 1. Подготовка тестовых данных ---

	// Создаем полностью валидный объект Order.
	// Это наш "эталон" для "счастливого пути".
	validOrder := Order{
		OrderUID:        "b563feb7b2b84b6test",
		TrackNumber:     "WBILMTESTTRACK",
		Entry:           "WBIL",
		Delivery:        Delivery{Name: "Test", Phone: "+123", Zip: "12345", City: "Test City", Address: "Test Address", Email: "test@test.com"},
		Payment:         Payment{Transaction: "b563feb7b2b84b6test", Currency: "USD", Provider: "wbpay", Amount: 1, PaymentDt: 1, Bank: "alpha", DeliveryCost: 1, GoodsTotal: 1},
		Items:           []Item{{ChrtID: 1, TrackNumber: "WBILMTESTTRACK", Price: 1, Rid: "rid", Name: "name", Sale: 0, Size: "0", TotalPrice: 1, NmID: 1, Status: 200, Brand: "brand"}},
		Locale:          "en",
		CustomerID:      "test",
		DeliveryService: "meest",
		DateCreated:     time.Now(),
	}

	// --- 2. Определение тестовых сценариев (Test Cases) ---

	// Создаем срез структур, где каждая структура - это один тест.
	// Этот подход позволяет легко добавлять новые тесты.
	testCases := []struct {
		name    string // Название теста для понятных логов
		order   Order  // Данные для теста
		wantErr bool   // Ожидаем ли мы ошибку? (true/false)
	}{
		{
			name:    "Корректный заказ",
			order:   validOrder, // Используем наш эталон
			wantErr: false,      // Ошибки быть не должно
		},
		{
			name: "Заказ без OrderUID",
			order: func() Order {
				o := validOrder // Берем за основу валидный заказ
				o.OrderUID = "" // "Ломаем" его - делаем одно поле невалидным
				return o
			}(),
			wantErr: true, // Ожидаем ошибку
		},
		{
			name: "Заказ без товаров (Items)",
			order: func() Order {
				o := validOrder
				o.Items = []Item{} // Делаем срез товаров пустым
				return o
			}(),
			wantErr: true, // Ожидаем ошибку (т.к. у нас `validate:"min=1"`)
		},
		{
			name: "Некорректный Email в доставке",
			order: func() Order {
				o := validOrder
				o.Delivery.Email = "not-an-email" // Ломаем вложенную структуру
				return o
			}(),
			wantErr: true, // Ожидаем ошибку
		},
	}

	// --- 3. Запуск тестов ---

	// Пробегаемся в цикле по всем нашим сценариям
	for _, tc := range testCases {
		// t.Run() запускает каждый сценарий как отдельный "подтест".
		// Это очень удобно для группировки и вывода результатов.
		t.Run(tc.name, func(t *testing.T) {
			// Выполняем тестируемую функцию
			err := tc.order.Validate()

			// Проверяем результат
			if (err != nil) != tc.wantErr {
				// (err != nil) вернет true, если ошибка есть, и false, если ее нет.
				// Мы сравниваем этот результат с тем, что ожидали (tc.wantErr).
				// Если они не совпадают - тест провален.
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
