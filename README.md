# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

```
go tool pprof -top -diff_base=./profiles/base.pprof  ./profiles/result.pprof 
File: __debug_bin1529257225
Type: inuse_space
Time: 2026-01-28 12:19:12 MSK
Showing nodes accounting for -1158.25kB, 33.41% of 3466.76kB total
      flat  flat%   sum%        cum   cum%
   -1026kB 29.60% 29.60%    -1026kB 29.60%  runtime.allocm
 -902.59kB 26.04% 55.63%  -902.59kB 26.04%  compress/flate.NewWriter
  768.26kB 22.16% 33.47%   768.26kB 22.16%  go.uber.org/zap/zapcore.newCounters
  515.19kB 14.86% 18.61%   515.19kB 14.86%  unicode.map.init.1
 -513.12kB 14.80% 33.41%  -513.12kB 14.80%  compress/flate.(*huffmanEncoder).generate
  512.05kB 14.77% 18.64%   512.05kB 14.77%  context.(*cancelCtx).Done
 -512.05kB 14.77% 33.41%  -512.05kB 14.77%  runtime.(*scavengerState).init
         0     0% 33.41%  -513.12kB 14.80%  compress/flate.(*Writer).Close
         0     0% 33.41%  -513.12kB 14.80%  compress/flate.(*compressor).close
         0     0% 33.41%  -513.12kB 14.80%  compress/flate.(*compressor).encSpeed
         0     0% 33.41%  -513.12kB 14.80%  compress/flate.(*huffmanBitWriter).indexTokens
         0     0% 33.41%  -513.12kB 14.80%  compress/flate.(*huffmanBitWriter).writeBlockDynamic
         0     0% 33.41%  -513.12kB 14.80%  compress/gzip.(*Writer).Close
         0     0% 33.41%  -902.59kB 26.04%  compress/gzip.(*Writer).Write
         0     0% 33.41%   512.05kB 14.77%  database/sql.(*DB).connectionOpener
         0     0% 33.41%   768.26kB 22.16%  github.com/f044fs3t5w3f/metrics/internal/logger.Initialize
         0     0% 33.41%   768.26kB 22.16%  go.uber.org/zap.(*Logger).WithOptions
         0     0% 33.41%   768.26kB 22.16%  go.uber.org/zap.Config.Build
         0     0% 33.41%   768.26kB 22.16%  go.uber.org/zap.Config.buildOptions.func1
         0     0% 33.41%   768.26kB 22.16%  go.uber.org/zap.New
         0     0% 33.41%   768.26kB 22.16%  go.uber.org/zap.WrapCore.func1
         0     0% 33.41%   768.26kB 22.16%  go.uber.org/zap.optionFunc.apply
         0     0% 33.41%   768.26kB 22.16%  go.uber.org/zap/zapcore.NewSamplerWithOptions
         0     0% 33.41%   768.26kB 22.16%  main.main
         0     0% 33.41% -1415.71kB 40.84%  net/http.(*ServeMux).ServeHTTP
         0     0% 33.41% -1415.71kB 40.84%  net/http.(*conn).serve
         0     0% 33.41% -1415.71kB 40.84%  net/http.HandlerFunc.ServeHTTP
         0     0% 33.41% -1415.71kB 40.84%  net/http.serverHandler.ServeHTTP
         0     0% 33.41% -1415.71kB 40.84%  net/http/pprof.Index
         0     0% 33.41% -1415.71kB 40.84%  net/http/pprof.collectProfile
         0     0% 33.41% -1415.71kB 40.84%  net/http/pprof.handler.ServeHTTP
         0     0% 33.41% -1415.71kB 40.84%  net/http/pprof.handler.serveDeltaProfile
         0     0% 33.41%  -512.05kB 14.77%  runtime.bgscavenge
         0     0% 33.41%   515.19kB 14.86%  runtime.doInit
         0     0% 33.41%   515.19kB 14.86%  runtime.doInit1
         0     0% 33.41%  1283.45kB 37.02%  runtime.main
         0     0% 33.41%     -513kB 14.80%  runtime.mcall
         0     0% 33.41%     -513kB 14.80%  runtime.mstart
         0     0% 33.41%     -513kB 14.80%  runtime.mstart0
         0     0% 33.41%     -513kB 14.80%  runtime.mstart1
         0     0% 33.41%    -1026kB 29.60%  runtime.newm
         0     0% 33.41%     -513kB 14.80%  runtime.park_m
         0     0% 33.41%    -1026kB 29.60%  runtime.resetspinning
         0     0% 33.41%    -1026kB 29.60%  runtime.schedule
         0     0% 33.41%    -1026kB 29.60%  runtime.startm
         0     0% 33.41%    -1026kB 29.60%  runtime.wakep
         0     0% 33.41% -1415.71kB 40.84%  runtime/pprof.(*Profile).WriteTo
         0     0% 33.41% -1415.71kB 40.84%  runtime/pprof.(*profileBuilder).build
         0     0% 33.41% -1415.71kB 40.84%  runtime/pprof.writeHeap
         0     0% 33.41% -1415.71kB 40.84%  runtime/pprof.writeHeapInternal
         0     0% 33.41% -1415.71kB 40.84%  runtime/pprof.writeHeapProto
         0     0% 33.41%   515.19kB 14.86%  unicode.init
```