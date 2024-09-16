#  Начало работы:
Для размещения приложения на сервере можно воспользоваться bashScript он находится в папке deployments под названием 
installApp.sh. Данный bashScript имеет 3 входных аргумента:
1. Имя пользователя на сервере <name>
2. Ip адрес сервера <ip>
3. Путь где находится бинарный файл **antifraud** <patch> 

Пример:

**`bash installApp.sh kaban 192.168.122.171 ~/"GolandProjects/Antifraud-payment-system/build/antifraud"`**

Чтобы все работало корректно убедитесь что сервер работает и находится в сети, убедитесь что бинарный файл **antifraud**
находится в папке которую вы указали в 3 аргументе файла bashScript. Убедитесь что вы ввели все аргументы, корректно.
После bashScript проверит, работает ли программа **antifraud** на сервере в текущий момент если работает то тогда в свою
очередь bashScript остановит программу. Так же он проверит, имеется ли на сервере бинарный файл **antifraud** если 
имеется то он будет принудительно удален, а вместо него будет скопирован новый. Бинарный файл **antifraud** по default 
устанавливается в папку ~/home/<UserName>, а так же запускается в работу.


Также Вы можете клонировать репозиторий с сайта GitHub с помощью команды 
**`git clone git@github.com:1Kabman1/antifraud-payment-system.git`**
После чего вам следует произвести тестирование приложения с помощью команды **`go test ./...`** если тестирование
пройдет успешно, запустите приложение с помощью команды **`go run .`** обязательно убедитесь что Вы находитесь в 
директории приложения **`/cmd/app`**. Вы так же можете создать бинарный файл командой **`go build`**

## POST// создание  Rule
Чтобы создать правило агрегации, нужно отправить POST запрос на url "http://your_domen/aggregation_rule/create"

Host:    
   **`curl -XPOST -v  'http://your_domen/aggregation_rule/create' -H 'Content-Type: appliction/json' -d '{}'`**
 
в формате JSON:

Пример №1:
```json
{
    "Name": "<имя для пользователя, для удобства>",
    "AggregateBy": "<агрегируемые свойства>",
    "AggregateValue": "<сумма или счетчик>",
    "ExpirationTime": "<время жизни агрегируемого в минутах>",
    "TimePeriod": "<точность подсчета агрегируемого в секундах>"
}
```
Пример №2:

```json
{
    "Name": "Amount per client",
    "AggregateBy": ["clientId", "bank_card"],
    "AggregateValue": "amount",
    "ExpirationTime": 2,
    "TimePeriod": 22
} 
```

В следствии программа создаст правила агрегации с уникальным id. Имя правила Вы задаете сами, имена могут быть 
одинаковые и агрегируемые свойства, но при этом правила все равно будут уникальны по отношению друг к другу, 
уникальность правилам придает уникальный id для каждого правила. 

Пример:
```json
{
   "1": {
      "AggregationRuleId": 1,
      "Name": "Amount per client",
      "AggregateBy": [
         "clientId",
         "bank_card"
      ],
      "AggregateValue": "count",
      "ExpirationTime":2,
      "TimePeriod":22
      
   },
   "2": {
      "AggregationRuleId": 2,
      "Name": "Amount per client",
      "AggregateBy": [
         "clientId",
         "bank_card"
      ],
      "AggregateValue": "amount",
      "ExpirationTime":2,
      "TimePeriod":22
   }
}
```

## GET// запрос Rule 
Чтобы получить данные по созданным правилам, нужно отправить GET запрос на 
url "http://your_domen/aggregation_rules/get".

Host:
    **`curl -XGET -v 'http://your_domen/aggregation_rules/get'`**
 

В следствии программа вернет ответ в JSON формате. (! Программа возвращает полный список правил)

Пример:

```json
{
  "1": {
    "AggregationRuleId": 1,
    "Name": "Amount per client",
    "AggregateBy": [
      "clientId",
      "bank_card"
    ],
    "AggregateValue": "amount",
     "ExpirationTime":2,
     "TimePeriod":22
  },
  "2": {
    "AggregationRuleId": 2,
    "Name": "Amount per client",
    "AggregateBy": [
      "clientId",
      "bank_card"
    ],
    "AggregateValue": "amount",
     "ExpirationTime":2,
     "TimePeriod":22
  }
}
```

## POST // Отслеживание агрегируемого  

Чтобы начать отслеживать агрегируемое на основании созданных правил, отправьте POST запрос на 
url "http://your_domen/register"

Host:
     **`curl -POST -v  'http://your_domen/register' -H 'Content-Type: appliction/json' -d '{}'`**


Пример:

```json
{
"payment_id": 123123,
"client_id": "abasdi-1923",
"payment_method_type": "bank_card",
"payment_method_id": "aopiasdxscnlojcxzoqwe",
"amount": 1000,
"currency": "RUB"
}
```

Программа проверит данное агрегируемое по всем правилам, а так же зафиксирует "count" OR "amount" в специально созданном
уникальном счетчике.
 
 

