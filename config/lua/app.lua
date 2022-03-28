#!/usr/bin/env tarantool

-- Настроить базу данных
box.cfg {
    listen = 3301
}

-- При поднятии БД создаем спейсы и индексы
box.once('init', function()
    box.schema.space.create('sessions')
    box.schema.user.passwd('pass')
    box.space.sessions:create_index('primary',
        { type = 'TREE', parts = {1, 'string'}})

end)

function check_session(session_id)
    local session_id = box.space.sessions:select{session_id}[1]
    print('found session', session_id)
    return session_id
end