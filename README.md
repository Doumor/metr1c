# metr1c
Утилита для сбора метрик из 1С rac для Prometheus.

В данный момент собирает информацию
количествах сеансов (активных и спящих), используемых лицензиях (soft и HASP), соединениях, процессах, памяти и сессиях на каждую информационную базу.

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

systemctl enable metr1c
systemctl start metr1c
```

### Установка из исходников
```shell
make install
make clean
```

## Как использовать?

Информацию выдаёт на порт `:1599` по эндпоинту `/metric`. Имя службы: `metr1c.service`.

Выдает метрики:
1) sessionCount - количество сеансов на сервере
2) activeSessionCount - количество активных сеансов
3) hibernatedSessionCount - количество спящих сеансов
4) softLicensesCount - количество soft лицензий
5) haspLicensesCount - количество HASP лицензий
6) connectionCount - количество соединений
7) processCount - количество процессов
8) processMemTotal - сколько памяти потребляется всеми процессами, в кб
9) sessionsPerInfobase - отображает то, как каждая из существующих баз используется в данный момент, сколько активных сеансов к ней обращено

## Безопасность

Пароль и имя админа передаются в rac как аргументы. и при установленном `hidepid=0` другие пользователи на сервере смогут их увидеть. Рекомендуется установить `hidepid=1`.

## Авторство
Лицензия - GPLv3.

Если интересует совместная разработка или проект в целом, то почта для связи - <doumor@vk.com>.
