var iUp = (function () {
	var time = 0,
		duration = 150,
		clean = function () {
			time = 0;
		},
		up = function (element) {
			setTimeout(function () {
				element.classList.add("up");
			}, time);
			time += duration;
		},
		down = function (element) {
			element.classList.remove("up");
		},
		toggle = function (element) {
			setTimeout(function () {
				element.classList.toggle("up");
			}, time);
			time += duration;
		};
	return {
		clean: clean,
		up: up,
		down: down,
		toggle: toggle
	};
})();

// Bing image URL pattern: validates format and prevents CSS injection
var BING_IMAGE_URL_PATTERN = /^\/th\?id=OHR\.[a-zA-Z0-9_\-]+\.jpg(&[a-zA-Z0-9=._\-]+)*$/;

function getBingImages(imgUrls) {
	/**
	 * 加载本地随机壁纸
	 * 从 assets/wallpapers 目录随机选择一张本地图片
	 */
	var panel = document.querySelector('#panel');
	if (!panel || !imgUrls || !Array.isArray(imgUrls) || imgUrls.length === 0) {
		return;
	}
	
	var indexName = "home-bg-index";
	var index = parseInt(sessionStorage.getItem(indexName), 10);

	if (isNaN(index) || index >= imgUrls.length) {
		index = Math.floor(Math.random() * imgUrls.length);
		sessionStorage.setItem(indexName, index);
	}
	
	var imgUrl = imgUrls[index];
	if (!imgUrl || typeof imgUrl !== 'string') {
		return;
	}
	
	// 本地壁纸路径（首页在根目录）
	var url = "assets/" + imgUrl;
	panel.style.backgroundImage = "url('" + url.replace(/['\\]/g, '\\$&') + "')";
	panel.style.backgroundPosition = "center center";
	panel.style.backgroundRepeat = "no-repeat";
	panel.style.backgroundColor = "#666";
	panel.style.backgroundSize = "cover";
	sessionStorage.setItem(indexName, index);
}

function decryptEmail(encoded) {
	var address = atob(encoded);
	window.location.href = "mailto:" + address;
}

document.addEventListener('DOMContentLoaded', function () {
	// 获取一言数据
	fetch("https://v1.hitokoto.cn")
		.then(function(response) {
			return response.json();
		})
		.then(function(res) {
			var descElement = document.getElementById('description');
			if (descElement && res.hitokoto && res.from) {
				// Create text nodes to prevent XSS
				var textNode = document.createTextNode(res.hitokoto);
				var br = document.createElement('br');
				var fromText = document.createTextNode(' -「');
				var strong = document.createElement('strong');
				strong.textContent = res.from;
				var endText = document.createTextNode('」');

				descElement.innerHTML = '';
				descElement.appendChild(textNode);
				descElement.appendChild(br);
				descElement.appendChild(fromText);
				descElement.appendChild(strong);
				descElement.appendChild(endText);
			}
		})
		.catch(function(error) {
			console.error('Error fetching hitokoto:', error);
		});

	var iUpElements = document.querySelectorAll(".iUp");
	for (var i = 0; i < iUpElements.length; i++) {
		iUp.up(iUpElements[i]);
	}

	var avatarElement = document.querySelector(".js-avatar");
	if (avatarElement) {
		avatarElement.addEventListener('load', function () {
			avatarElement.classList.add("show");
		});
	}
	(function () {
		function loadScriptOnce(url, callback) {
			var existing = document.querySelector('script[data-bing-images-loader="1"]');
			if (existing) {
				callback(null);
				return;
			}
			var script = document.createElement('script');
			script.setAttribute('data-bing-images-loader', '1');
			script.src = url + '?t=' + new Date().getTime();
			script.onload = function () { callback(null); };
			script.onerror = function () { callback(new Error('Failed to load ' + url)); };
			document.body.appendChild(script);
		}

		loadScriptOnce('./assets/json/images.js', function (err) {
			if (!err && typeof getBingImages === 'function' && window.BING_IMAGES && Array.isArray(window.BING_IMAGES)) {
				getBingImages(window.BING_IMAGES);
			}
		});
})();
});

// ====== About 弹窗 ======
function showAboutModal() {
	var modal = document.getElementById('aboutModal');
	if (modal) modal.classList.add('open');
}

function closeAboutModal() {
	var modal = document.getElementById('aboutModal');
	if (modal) modal.classList.remove('open');
}

document.addEventListener('keydown', function (e) {
	if (e.key === 'Escape') closeAboutModal();
});

var btnMobileMenu = document.querySelector('.btn-mobile-menu__icon');
var navigationWrapper = document.querySelector('.navigation-wrapper');

if (btnMobileMenu && navigationWrapper) {
	btnMobileMenu.addEventListener('click', function () {
		var isVisible = navigationWrapper.classList.contains('visible');
		
		function handleAnimationEnd() {
			navigationWrapper.classList.remove('visible', 'animated', 'bounceOutUp');
			navigationWrapper.removeEventListener('animationend', handleAnimationEnd);
		}
		
		if (isVisible) {
			navigationWrapper.addEventListener('animationend', handleAnimationEnd);
			navigationWrapper.classList.remove('bounceInDown');
			navigationWrapper.classList.add('animated', 'bounceOutUp');
		} else {
			navigationWrapper.classList.add('visible', 'animated', 'bounceInDown');
		}
		
		btnMobileMenu.classList.toggle('icon-list');
		btnMobileMenu.classList.toggle('icon-angleup');
	});
}
