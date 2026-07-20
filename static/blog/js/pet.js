(function () {
    var DEFAULT_COLS = 8;
    var DEFAULT_ROWS = 9;
    var DISPLAY_SCALE = 0.58;
    var PET_MANIFEST_CANDIDATES = ['/pet/index.json', '/static/pet/index.json'];
    var STORAGE_KEY = 'hins-blog-active-pet';

    var ROWS = {
        idle: 0,
        runRight: 1,
        runLeft: 2,
        wave: 3,
        jump: 4
    };

    var runtime = {
        manifest: null,
        manifestBasePath: '',
        activePetId: '',
        cleanup: null,
        renderToken: 0
    };

    function init() {
        if (window.__hinsBlogPetRuntimeBooted) return;
        window.__hinsBlogPetRuntimeBooted = true;
        injectPetStyles();
        var mount = ensurePetMount();
        if (!mount) return;

        loadPetManifest()
            .then(function (payload) {
                runtime.manifest = normalizeManifest(payload.manifest);
                runtime.manifestBasePath = payload.basePath;
                runtime.activePetId = getSavedPetId(runtime.manifest) || runtime.manifest.defaultPetId;
                return switchToPet(runtime.activePetId, mount, false);
            })
            .catch(function () {
                mount.setAttribute('data-pet-fallback', 'true');
            });
    }

    function ensurePetMount() {
        cleanupDuplicateMounts();
        var mount = document.getElementById('blogSidebarPet');
        if (mount) return mount;

        var dock = document.createElement('div');
        dock.className = 'blog-global-pet-dock';
        mount = document.createElement('div');
        mount.id = 'blogSidebarPet';
        mount.className = 'blog-sidebar-pet blog-sidebar-pet--global';
        mount.setAttribute('aria-hidden', 'true');
        dock.appendChild(mount);
        document.body.appendChild(dock);
        cleanupDuplicateMounts();
        return mount;
    }

    function cleanupDuplicateMounts() {
        var primary = document.getElementById('blogSidebarPet');
        var mounts = document.querySelectorAll('#blogSidebarPet');
        for (var i = 0; i < mounts.length; i++) {
            if (mounts[i] !== primary && mounts[i].parentNode) {
                mounts[i].parentNode.removeChild(mounts[i]);
            }
        }

        var docks = document.querySelectorAll('.blog-global-pet-dock');
        for (var j = 0; j < docks.length; j++) {
            if ((!primary || !docks[j].contains(primary)) && docks[j].parentNode) {
                docks[j].parentNode.removeChild(docks[j]);
            }
        }
    }

    function loadPetManifest() {
        var attempt = 0;

        function next() {
            if (attempt >= PET_MANIFEST_CANDIDATES.length) {
                return Promise.reject(new Error('pet manifest unavailable'));
            }
            var manifestPath = PET_MANIFEST_CANDIDATES[attempt++];
            return fetch(manifestPath)
                .then(function (response) {
                    if (!response.ok) throw new Error('pet manifest unavailable');
                    return response.json();
                })
                .then(function (manifest) {
                    return {
                        manifest: manifest,
                        basePath: manifestPath.replace(/\/index\.json$/, '')
                    };
                })
                .catch(next);
        }

        return next();
    }

    function normalizeManifest(manifest) {
        var pets = Array.isArray(manifest && manifest.pets) ? manifest.pets : [];
        pets = pets.filter(function (pet) { return pet && pet.id && pet.metaPath; });
        if (!pets.length) throw new Error('empty pet manifest');
        return {
            defaultPetId: manifest.defaultPetId || pets[0].id,
            pets: pets
        };
    }

    function getSavedPetId(manifest) {
        try {
            var saved = window.localStorage && window.localStorage.getItem(STORAGE_KEY);
            if (saved && findPet(manifest, saved)) return saved;
        } catch (error) {}
        return '';
    }

    function savePetId(id) {
        try {
            if (window.localStorage) window.localStorage.setItem(STORAGE_KEY, id);
        } catch (error) {}
    }

    function switchToPet(id, mount, persist) {
        mount = mount || document.getElementById('blogSidebarPet');
        if (!mount || !runtime.manifest) return Promise.reject(new Error('pet runtime unavailable'));

        var token = ++runtime.renderToken;
        var pet = findPet(runtime.manifest, id) || findPet(runtime.manifest, runtime.manifest.defaultPetId) || runtime.manifest.pets[0];
        runtime.activePetId = pet.id;
        if (persist !== false) savePetId(pet.id);

        if (typeof runtime.cleanup === 'function') {
            runtime.cleanup();
            runtime.cleanup = null;
        }
        mount.classList.remove('is-visible', 'is-ready');
        mount.classList.add('is-loading');
        mount.innerHTML = '';

        return loadPetMeta(pet).then(function (payload) {
            return loadImage(payload.basePath + '/' + payload.meta.spritesheetPath).then(function (image) {
                if (token !== runtime.renderToken) return null;
                mount.classList.remove('is-loading');
                runtime.cleanup = renderPet(mount, image, payload.meta, pet);
                return runtime.cleanup;
            });
        }).catch(function (error) {
            if (token !== runtime.renderToken) return null;
            mount.classList.remove('is-loading');
            mount.setAttribute('data-pet-fallback', 'true');
            throw error;
        });
    }

    function findPet(manifest, id) {
        for (var i = 0; i < manifest.pets.length; i++) {
            if (manifest.pets[i].id === id) return manifest.pets[i];
        }
        return null;
    }

    function nextPetId(currentId) {
        var pets = runtime.manifest.pets;
        for (var i = 0; i < pets.length; i++) {
            if (pets[i].id === currentId) return pets[(i + 1) % pets.length].id;
        }
        return pets[0].id;
    }

    function loadPetMeta(pet) {
        var metaPath = resolvePetPath(pet.metaPath);
        return fetch(metaPath)
            .then(function (response) {
                if (!response.ok) throw new Error('pet meta unavailable');
                return response.json();
            })
            .then(function (meta) {
                return {
                    meta: meta,
                    basePath: metaPath.replace(/\/pet\.json$/, '')
                };
            });
    }

    function resolvePetPath(path) {
        if (/^https?:\/\//.test(path) || path.charAt(0) === '/') return path;
        if (path.indexOf('pet/') === 0) return '/' + path;
        return runtime.manifestBasePath + '/' + path;
    }

    function loadImage(src) {
        return new Promise(function (resolve, reject) {
            var image = new Image();
            image.decoding = 'async';
            image.onload = function () {
                if (image.decode) {
                    image.decode().then(function () { resolve(image); }).catch(function () { resolve(image); });
                } else {
                    resolve(image);
                }
            };
            image.onerror = reject;
            image.src = src;
        });
    }

    function renderPet(mount, image, meta, pet) {
        var cols = Math.max(1, Number(meta.columns || DEFAULT_COLS));
        var rows = Math.max(1, Number(meta.rows || DEFAULT_ROWS));
        var frameWidth = Math.floor(image.naturalWidth / cols);
        var frameHeight = Math.floor(image.naturalHeight / rows);
        var sprite = document.createElement('div');
        var state = 'idle';
        var frame = 0;
        var frameInterval = 420;
        var lastFrameAt = 0;
        var rafId = 0;
        var idleActionTimer = 0;
        var oneShot = null;
        var dragging = false;
        var floating = mount.classList.contains('is-floating');
        var dragOffsetX = 0;
        var dragOffsetY = 0;
        var lastPointerX = 0;
        var lastDirection = 'run-right';
        var frameCounts = {};
        var eventCleanups = [];

        ROWS.idle = rowNumber(meta.idleRow, ROWS.idle, rows);
        ROWS.runRight = rowNumber(meta.runningRightRow, ROWS.runRight, rows);
        ROWS.runLeft = rowNumber(meta.runningLeftRow, ROWS.runLeft, rows);
        ROWS.jump = rowNumber(meta.jumpingRow, ROWS.jump, rows);
        ROWS.wave = rowNumber(meta.wavingRow, ROWS.wave, rows);

        frameCounts.idle = frameCount(meta.idleFrames, cols);
        frameCounts.runRight = frameCount(meta.runningRightFrames, cols);
        frameCounts.runLeft = frameCount(meta.runningLeftFrames, frameCounts.runRight);
        frameCounts.jump = frameCount(meta.jumpingFrames, cols);
        frameCounts.wave = frameCount(meta.wavingFrames, cols);

        if (!frameWidth || !frameHeight) throw new Error('invalid pet frame size');

        mount.removeAttribute('data-pet-fallback');
        mount.classList.add('is-ready');
        mount.dataset.petId = pet.id;
        mount.setAttribute('role', 'button');
        mount.setAttribute('tabindex', '0');
        mount.setAttribute('aria-label', 'drag, hover, or switch the pet');

        sprite.className = 'blog-sidebar-pet__sprite';
        sprite.style.backgroundImage = 'url(' + image.src + ')';
        sprite.style.backgroundSize = (frameWidth * cols * DISPLAY_SCALE) + 'px ' + (frameHeight * rows * DISPLAY_SCALE) + 'px';
        sprite.style.width = Math.round(frameWidth * DISPLAY_SCALE) + 'px';
        sprite.style.height = Math.round(frameHeight * DISPLAY_SCALE) + 'px';
        mount.appendChild(sprite);
        draw();
        mount.classList.add('is-visible');

        function on(target, type, handler) {
            target.addEventListener(type, handler);
            eventCleanups.push(function () { target.removeEventListener(type, handler); });
        }

        function setState(nextState, interval, options) {
            options = options || {};
            if (state !== nextState) {
                frame = 0;
                lastFrameAt = 0;
            }
            state = nextState;
            frameInterval = interval;
            oneShot = options.once ? { loops: 0, after: options.after || 'idle' } : null;
            draw();
        }

        function rowForState() {
            if (state === 'run-right') return ROWS.runRight;
            if (state === 'run-left') return ROWS.runLeft;
            if (state === 'jump') return ROWS.jump;
            if (state === 'wave') return ROWS.wave;
            return ROWS.idle;
        }

        function frameLimitForState() {
            if (state === 'run-right') return frameCounts.runRight;
            if (state === 'run-left') return frameCounts.runLeft;
            if (state === 'jump') return frameCounts.jump;
            if (state === 'wave') return frameCounts.wave;
            return frameCounts.idle;
        }

        function draw() {
            sprite.style.backgroundPosition = (-frame * frameWidth * DISPLAY_SCALE) + 'px ' + (-rowForState() * frameHeight * DISPLAY_SCALE) + 'px';
        }

        function tick(now) {
            if (!lastFrameAt) lastFrameAt = now;
            if (now - lastFrameAt >= frameInterval) {
                lastFrameAt = now;
                frame = (frame + 1) % frameLimitForState();
                if (oneShot && frame === 0) {
                    oneShot.loops += 1;
                    if (oneShot.loops >= 1) setState(oneShot.after, 420);
                }
                draw();
            }
            rafId = window.requestAnimationFrame(tick);
        }

        function playOnce(nextState, interval, afterState) {
            setState(nextState, interval, { once: true, after: afterState || 'idle' });
        }

        function liftToViewport(rect) {
            if (!floating) {
                document.body.appendChild(mount);
                floating = true;
            }
            mount.classList.add('is-floating');
            mount.style.left = rect.left + 'px';
            mount.style.top = rect.top + 'px';
            mount.style.width = rect.width + 'px';
            mount.style.height = rect.height + 'px';
        }

        function beginDrag(event) {
            var rect = mount.getBoundingClientRect();
            dragging = true;
            dragOffsetX = event.clientX - rect.left;
            dragOffsetY = event.clientY - rect.top;
            lastPointerX = event.clientX;
            lastDirection = 'run-right';
            liftToViewport(rect);
            window.clearTimeout(idleActionTimer);
            mount.classList.add('is-dragging');
            if (mount.setPointerCapture) mount.setPointerCapture(event.pointerId);
            setState(lastDirection, 86);
            event.preventDefault();
        }

        function moveDrag(event) {
            if (!dragging) return;
            var x = clamp(event.clientX - dragOffsetX, 0, window.innerWidth - mount.offsetWidth);
            var y = clamp(event.clientY - dragOffsetY, 0, window.innerHeight - mount.offsetHeight);
            var dx = event.clientX - lastPointerX;
            mount.style.left = x + 'px';
            mount.style.top = y + 'px';
            if (Math.abs(dx) > 2) {
                var direction = dx < 0 ? 'run-left' : 'run-right';
                if (direction !== lastDirection) {
                    lastDirection = direction;
                    setState(direction, 78);
                }
                lastPointerX = event.clientX;
            }
        }

        function endDrag(event) {
            if (!dragging) return;
            dragging = false;
            mount.classList.remove('is-dragging');
            if (mount.releasePointerCapture) {
                try { mount.releasePointerCapture(event.pointerId); } catch (error) {}
            }
            setState('idle', 420);
            scheduleIdleAction();
        }

        function scheduleIdleAction() {
            window.clearTimeout(idleActionTimer);
            idleActionTimer = window.setTimeout(function () {
                if (!dragging && state === 'idle') playOnce('wave', 120, 'idle');
                scheduleIdleAction();
            }, 7000 + Math.round(Math.random() * 4000));
        }

        on(mount, 'pointerdown', beginDrag);
        on(mount, 'pointermove', moveDrag);
        on(mount, 'pointerup', endDrag);
        on(mount, 'pointercancel', endDrag);
        on(mount, 'pointerenter', function () {
            if (!dragging && state !== 'jump') {
                window.clearTimeout(idleActionTimer);
                playOnce('jump', 78, 'idle');
                scheduleIdleAction();
            }
        });
        on(mount, 'keydown', function (event) {
            if (event.key === 'Enter' || event.key === ' ') {
                event.preventDefault();
                playOnce('jump', 78, 'idle');
            }
        });
        on(mount, 'dblclick', function (event) {
            event.preventDefault();
            switchToPet(nextPetId(runtime.activePetId), mount, true);
        });

        setState('idle', 420);
        rafId = window.requestAnimationFrame(tick);
        scheduleIdleAction();

        return function cleanupPet() {
            window.cancelAnimationFrame(rafId);
            window.clearTimeout(idleActionTimer);
            eventCleanups.forEach(function (cleanup) { cleanup(); });
        };
    }

    function injectPetStyles() {
        if (document.getElementById('blogPetRuntimeStyles')) return;
        var style = document.createElement('style');
        style.id = 'blogPetRuntimeStyles';
        style.textContent = '' +
            '.blog-global-pet-dock{position:fixed;left:18px;bottom:18px;z-index:140;pointer-events:none}' +
            '.blog-sidebar-pet--global{min-width:112px;min-height:128px;display:flex;align-items:flex-end;justify-content:center;pointer-events:auto}' +
            '.blog-sidebar-pet{position:relative}' +
            '.blog-sidebar-pet__sprite{background-repeat:no-repeat;image-rendering:auto;filter:drop-shadow(0 10px 18px rgba(0,0,0,.28));transform-origin:center bottom}' +            '.blog-sidebar-pet.is-floating{position:fixed;z-index:140;pointer-events:auto;overflow:visible}' +
            '.blog-sidebar-pet.is-dragging{cursor:grabbing}' +
            '.blog-sidebar-pet{opacity:0;overflow:visible;transition:opacity .18s ease;touch-action:none;user-select:none}' +
            '.blog-sidebar-pet.is-visible,.blog-sidebar-pet[data-pet-fallback="true"]{opacity:1}' +
            '.blog-sidebar-pet__sprite{animation:none!important;backface-visibility:hidden;transform:translateZ(0);will-change:background-position,transform}' +
            '.blog-sidebar-pet.is-dragging .blog-sidebar-pet__sprite{transform:translate3d(0,-2px,0)}';
        document.head.appendChild(style);
    }

    function frameCount(value, fallback) {
        var number = Number(value);
        if (!Number.isFinite(number)) number = fallback;
        return Math.max(1, Math.min(DEFAULT_COLS, Math.floor(number)));
    }

    function rowNumber(value, fallback, rowCount) {
        var number = Number(value);
        if (!Number.isFinite(number)) number = fallback;
        return Math.max(0, Math.min(rowCount - 1, number));
    }

    function clamp(value, min, max) {
        return Math.min(Math.max(value, min), max);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();