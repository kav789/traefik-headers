# ratelimiter - traefik_ratelimiter

Плагин "ratelimiter" для сервиса traefik.

## 1.1. Хранение конфигурации

конфигурация "ratelimiter" загружается из keeper для оперативного управления "ratelimiter", или из параметров middleware **ratelimitData** при инициализации 
плагина, в случае недоступности keeper.
Данная конфигурация предназначена для конфигурирования лимитов скорости обработки запросов, в зависимости от пути и/или содержимого заголовка запроса.
конфигурация включает следующие параметры:

- **Значение лимита (`limits`)**
    - *Тип:* Массив структур
    - *Обязательность:* Да

  - **Правила (`rules`)**
      - *Тип:* Массив структур
      - *Обязательность:* Да
      - *Примечание:* Содержит правила определения запросов, для которых будет действовать ограничение RPS. Лимит будет применен, если запрос подпадет хотя бы под одно правило.
        В любом правиле должны быть указаны **urlpathpattern** и/или **headerkey** и **headerval**. 
        Например, если значение **urlpathpattern** отсутствует, то сравнение производится только по **headerkey** и **headerval** и
        наоборот, если **headerkey** или **headerval** отсутствуют, то сравнение производится только по **urlpathpattern**.
        Если присутствуют и **urlpathpattern**, и **headerkey** + **headerval**, то сравнение производится одновременно по **urlpathpattern**, и **headerkey** + **headerval** и лимит будет действовать только при полном совпадении значений **urlpathpattern**, **headerkey** + **headerval**.
        Если в лимитах присутствуют полностью идентичные правила, то срабатывать будет первое попавшееся по очереди правило.

      - **Паттерн пути (`urlpathpattern`)**
        - *Тип:* Строка
        - *Обязательность:* Нет
        - *Чуствительность к регистру значения:* Да
        - *Примечание:* Если значение присутствует, то правило будет применено, если путь http(s) запроса соответствует паттерну. 
          Описание:
          - паттерн содержит элементы пути http(s) запроса, разделенные символом **/**. Элементы пути запроса в паттерне сравниваются с соответствующими частями пути запроса на полное равенство.
            например:
              - паттерну ```"/api/v2/work/space/methods"``` будут соответствовать пути запросов  начинающихся с /api/v2/work/space/methods.
                То есть данному паттерну будут соответствовать пути: ```"/api/v2/work/space/methods"``` и  ```"/api/v2/work/space/methods/type1"```, 
                но не будет соответствовать путь ```"/api/v2/work/space/methods_type1"```
          - паттерн может содержать в элементе пути ```*```. Это означает, что эта часть пути запроса может содержать любое значение, но эта часть пути должна присутствовать обязательно.
            например:
              - паттерн ```"/api/v2/*/*/methods"``` будет соответствовать пути запроса, начинающегося с /api/v2, далее два следующие элемента должны присутствовать, но их значение может быть любым, исключая пустое (/api/v2/work//methods – превращается в /api/v2/work/methods и отрабатывает согласно алгоритма)
                и следующий элемент должен иметь значение "methods". При этом количество элементов пути может быть больше, либо равно 5.
                То есть данному паттерну будут соответствовать пути: ```"/api/v2/work/space/methods"```, ```"/api/v2/work/space/methods/type1"```, ```"/api/v2/home/space/methods/type134/somewords/othersome"``` и ```"/api/v2/home/space/methods/type1"```,
                но не будет соответствовать путь ```"/api/v2/work/space/methods_type1"```
          - паттерн может содержать в конце символ ```$```
            это означает, что путь запроса должен иметь определенную длинну в элементах пути
            например:
              - паттерн ```"/api/v2/*/*/methods$"``` будет соответствовать пути запроса, начинающегося с /api/v2 далее два следующих элемента должны присутствовать, но их значение может быть любым, исключая пустое
                и путь должен состоять ровно из 5 элементов и заканчиваться словом "methods".
                То есть данному паттерну будут соответствовать пути: ```"/api/v2/work/space/methods"```, ```"/api/v2/home/space/methods"```,
                но не будут соответствовать пути ```"/api/v2/work/space/methods_type1"``` и ```"/api/v2/home/space/methods/type1"```
          - примеры других паттернов:
            - ```"$"``` - соответствует пустому пути (срабатывает только на hostname)
            - ```"/$"```- соответствует ```/``` пути запроса (срабатывает только на hostname со слэшем в конце)
      - **Ключ из заголовка http(s) запроса (`headerkey`)**
        - *Тип:* Строка
        - *Обязательность:* Нет
        - *Чуствительность к регистру значения:* Нет
        - *Примечание:* Значение ключа из http(s) запроса. Если значение присутствует, то правило будет применено в том случае, если в запросе поле headerval не пустое и соответствует правилу.

      - **Значение соответствующего ключа из заголовка (`headerval`)**
        - *Тип:* Строка
        - *Обязательность:* Да, при наличии ключа `headerkey`
        - *Чуствительность к регистру значения:* Нет
        - *Примечание:* Значение соответствующего ключа в запросе используется только в случае, если оно содержит не пустое значение.
          Данное значение используется только в том случае если указано не пустое значение headerkey

  - **Лимит (`limit`)**
      - *Тип:* Целое число больше нуля
      - *Обязательность:* Да
      - *Примечание:*  Лимит ограничения RPS. На запросы сверх лимита будет отправлен ответ со статусом: 429 Too Many Requests.

