local log = require("log")

local source = "docker://localhost:5000/alpine:3"
local destination = "docker://localhost:5000/test:latest"
local cachePrefix = "localhost:5000/cache/"

local alpine = pull(source)
log.info("Loaded image", { image = alpine:__tostring() })

log.info("Cloning image")
local new = alpine:clone(destination)

applyCached(new, cachePrefix, "test:snapshot1", function()
    local files = snapshot("./script_interface")
    local layer = files:asLayer("./temp.tar")
    return { layer }
end)

log.info("Pushing image", { destination = destination })
new:push()
