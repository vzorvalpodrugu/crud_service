-- счётчик для уникальных имён постов
local counter = 0

-- функция вызывается для каждого запроса
request = function()


    local author_id = 1

    local body = string.format(
        '{"name": "Пост %d", "text": "Нагрузочный тест", "author_id": %d}',
        counter,
        author_id
    )

    local headers = {
        ["Content-Type"] = "application/json",
        ["X-User-ID"]    = tostring(author_id),
    }

    return wrk.format("POST", "/posts", headers, body)
end

-- функция вызывается после завершения теста
done = function(summary, latency, requests)
    print("\n---- Результаты ----")
    print(string.format("Запросов всего:     %d", summary.requests))
    print(string.format("Ошибок:             %d", summary.errors.status))
    print(string.format("RPS:                %.2f", summary.requests / (summary.duration / 1e6)))
    print(string.format("Латентность p50:    %.2fms", latency:percentile(50) / 1e3))
    print(string.format("Латентность p95:    %.2fms", latency:percentile(95) / 1e3))
    print(string.format("Латентность p99:    %.2fms", latency:percentile(99) / 1e3))
    print(string.format("Макс латентность:   %.2fms", latency.max / 1e3))
end