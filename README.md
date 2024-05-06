# metr1c
Утилита для сбора метрик из 1С rac для Prometheus.

В данный момент собирает только информацию об активных сессиях (читать используемых пользовательских лицензиях).

## Roadmap

```
+ 0.0.1 : Первый релиз.
+ 0.1.0 : Добавить службу systemd.
0.2.0 : Добавить ansible роль. (Будет в отдельном репозитории)
...
? : Уменьшить исполняемый файл максимально насколько возможно. Сейчас (0.0.1) он весит 7 мегабайт, что много.
? : Привести структуру проекта ближе к стандарту, разделить на части.
...
1.0.0 : metric выдаёт всю информацию 1С rac и позволяет конфигурировать то, что нужно отправлять через конфиг/переменные окружения в службе. Установка в /opt через роль или скрипт. Проверено и работает на проде.
...
? : Порт на Windows (А нужно ли?) / инструкции по сборке для Windows и использованию.
? : Отдельный проект / дополнительный функционал (Какой?) через расширение информационной базы. Текущие проекты, которые я видел, не были безопасны.
? : Сборка для ALT Linux Сизиф (А это точно нужно внутри дистрибутива?)
```

## Как собрать?
```shell
make build
```

Выходной файл: `metr1c`.

## Как установить?
На сервере должен работать 1С ras (systemctl link /path/to/1c/ras-...).

### Ручная установка

Скопировать `metr1c.tar.gz` на сервер.

Далее с правами `root` выполнить:

```shell
tar -zxvf ./metr1c.tar.gz
rm ./metr1c.tar.gz
mkdir /opt/metr1c
mv ./metr1c /opt/metr1c/metr1c
mv metr1c.service /etc/systemd/system/
chown root:root /etc/systemd/system/metr1c.service
chmod 770 /etc/systemd/system/metr1c.service
nano /etc/systemd/system/metr1c.service # Указать переменные

systemctl start metr1c
```

### Установка из исходников
```
# make install
# make clean
```

# Как использовать?

Информацию выдаёт на порт `:1599` по эндпоинту `/metric`. Имя службы: `metr1c.service`.

## Безопасность

Пароль и имя админа передаются в rac как аргументы. и при установленном `hidepid=0` другие пользователи на сервере смогут их увидеть. Рекомендуется установить `hidepid=1`.

## Авторство
Лицензия - GPLv3.

Если интересует совместная разработка или проект в целом, то почта для связи - <doumor@vk.com>.
