(function () {
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

    document.addEventListener('DOMContentLoaded', initSidebarInteractions);
})();