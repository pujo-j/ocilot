local log = require("log")
local box=require("luabox")

-- Layer

local layermt = {}

local enrichLayer = function(layer)
    setmetatable(layer, layermt)
    return layer
end

layermt.__index = function(_, key)
    return layermt[key]
end

layermt.info = function(this)
    return ocisys.layerInfo(this._ud)
end

layermt.__tostring = function(this)
    return ocisys.layerString(this._ud)
end


-- Snapshot
local snapmt = {}

local enrichSnapshot = function(snap)
    setmetatable(snap, snapmt)
    return snap
end

snapmt.__index = function(_, key)
    return snapmt[key]
end

snapmt.asLayer = function(this, workFile)
    log.debug("Generating layer")
    local layer_ud = ocisys.snapAsLayer(this._ud, workFile)
    return enrichLayer({ _ud = layer_ud })
end

snapmt.diff = function(this, other)
    local snap_ud = ocisys.snapDiff(this._ud, other._ud)
    return enrichSnapshot({ _ud = snap_ud })
end

snapmt.files = function(this)
    return ocisys.snapFiles(this._ud)
end

snapmt.__tostring = function(this)
    return ocisys.snapString(this._ud)
end

-- Image
local imagemt = {}

local function enrichImage(image)
    setmetatable(image, imagemt)
    return image
end

imagemt.__index = function(_, key)
    return imagemt[key]
end

imagemt.push = function(this)
    ocisys.imagePush(this._ud)
end

imagemt.clone = function(this, name)
    local clone_ud = ocisys.imageClone(this._ud, name)
    return enrichImage({ _ud = clone_ud })
end

imagemt.config = function(this)
    return ocisys.imageGetConfig(this._ud)
end

imagemt.setConfig = function(this,config)
    ocisys.imageSetConfig(this._ud,config)
end

imagemt.layers = function(this)
    local layers = ocisys.imageGetLayers(this._ud)
    local res = {}
    for _, l in ipairs(layers) do
        res[#res + 1] = enrichLayer({ _ud = l })
    end
    return res
end

imagemt.append = function(this, layer)
    ocisys.imageAppendLayer(this._ud, layer._ud)
end

imagemt.__tostring = function(this)
    return ocisys.imageString(this._ud)
end


-- Globals

function snapshot(dir)
    log.debug("Generating snapshot ", { dir = dir })
    local snap_ud = ocisys.snap(dir)
    return enrichSnapshot({
        _ud = snap_ud
    })
end

pull = function(name)
    local image_ud = ocisys.imagePull(name)
    return enrichImage({ _ud = image_ud })
end

shell = function(bin, args)
    local shell_cmd = bin or "/bin/sh"
    local shell_args = args or "-c"
    s = {}
    s.exec = function(cmd)
        ocisys.shellExec(shell_cmd, shell_args, cmd)
    end
    s.eval = function(cmd)
        ocisys.shellEval(shell_cmd, shell_args, cmd)
    end
    return s
end

read = function(fileName)
    return ocisys.read(fileName)
end

write = function(fileName, data)
    return ocisys.write(fileName, data)
end

hash = function(data)
    return ocisys.hash(data)
end

applyCached = function(to_image, cache_url, cache_key, buildlayers)
    local image_ud, found = ocisys.cacheGetImage(cache_url, cache_key)
    local image = enrichImage({ _ud = image_ud })
    if found then
        log.debug("cache hit for " .. cache_url .. ":" .. cache_key)
        local layers = image:layers()
        for i = 1, #layers, 1 do
            log.debug("appending cached layer " .. layers[i]:__tostring())
            to_image:append(layers[i])
        end
    else
        log.debug("cache miss for " .. cache_url .. ":" .. cache_key)
        local layers = buildlayers()
        for i = 1, #layers, 1 do
            log.debug("appending layer " .. layers[i]:__tostring())
            image:append(layers[i])
            to_image:append(layers[i])
        end
        log.debug("saving cache image")
        image:push()
    end
end

env=box.getEnv()
args=box.getArgs()