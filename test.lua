local log = require("log")

--local source = "docker://localhost:5000/alpine:3"
local source = "docker.pkg.github.com/pujo-j/ocilot-builders/python-base:latest"
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

local config = new:config()
local e = config["Env"]
e[#e] = 'TEST=TEST_ENV_VAR'
log.info("Applying config", { config = config })
new:setConfig(config)

log.info("Pushing image", { destination = destination })
new:push()

