const snowBuildUpCanvas = document.getElementById("snow-build-up");
const ctx = snowBuildUpCanvas.getContext("2d");
const snowflakes = [];

snowBuildUpCanvas.width = getViewportWidth();
snowBuildUpCanvas.height = document.documentElement.scrollHeight;

let snowflakeCharacters = ["❄", "*", "❉", "❃", "❅"];
const snowLevels = new Array(snowBuildUpCanvas.width).fill(0);
const maxSnowDepth = 40;

let snowflakeCounter = 0;
let snowflakeCreationRate;
let isMelting = false;

const SNOW_ENABLED_COOKIE = "snow_enabled";
const SNOW_DENSITY_COOKIE = "snow_density";

function ensureSnowCookiesExist() {
    const oneYearSeconds = 60 * 60 * 24 * 365;
    if (getCookie(SNOW_ENABLED_COOKIE) === "") {
        setCookie(SNOW_ENABLED_COOKIE, "1", oneYearSeconds);
    }
    if (getCookie(SNOW_DENSITY_COOKIE) === "") {
        setCookie(SNOW_DENSITY_COOKIE, "1", oneYearSeconds);
    }
}

function parseSnowDensity(rawValue) {
    const v = (rawValue || "").trim();
    const parsed = Number(v);

    if (!Number.isFinite(parsed)) return 1;

    // Clamp to avoid pathological values.
    return Math.min(5, Math.max(0.1, parsed));
}

function computeBaseCreationRate() {
    if (window.innerWidth > 768) {
        // More snowflakes on desktop
        return 4;
    }
    // Fewer snowflakes on mobile
    return 8;
}

function applyDensityToCreationRate(baseRate, density) {
    // Lower rate -> more flakes
    // Density > 1 should increase snowfall
    return Math.max(1, Math.round(baseRate / density));
}

ensureSnowCookiesExist();

let snowEnabled = parseBooleanFromString(getCookie(SNOW_ENABLED_COOKIE));
let snowDensity = parseSnowDensity(getCookie(SNOW_DENSITY_COOKIE));
let lastFrameTime = performance.now();

snowflakeCreationRate = applyDensityToCreationRate(computeBaseCreationRate(), snowDensity);

function createSnowflake() {
    const snowflake = {
        x: Math.random() * snowBuildUpCanvas.width,
        y: -10,
        size: Math.random() * 20 + 10,
        speed: Math.random() * 1.5 + 1,
        opacity: 1,
        character: snowflakeCharacters[Math.floor(Math.random() * snowflakeCharacters.length)]
    };
    snowflakes.push(snowflake);
}

function smoothAccumulation() {
    // Smooth out the snow accumulation over time by averaging snow levels
    for (let i = 1; i < snowLevels.length - 1; i++) {
        snowLevels[i] = (snowLevels[i - 1] + snowLevels[i] + snowLevels[i + 1]) / 3;
    }
}

function drawSnowflakes(deltaTime) {
    ctx.clearRect(0, 0, snowBuildUpCanvas.width, snowBuildUpCanvas.height);

    // Draw buildup
    ctx.fillStyle = "white";
    for (let x = 0; x < snowLevels.length; x++) {
        ctx.fillRect(x, snowBuildUpCanvas.height - snowLevels[x], 1, snowLevels[x]);
    }

    // Draw falling snowflakes
    ctx.fillStyle = "white";
    snowflakes.forEach((snowflake, index) => {
        ctx.globalAlpha = snowflake.opacity;
        ctx.font = `${snowflake.size}px sans-serif`;
        ctx.fillText(snowflake.character, snowflake.x, snowflake.y);
        ctx.globalAlpha = 1;

        // Adjust speed by deltaTime
        snowflake.y += snowflake.speed * deltaTime * 60;

        // If snowflake reaches the bottom, accumulate snow
        if (snowflake.y + snowflake.size / 2 >= snowBuildUpCanvas.height - snowLevels[Math.floor(snowflake.x)]) {
            const snowflakeX = Math.floor(snowflake.x);
            const snowflakeSize = snowflake.size / 2;

            // Accumulate snow more smoothly across a range with a weighted falloff effect
            const accumulationWidth = Math.ceil(snowflakeSize * 2);
            for (let i = -accumulationWidth; i <= accumulationWidth; i++) {
                const xIndex = Math.min(Math.max(snowflakeX + i, 0), snowLevels.length - 1);
                const distance = Math.abs(i); // The distance from the center of the snowflake
                const falloff = Math.exp(-distance / 5); // Gaussian falloff (smooth transition)

                // Apply the falloff to the accumulation to smooth the edges
                if (snowLevels[xIndex] < maxSnowDepth) {
                    snowLevels[xIndex] += (snowflake.size / 4) * falloff; // Quicker accumulation with smoothing
                }
            }

            // Remove snowflake once it lands
            snowflakes.splice(index, 1);
        }
    });

    // Apply smoothing after every update
    smoothAccumulation();
}

