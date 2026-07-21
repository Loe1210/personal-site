(function () {
    function decorateProfileText() {
        Array.prototype.forEach.call(document.querySelectorAll('.blog-profile-name, .blog-profile-tagline, .blog-profile-intro'), function (element) {
            if (element.dataset.characterDecorated === 'true') return;
            var text = element.textContent || '';
            element.dataset.characterDecorated = 'true';
            element.setAttribute('aria-label', text);
            element.textContent = '';
            Array.from(text).forEach(function (character) {
                if (/\s/.test(character)) {
                    element.appendChild(document.createTextNode(character));
                    return;
                }
                var letter = document.createElement('span');
                letter.className = 'blog-profile-character';
                letter.textContent = character;
                letter.setAttribute('aria-hidden', 'true');
                element.appendChild(letter);
            });
        });
    }
    function initProfileCharacterHops() {
        Array.prototype.forEach.call(document.querySelectorAll('.blog-profile-character'), function (letter) {
            letter.addEventListener('pointerenter', function () {
                letter.classList.remove('is-falling');
                letter.classList.add('is-hovered');
            });
            letter.addEventListener('pointerleave', function () {
                letter.classList.remove('is-hovered', 'is-falling');
                // Reflow restarts only this character's release animation.
                void letter.offsetWidth;
                letter.classList.add('is-falling');
            });
            letter.addEventListener('animationend', function (event) {
                if (event.animationName === 'profileCharacterDrop') {
                    letter.classList.remove('is-falling');
                }
            });
        });
    }
    function initSidebarNavigationIndicator() {
        var navigation = document.querySelector('.blog-profile-nav');
        if (!navigation || navigation.querySelector('.blog-profile-nav__indicator')) return;

        var indicator = document.createElement('span');
        indicator.className = 'blog-profile-nav__indicator';
        indicator.setAttribute('aria-hidden', 'true');
        navigation.insertBefore(indicator, navigation.firstChild);

        var compact = window.matchMedia ? window.matchMedia('(max-width: 980px)') : null;
        function moveIndicator(item) {
            if (!item || (compact && compact.matches)) return;
            indicator.style.height = item.offsetHeight + 'px';
            indicator.style.transform = 'translateY(' + item.offsetTop + 'px)';
            indicator.classList.add('is-visible');
        }
        function moveToActive() {
            moveIndicator(navigation.querySelector('.blog-profile-nav__item.is-active') || navigation.querySelector('.blog-profile-nav__item'));
        }

        Array.prototype.forEach.call(navigation.querySelectorAll('.blog-profile-nav__item'), function (item) {
            item.addEventListener('pointerenter', function () { moveIndicator(item); });
            item.addEventListener('focus', function () { moveIndicator(item); });
        });
        navigation.addEventListener('pointerleave', moveToActive);
        navigation.addEventListener('focusout', function () {
            window.setTimeout(function () {
                if (!navigation.contains(document.activeElement)) moveToActive();
            }, 0);
        });
        window.addEventListener('resize', moveToActive);
        if (compact) {
            compact.addEventListener('change', function () {
                indicator.classList.toggle('is-visible', !compact.matches);
                if (!compact.matches) moveToActive();
            });
        }
        moveToActive();
    }

    function initSidebarInteractions() {
        var modal = document.getElementById('blogWechatModal');
        var openTrigger = document.querySelector('[data-wechat-open]');
        var overlay = document.querySelector('[data-wechat-overlay]');
        var closeButton = document.querySelector('[data-wechat-close]');
        var lastActiveElement = null;

        if (!modal || !openTrigger) {
            return;
        }

        function setModalOpen(isOpen) {
            modal.hidden = !isOpen;
            document.body.classList.toggle('blog-wechat-modal-open', isOpen);
            openTrigger.setAttribute('aria-expanded', String(isOpen));

            if (isOpen) {
                lastActiveElement = document.activeElement;
                if (closeButton) {
                    window.setTimeout(function () {
                        closeButton.focus();
                    }, 0);
                }
            } else if (lastActiveElement && lastActiveElement.focus) {
                lastActiveElement.focus();
            }
        }

        openTrigger.addEventListener('click', function () {
            setModalOpen(true);
        });

        [overlay, closeButton].forEach(function (node) {
            if (!node) {
                return;
            }
            node.addEventListener('click', function () {
                setModalOpen(false);
            });
        });

        document.addEventListener('keydown', function (event) {
            if (event.key === 'Escape' && !modal.hidden) {
                setModalOpen(false);
            }
        });
    }

    document.addEventListener('DOMContentLoaded', function () {
        decorateProfileText();
        initProfileCharacterHops();
        initSidebarNavigationIndicator();
        initSidebarInteractions();
    });
})();