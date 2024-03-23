local tag = 0
function T(l, v, r)
	tag = tag +  1
    return {left = l, value = v, right = r, tag = tag}
end

local e = nil
local t1 = T(T(T(e, 1, e), 2, T(e, 3, e)), 4, T(e, 5, e))
local t2 = T(e, 1, T(e, 2, T(e, 3, T(e, 4, T(e, 5, e)))))
local t3 = T(e, 1, T(e, 2, T(e, 3, T(e, 4, T(e, 6, e)))))


function Visit(t)
    if t ~= nil then  -- note: ~= is "not equal"
        Visit(t.left)
        coroutine.yield(t.value)
        Visit(t.right)
    end
end

function Cmp(t1, t2)
    local co1 = coroutine.create(Visit)
    local co2 = coroutine.create(Visit)
    while true
    do
        local ok1, v1 = coroutine.resume(co1, t1)
        local ok2, v2 = coroutine.resume(co2, t2)
        if ok1 ~= ok2 or v1 ~= v2 then
            return false
        end
        if not ok1 and not ok2 then
            return true
        end
    end
end

function Cmp(t1, t2)
    local next1 = coroutine.wrap(function() Visit(t1) end)
    local next2 = coroutine.wrap(function() Visit(t2) end)
    while true
    do
        local v1 = next1()
        local v2 = next2()
        if v1 ~= v2 then
            return false
        end
        if v1 == nil and v2 == nil then
            return true
        end
    end
end


print(Cmp(t1, t3))