function meltSnow() {
    let allMelted = true;
    for (let i = 0; i < snowLevels.length; i++) {
        if (snowLevels[i] > 0) {
            const meltRate = Math.random() * 0.3 + 0.1;
            snowLevels[i] -= meltRate;
            if (snowLevels[i] < 0) snowLevels[i] = 0;
            allMelted = false;
        }
    }

    // Introduce additional random offsets to simulate flow and uneven melting
    for (let i = 0; i < snowLevels.length; i++) {
        if (i > 0 && i < snowLevels.length - 1) {
            const leftNeighbor = snowLevels[i - 1];
            const rightNeighbor = snowLevels[i + 1];
            const average = (leftNeighbor + rightNeighbor) / 2;

            // Adjust current level slightly toward the average of neighbors
            const adjustment = (average - snowLevels[i]) * 0.2;
            snowLevels[i] += adjustment;
        }
    }

    return allMelted;
}

function checkFullSnow() {
    return snowLevels.every((level) => level >= maxSnowDepth);
}

function animate() {
    const currentTime = performance.now();
    const deltaTime = (currentTime - lastFrameTime) / 1000;
    lastFrameTime = currentTime;

    // Throttle cookie reads
    if (!animate._lastPrefCheck || currentTime - animate._lastPrefCheck > 500) {
        animate._lastPrefCheck = currentTime;
        snowEnabled = parseBooleanFromString(getCookie(SNOW_ENABLED_COOKIE));
        snowDensity = parseSnowDensity(getCookie(SNOW_DENSITY_COOKIE));
        snowflakeCreationRate = applyDensityToCreationRate(computeBaseCreationRate(), snowDensity);
    }

    if (!snowEnabled) {
        snowflakes.length = 0;
        snowLevels.fill(0);
        isMelting = false;
        ctx.clearRect(0, 0, snowBuildUpCanvas.width, snowBuildUpCanvas.height);
        requestAnimationFrame(animate);
        return;
    }

    if (!isMelting) {
        if (checkFullSnow()) {
            isMelting = true;
        } else if (snowflakeCounter % snowflakeCreationRate === 0) {
            createSnowflake();
        }
        snowflakeCounter++;
    } else {
        if (meltSnow()) {
            // Reset state when snow is fully melted
            snowflakes.length = 0;
            snowLevels.fill(0);
            isMelting = false;
        }
    }

    drawSnowflakes(deltaTime);
    requestAnimationFrame(animate);
}

function loadSnowState() {
    if (!snowEnabled) {
        return;
    }
    const savedState = sessionStorage.getItem("snowState");
    if (!savedState) {
        return;
    }

    const state = JSON.parse(savedState);

    // Check if saved state matches current canvas width
    if (state.canvasWidth !== snowBuildUpCanvas.width) {
        return;
    }

    // Restore snow levels
    for (let i = 0; i < Math.min(state.snowLevels.length, snowLevels.length); i++) {
        snowLevels[i] = state.snowLevels[i];
    }

    // Restore falling snowflakes
    state.snowflakes.forEach((flake) => {
        // Offset y position to account for time passed since save
        const timeSinceSave = (Date.now() - state.timestamp) / 1000;
        flake.y += flake.speed * timeSinceSave * 60;
        // Add snowflake back to array
        snowflakes.push({ ...flake });
    });

    // Restore other state
    isMelting = state.isMelting || false;
    snowflakeCounter = state.snowflakeCounter || 0;
}

function saveSnowState() {
    if (!snowEnabled) {
        sessionStorage.removeItem("snowState");
        return;
    }
    const state = {
        snowLevels: Array.from(snowLevels),
        snowflakes: snowflakes.map((flake) => ({ ...flake })),
        snowflakeCounter: snowflakeCounter,
        canvasWidth: snowBuildUpCanvas.width,
        isMelting: isMelting,
        timestamp: Date.now()
    };
    sessionStorage.setItem("snowState", JSON.stringify(state));
}

// Load state on initialization
loadSnowState();

// Save state before page unload
window.addEventListener("beforeunload", () => {
    saveSnowState();
});

// Update array on window resize
window.addEventListener("resize", () => {
    snowBuildUpCanvas.width = getViewportWidth();
    snowBuildUpCanvas.height = document.documentElement.scrollHeight;

    // Recompute snowfall rate when viewport changes.
    snowflakeCreationRate = applyDensityToCreationRate(computeBaseCreationRate(), snowDensity);

    // Don't reset snow levels on resize, let it adapt
    if (snowLevels.length === snowBuildUpCanvas.width) {
        return;
    }

    const oldLevels = Array.from(snowLevels);
    snowLevels.length = snowBuildUpCanvas.width;
    snowLevels.fill(0);

    // Try to preserve some state if resizing
    for (let i = 0; i < Math.min(oldLevels.length, snowLevels.length); i++) {
        snowLevels[i] = oldLevels[i];
    }
});

// Update canvas height on scroll
window.addEventListener("scroll", () => {
    snowBuildUpCanvas.height = document.documentElement.scrollHeight;
});

animate();
