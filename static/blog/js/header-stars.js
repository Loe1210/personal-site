(function () {
    'use strict';

    function addHeaderStars() {
        var header = document.querySelector('.blog-hero, .directory-hero');
        if (!header || header.querySelector('.header-star-field')) return;

        var field = document.createElement('div');
        field.className = 'header-star-field';
        field.setAttribute('aria-hidden', 'true');

        // A dense dot canopy with deterministic, individually staggered breathing.
        for (var index = 0; index < 960; index += 1) {
            var star = document.createElement('i');
            var column = index % 60;
            var row = Math.floor(index / 60);
            var x = 1 + column * 1.66;
            var y = 3 + row * 6.2;
            star.className = 'header-star';
            star.style.setProperty('--star-x', x + '%');
            star.style.setProperty('--star-y', y + '%');
            star.style.setProperty('--star-size', '1px');
            star.style.setProperty('--star-color', index % 5 === 0 ? 'rgba(183, 225, 255, .94)' : 'rgba(132, 195, 240, .86)');
            star.style.setProperty('--star-duration', (4.4 + (index % 13) * 0.37) + 's');
            star.style.setProperty('--star-delay', (-((index * 0.53) % 7.3)) + 's');
            field.appendChild(star);
        }

        header.prepend(field);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', addHeaderStars);
    } else {
        addHeaderStars();
    }
}());