-  примеры правил:
   - ```
     { 
       "limits": [
         {
           "rules": [
             {"urlpathpattern": "/api/v2/merchants/*/users/*/payments/methods$"}
           ],
           "limit": 10000
         }
       ]
     }
     ```
     лимит 10000 rps будет применен к запросам, совпадаюшим с правилом имеющим вид патерна пути

   - ``` 
     { 
       "limits": [
         {
           "rules": [
             {"urlpathpattern": "/api/v2/merchants/*/users/*/payments/methods$", "headerkey": "key", "headerval": "val" }
           ],
           "limit": 10000
         }
       ]
     }
     ```
     лимит 10000 rps будет применен к запросам, совпадаюшим с правилом, имеющим вид патерна,
     если заголовок запроса имеет ключ с названием "key" и его значением "val", игнорируя регистр написания ключа и значения

   - ```
     { 
       "limits": [
         {
           "rules": [
             {"urlpathpattern": "", "headerkey": "key", "headerval": "val" }
           ],
           "limit": 10000
         }
       ]
     }
     ```
     лимит 10000 rps будет применен к любым запросам, если заголовок запроса имеет ключ "key" и его значением "val", игнорируя регистр написания ключа и значения

   - ```
     { 
       "limits": [
         {
           "rules": [
             {"urlpathpattern": "/api/v2/merchants/*/users/*/payments/methods$"},
             {"urlpathpattern": "/api/v1/merchants/*/users/*/payments/methods$"}
           ],
           "limit": 10000
         }
       ]
     }
     ``` 
     лимит 10000 rps будет применен ко всем запросам, совпадаюшим с правилами имеющими вид патернов пути, то есть общее количество запросов, имеющих паттерны
     ```/api/v2/merchants/*/users/*/payments/methods$``` и ```/api/v1/merchants/*/users/*/payments/methods$```
   не будет превышать 10000 RPS/ Распределение соотношения лимитов для каждого из запросов будет рандомным 


   - ```
     { 
       "limits": [
         {
           "rules": [
             {"urlpathpattern": "/api/v2/merchants/*/users/*/payments/methods$"},
           ],
           "limit": 1000
         },

         {
           "rules": [
             {"urlpathpattern": "/api/v2/merchants/*/users/*/payments/methods$"},
             {"urlpathpattern": "/api/v1/merchants/*/users/*/payments/methods$"}
           ],
           "limit": 10000
         }
       ]
     }
     ```
     лимит 1000 rps будет применен к первому правилу с паттерном ```/api/v2/merchants/*/users/*/payments/methods$```.
     т.к. во втором лимитном правиле присуствует правило, идентичное правилу из первого лимитного правила, то оно будет проигнорировано и лимит 10000 будет применен только к правилу
     с паттерном ```/api/v1/merchants/*/users/*/payments/methods$```



Конфигурация хранится в локальном кэше ratelimiter

## 1.2. Обновление конфигурации

Конфигурация обновляется периодически, 1 раз в 30 сек из keeper

## 1.2. Параметры плагина

```
traefikMiddleware: 
    traefik-ratelimit:
    spec:
      keeperRateLimitKey: ratelimits
      keeperURL: https://keeper-ext-feature-wg-8238.k8s.dev.paywb.lan
      keeperReqTimeout: 100s
      keeperAdminPassword: Pas$w0rd
      keeperReloadInterval: 45s
      ratelimitData: |
        { 
          "limits": [
            {
              "rules": [
                {"urlpathpattern": "/api/v2/merchants/*/users/*/payments/methods$"},
              ],
              "limit": 1000
            },

            {
              "rules": [
                {"urlpathpattern": "/api/v2/merchants/*/users/*/payments/methods$"},
                {"urlpathpattern": "/api/v1/merchants/*/users/*/payments/methods$"}
              ],
              "limit": 10000
            }
          ]
        }
```
- *keeperRateLimitKey* - ключ в keeper, под которым хранится json конфигурация
- *keeperURL* - url keeper, в котором хранится json кофиграция
- *keeperReqTimeout* - таймаут ожидания ответа при запросе к keeper. По умолчанию 300s
- *keeperAdminPassword* - пароль keeper
- *keeperReloadInterval* - интервал опроса keeper для получения обновлений конфигурации. По умолчанию 30s
- *ratelimitData* - json конфигурации плагина, который будет использоваться в случае недоступности keeper при инициализации плагина

## Логика работы "ratelimiter"

Плагин сравнивает входяшие запросы со списком правил, полученых из кипер или из конфигурации middleware (при недоступности keeper в момент инициализации)
при совпадении с правилами, подсчитывает текущий RPS по правилу,
и если скорость превышает указанный лимит , то запрос не передается на дальнейшую обработку,
а создается ответ на запрос со статусом 429 Too Many Requests.