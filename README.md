# **Биллинговый микросервис Avito**
Микросервис реализует пополнение счета пользователя,получение текущего баланса пользователя,резервирование денег за услугу,ведение таблицы с доходами за услуги в отдельный .csv файл.
#### Примеры cURL запросов для обращений к микросервису
* Новый пользователь с пополнением баланса : curl -X 'POST' -d '{"money":"5000"}' 'localhost:8080/add' 
* Добавление средств на баланс существующего пользователя : curl -X 'PUT' -d '{"id":"3","money":"5000"}' 'localhost:8080/add'
* Получить баланс пользователя : curl -X 'GET' -d '{"id":"3"}’ 'localhost:8080/getBalanceOfUser' 
* Зарезервировать деньги за услугу : curl -X 'POST' -d '{"id_user":"3","id_service":"435346","id_order":"9886","cost":"2000"}' 'localhost:8080/reserve'
* Признание выручки за услугу : curl -X 'POST' -d '{"id_user":"3","id_service":"435346","id_order":"9886","cost":"2000"}' 'localhost:8080/profit' 

Записи о получении прибыли сохраняются в **Report.csv**

Дамп базы данных - файл **Billing.sql** 
